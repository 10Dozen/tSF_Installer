package application

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type BackupStrategy struct {
	dir    []string
	always []string
	onDiff []string
}

var backupStrategy = map[string]BackupStrategy{
	SLUG_ROOT: {
		dir: []string{""},
		always: []string{
			"overview.jpg",
		},
		onDiff: []string{
			"init.sqf",
			"initServer.sqf",
			"description.ext",
		},
	},
	SLUG_CONFIG_GEAR: {
		dir: []string{SLUG_CONFIG, SLUG_CONFIG_GEAR},
		always: []string{
			"Kits.sqf",
			"GearAssignementTable.sqf",
			"GearAssignmentTable.yaml",
		},
		onDiff: []string{
			"Settings.sqf",
			"PluginSettings.yaml",
		},
	},
	SLUG_DYNAI: {
		dir:    []string{SLUG_CONFIG, SLUG_DYNAI},
		always: nil,
		onDiff: []string{
			"Settings.sqf",
			"Zones.sqf",
		},
	},

	SLUG_CONFIG: {
		dir:    []string{SLUG_CONFIG},
		always: nil,
		onDiff: []string{"*"},
	},
}

/* -- For each strategy - run worker

	   -- Compose diff info as slice:
	   		{
	   			pathnameCurrent string,  ((missionRoot ommited)/Config/file.sqf)
				pathnameNew string,      ((missionRoot ommited)/Config/file.sqf)
				diffHtml string,
			}

       -- On UI render page as /diff/{index}
	      Get from slice by index and render
*/

func (a *Application) BackupCustomFiles(rootDir, tempDir string) {
	// -- Rename settings files and move to Config
	a.Log(LOG_TAG_INFO, "Делаем резервные копии измененных файлов миссии")

	const numWorkers = 1
	numJobs := len(backupStrategy)
	jobs := make(chan BackupStrategy, numJobs)
	results := make(chan []*Diff, numJobs)

	for w := 1; w <= numWorkers; w++ {
		go backupWorker(a, rootDir, tempDir, jobs, results)
	}
	for _, strategy := range backupStrategy {
		jobs <- strategy
	}
	close(jobs)

	totalDiffs := make([]*Diff, 0)
	for a := 1; a <= numJobs; a++ {
		totalDiffs = append(totalDiffs, <-results...)
	}
	a.Diffs = totalDiffs
	close(results)
}

func backupWorker(a *Application, originalDir, newDir string, jobs <-chan BackupStrategy, results chan<- []*Diff) {
	/*
	   Scan (rootDir + BackupStrategy.dir) for all files in directory by filter "onDoff"

	   	If file in "always" files - rename to *.backup.*
	   	If file match "onDiff" (e.g. "file.name" or "*.*")
	   	Otherwise:
	   	    - Find same file at tempDir:
	   			If cur file is .sqf -> find same .sqf,
	   			If cur file is .yaml -> find same .yaml
	   				if exists -- make DIFF

	   				if not exists and file .sqf -> find .yaml
	   					if exists -- rename to *.backup.*
	   					otherwise -- error?
	*/
	for strategy := range jobs {
		targetSubPath := path.Join(strategy.dir...)
		originalDir := path.Join(originalDir, targetSubPath)
		newDir := path.Join(newDir, targetSubPath)

		fmt.Println("[BackupWorker] Strategy:", strategy)
		fmt.Println("[BackupWorker] Taget subpath=", targetSubPath, ", Original dir=", originalDir, ", New dir=", newDir)

		filename := ""
		diffData := make([]*Diff, 0)

		items, _ := os.ReadDir(originalDir)
		for _, item := range items {

			fmt.Println("[BackupWorker] backup item=", item)
			if item.IsDir() {
				continue
			}

			filename = item.Name()

			fmt.Println("[BackupWorker] filename=", filename)

			// File is ALWAYS backup
			if slices.Contains(strategy.always, filename) {
				fmt.Println("[BackupWorker] always backup file! =", filename)
				backupFile(filename, originalDir)
				continue
			}

			for _, diffFilePattern := range strategy.onDiff {

				fmt.Println("[BackupWorker] on diff backup file! =", filename, diffFilePattern)
				match, _ := path.Match(diffFilePattern, filename)
				if !match {
					fmt.Println("[BackupWorker] Path not matched. Continue...")
					continue
				}
				fmt.Println("[BackupWorker] Path matched")

				// Same extension - save diff for user to decide

				fmt.Println(
					"[BackupWorker] Going to check - if file exists in new dir:",
					path.Join(newDir, filename),
					"=>",
					PathExists(path.Join(newDir, filename)),
				)
				if PathExists(path.Join(newDir, filename)) {
					diff := getFilesDiffer(filename, originalDir, newDir)
					if diff != nil {
						fmt.Println("[BackupWorker] Diff is not empty - save diffData")
						diffData = append(diffData, diff)
					}
					break
				}

				fmt.Println("[BackupWorker] Different extension")
				// Different extension (sqf vs yaml) or file is not used anymore
				backupFile(filename, originalDir)
			}
		}

		results <- diffData
	}
}

func backupFile(filename, dir string) {
	// Renames given 'src'/'filename' and add '.backup.' suffix before extension
	// e.g. MyFile.txt -> MyFile.backup.txt
	fmt.Println("[BackupWorker.backupFile] Filename=", filename, "Dir=", dir)
	if strings.Contains(filename, "backup") {
		return
	}

	parts := strings.Split(filename, ".")

	fmt.Println("[BackupWorker.backupFile] parts=", parts)
	fmt.Println("[BackupWorker.backupFile] Going to rename=", path.Join(dir, filename), "to=", path.Join(dir, fmt.Sprintf("%s.backup.%s", parts[0], parts[1])))
	os.Rename(
		path.Join(dir, filename),
		path.Join(dir, fmt.Sprintf("%s.backup.%s", parts[0], parts[1])),
	)
}

func getFilesDiffer(filename, src, test string) *Diff {
	// Tests 2 files ('src/filename' vs 'test/filename') to have differences.
	// If difference found - returns diff data,
	// if there is no 'test/filename' or no differences found - returns empty diff data
	// Returns:
	//   bool   Flag that files differs
	//   string Diff data (HTML formatted) or empty sring

	currentFile := path.Join(src, filename)
	testAgainstFile := path.Join(test, filename)

	if checkFileHashesEquals(currentFile, testAgainstFile) {
		return nil
	}

	contents := []string{"", ""}
	for i, f := range []string{currentFile, testAgainstFile} {
		file, err := os.Open(f)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			panic(err)
		}
		contents[i] = string(content)
	}

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(contents[0], contents[1], false)
	diffs = dmp.DiffCleanupSemantic(diffs)
	diffs = dmp.DiffCleanupMerge(diffs)

	return &Diff{
		OriginalFile: currentFile,
		NewFile:      testAgainstFile,
		Diff:         formatDiffHTML(dmp.DiffPrettyHtml(diffs)),
	}
}

func checkFileHashesEquals(a, b string) bool {
	// Check that hashes of the files 'a' and 'b' match.
	// Return:
	//   bool True if hashes match, otherwise - false

	hashes := [2]string{"", ""}
	for i, f := range [2]string{a, b} {
		file, err := os.Open(f)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		hash := sha256.New()
		if _, err := io.Copy(hash, file); err != nil {
			log.Fatal(err)
		}
		hashes[i] = hex.EncodeToString(hash.Sum(nil))
	}
	return hashes[0] == hashes[1]
}

func formatDiffHTML(content string) string {
	// Apply additional formatting to HTML diff data ( line numbers)
	// Return:
	//   string HTML content

	var indexed []string
	for i, l := range strings.Split(content, "<br>") {
		indexed = append(
			indexed,
			fmt.Sprintf(
				"%3d | %s",
				i+1,
				strings.Replace(l, "&para;", "", 1),
			),
		)
	}
	return fmt.Sprintf(`<div>%s</div>`, strings.Join(indexed, "<br>"))
}
