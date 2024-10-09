/*
Copyright © 2024 10Dozen <steel.mjr@gmai.com>
*/
package cmd

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// installCmd represents the install command
const (
	DZN_COMMON_FUNCTIONS = "dzn_commonFunctions"
	DZN_GEAR             = "dzn_gear"
	DZN_DYNAI            = "dzn_dynai"
	DZN_TSFRAMEWORK      = "dzn_tSFramework"
	DZN_BRV              = "dzn_brv"

	HTTPS_PREFIX           = "https://"
	GITHUB_DOWNLOAD_SUFFIX = "archive/refs/heads/%s.zip"
	GITHUB_DEFAULT_BRANCH  = "master"

	TEMP_DIR = "FrameworkLatest"
)

type backupStrategy struct {
	dir    string
	deep   bool
	wipe   bool
	always []string
	onDiff []string
}

var (
	cfgFile  string
	defaults = map[string]string{
		DZN_COMMON_FUNCTIONS: "github.com/10Dozen/dzn_commonFunctions",
		DZN_GEAR:             "github.com/10Dozen/dzn_gear",
		DZN_DYNAI:            "github.com/10Dozen/dzn_dynai",
		DZN_TSFRAMEWORK:      "github.com/10Dozen/dzn_tSFramework",
	}
	updateStrategy = map[string]backupStrategy{
		"root": backupStrategy{
			dir:    "",
			deep:   false,
			wipe:   false,
			always: []string{"overview.jpg"},
			onDiff: []string{
				"init.sqf",
				"initServer.sqf",
				"description.ext",
			},
		},
		DZN_BRV: backupStrategy{
			dir:  DZN_BRV,
			deep: false,
			wipe: true,
		},
		DZN_COMMON_FUNCTIONS: backupStrategy{
			dir:  DZN_COMMON_FUNCTIONS,
			deep: false,
			wipe: true,
		},
		DZN_GEAR: backupStrategy{
			dir:  DZN_GEAR,
			deep: false,
			wipe: true,
			always: []string{
				"Kits.sqf",
				"GearAssignementTable.sqf",
			},
			onDiff: []string{
				"Settings.sqf",
			},
		},
		DZN_DYNAI: backupStrategy{
			dir:  DZN_DYNAI,
			deep: false,
			wipe: true,
			always: []string{
				"Zones.sqf",
			},
			onDiff: []string{
				"Settings.sqf",
			},
		},
		DZN_TSFRAMEWORK: backupStrategy{
			dir:  DZN_TSFRAMEWORK,
			deep: true,
			wipe: true,
			always: []string{
				"tSF_briefing.sqf",
				"Endings.hpp",
			},
			onDiff: []string{
				"Settings.yaml",
				"Settings.sqf",
				"CCP Compositions.sqf",
				"FARP Compositions.sqf",
			},
		},
	}
	// --- Order of applying components
	repoOrder = [4]string{
		DZN_COMMON_FUNCTIONS,
		DZN_GEAR,
		DZN_DYNAI,
		DZN_TSFRAMEWORK,
	}
	// --- Order of files on update validation
	updateOrder = [6]backupStrategy{
		updateStrategy["root"],
		updateStrategy[DZN_COMMON_FUNCTIONS],
		updateStrategy[DZN_BRV],
		updateStrategy[DZN_GEAR],
		updateStrategy[DZN_DYNAI],
		updateStrategy[DZN_TSFRAMEWORK],
	}

	installCmd = &cobra.Command{
		Use:   "install",
		Short: "Installs the tSFramework to given directory",
		Long: `Downloads, unpack and install required components of the tSFramewrok:
- dzn_commonFunctions
- dzn_gear
- dzn_dynai
- dzn_tSFramework
	
By default installs current version from main branch, but exact repo may be specified by flags.`,
		Run: func(cmd *cobra.Command, args []string) {
			install()
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	installCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tSFInstaller.yaml)")

	installCmd.PersistentFlags().String("dir", "", "output directory (where to install, required)")
	installCmd.PersistentFlags().StringP(DZN_COMMON_FUNCTIONS, "c", "", "repo path to dzn_commonFunctions")
	installCmd.PersistentFlags().StringP(DZN_GEAR, "g", "", "repo path to dzn_gear")
	installCmd.PersistentFlags().StringP(DZN_DYNAI, "d", "", "repo path to dzn_dynai")
	installCmd.PersistentFlags().StringP(DZN_TSFRAMEWORK, "f", "", "repo path to dzn_tSFramework")
	installCmd.PersistentFlags().Bool("nobackup", false, "flag to skip backup of the existing files in dir and force override files")

	viper.BindPFlag("dir", installCmd.PersistentFlags().Lookup("dir"))
	viper.BindPFlag("override", installCmd.PersistentFlags().Lookup("nobackup"))
	installCmd.MarkPersistentFlagRequired("dir")

	rootCmd.AddCommand(installCmd)
}

func initConfig() {
	viper.SetConfigFile(".tSFInstaller.yaml")

	if cfgFile != "" {
		// Use config passed as parameter
		viper.SetConfigFile(cfgFile)
	}

	// Defaults - allow configs to contain only diff
	viper.SetDefault(DZN_COMMON_FUNCTIONS, defaults[DZN_COMMON_FUNCTIONS])
	viper.SetDefault(DZN_GEAR, defaults[DZN_GEAR])
	viper.SetDefault(DZN_DYNAI, defaults[DZN_DYNAI])
	viper.SetDefault(DZN_TSFRAMEWORK, defaults[DZN_TSFRAMEWORK])

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else if cfgFile != "" {
		// Config was provided but not found
		panic(err)
	} else {
		// In case no config was provided and default config not yet exists
		viper.WriteConfig()
		println("Default config generated!")
	}
}

// Functions
func install() {
	// Main method of the Install command
	// Downloads actual components from GitHub repo, unzips and copies files to temporary directory.
	// Then backups data in 'dir' if it differs from actual framework, generated diff files for manual analysys.
	// Then copies actual framework into the 'dir'

	os.RemoveAll(TEMP_DIR)
	os.Mkdir(TEMP_DIR, os.ModeAppend)

	log.Println("Downloading repositories...")
	channels := [len(repoOrder)]chan string{}
	for idx, reponame := range repoOrder {
		channels[idx] = make(chan string)
		go func() {
			channels[idx] <- downloadRepo(reponame, viper.GetString(reponame))
		}()
	}
	// Wait for go-routines to finish and sequentially copy unziped data to Temp folder
	for _, ch := range channels {
		dir := <-ch
		close(ch)
		copyUnzippedDirectory(dir, TEMP_DIR)
		os.RemoveAll(dir)
	}

	outDir := viper.GetString("dir")
	if !viper.GetBool("override") {
		log.Println("Preparing existing installation")
		prepareUpdate(outDir, TEMP_DIR)
	}

	log.Println("Installing new framework to target directory")
	err := os.CopyFS(outDir, os.DirFS(TEMP_DIR))
	if err != nil {
		panic(err)
	}
}

func downloadRepo(name, url string) string {
	// Downloads repo of given 'name' and given 'url'.
	// URL may be modified by CLI params, by default master branches will be used.
	// Returns:
	//   string - temp directory of the downloaded and uncompressed repo.

	fqUrl := getDownloadUrl(name, url)
	log.Printf("Repo [%s], fqUrl = %s\n", name, fqUrl)

	err := downloadFile(name, fqUrl)
	if err != nil {
		panic(err)
	}
	dir, err := os.MkdirTemp("", "temp")
	if err != nil {
		panic(err)
	}

	err = unzipFile(name, dir)
	if err != nil {
		panic(err)
	}

	os.Remove(name)

	return dir
}

func getDownloadUrl(name, url string) string {
	// Pre-format url with protocol, branches and path to actual file at GitHub
	// Basically converts GitHub web url of the branch into link a branch archive file
	// Return:
	//   string URL to selected repository/branch archive file
	fqUrl := completeUrl(url)
	if fqUrl == completeUrl(defaults[name]) {
		fqUrl += fmt.Sprintf(GITHUB_DOWNLOAD_SUFFIX, GITHUB_DEFAULT_BRANCH)
		return fqUrl
	}

	// Non-default repo path
	urlParts := strings.Split(fqUrl, "/")
	branchName := urlParts[6]
	return strings.Join(urlParts[0:5], "/") + "/" + fmt.Sprintf(GITHUB_DOWNLOAD_SUFFIX, branchName)
}

func completeUrl(url string) string {
	// Returns fully qualified URL (with protocol prefix and valid end char)
	fqUrl := url
	if !strings.HasPrefix(url, HTTPS_PREFIX) {
		fqUrl = HTTPS_PREFIX + fqUrl
	}
	lastChar := string(fqUrl[len(fqUrl)-1])
	if lastChar != "/" {
		fqUrl += "/"
	}
	return fqUrl
}

func downloadFile(filepath, url string) error {
	// Downloads content from given 'url' and writes to 'filepath' file on disk
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func unzipFile(source, dest string) error {
	// Uncompresses files/dirs from 'source' archive and puts them to 'dest' directory
	read, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer read.Close()

	for _, file := range read.File {
		if file.Mode().IsDir() {
			continue
		}
		open, err := file.Open()
		if err != nil {
			return err
		}
		name := path.Join(dest, file.Name)
		os.MkdirAll(path.Dir(name), os.ModeDir)
		create, err := os.Create(name)
		if err != nil {
			return err
		}
		defer create.Close()
		create.ReadFrom(open)
	}
	return nil
}

func copyUnzippedDirectory(source, dest string) {
	// Copies content of uncompressed file at 'source' to 'dest'
	// Files are stored in directory with archive name,
	// so it's neccessary to dive into 'source' directory and copy
	// from there.

	f, err := os.Open(source)
	if err != nil {
		panic(err)
	}
	fileInfo, err := f.ReadDir(-1)
	dirToCopy := path.Join(source, fileInfo[0].Name())

	err = os.CopyFS(dest, os.DirFS(dirToCopy))
	if err != nil {
		panic(err)
	}
}

func prepareUpdate(source, temp string) {
	// Checks 'source' directory (old mission folder) against 'temp' directory (newly downloaded framework).
	// Scans for difference in specified files and makes backups and diff files if old files were changed.
	order := updateOrder

	channels := [len(order)]chan bool{}
	for idx, strt := range order {
		channels[idx] = make(chan bool)
		go func() {
			channels[idx] <- handleDirectory(source, temp, strt.dir, strt.deep, strt.wipe, strt.always, strt.onDiff)
		}()
	}

	for _, ch := range channels {
		<-ch
		close(ch)
	}
}

func handleDirectory(source, temp, subdir string, deep, wipe bool, toBackup, toBackupOnDiff []string) bool {
	// Scans directory 'source'/'subdir' according to given rules.
	// For directories - if 'deep' is true - invokes recursive deep scan of all nested files/dirs.
	// For files - check against 'toBackup' listing and backup mathing files, otherwise continue
	//             AND check against 'toBackupDiff' listing and check for diffs against same file in 'temp'/'subdir'.
	//             Backup and create diff file if not match, otherwise - ignore file.
	// If 'wipe' flag is true - all dirs/files without backuped files will be deleted.
	// Return:
	//   bool Flag that directory has backuped files.

	targetDir := path.Join(source, subdir)
	compareToDir := path.Join(temp, subdir)

	saveDirectory := false
	entries, err := os.ReadDir(targetDir)
	if err != nil {
		// Dir may be empty/not exists during fresh installation
		return saveDirectory
	}

	for _, e := range entries {
		name := e.Name()

		// Check dirs recursevily
		if e.IsDir() {
			toBeSaved := false
			if deep {
				toBeSaved = handleDirectory(targetDir, compareToDir, name, deep, wipe, toBackup, toBackupOnDiff)
				saveDirectory = saveDirectory || toBeSaved
			}

			// Remove directory only if there is nothing to save in it
			if wipe && !toBeSaved {
				os.RemoveAll(path.Join(targetDir, name))
			}

			continue
		}

		if slices.Contains(toBackup, name) {
			backupFile(name, targetDir)
			saveDirectory = true
			continue
		}
		if slices.Contains(toBackupOnDiff, name) {
			differs, details := checkFilesDiffer(name, targetDir, compareToDir)
			if !differs {
				if wipe {
					os.Remove(path.Join(targetDir, name))
				}
				continue
			}

			log.Printf("  [handleDirectory] File [%s] differs, making backup\n", name)
			backupFile(name, targetDir)
			saveDirectory = true
			// Write diff if present
			if details != "" {
				os.WriteFile(
					path.Join(targetDir, fmt.Sprintf("%s-diff.html", name)),
					[]byte(details),
					0666,
				)
			}

			continue
		}

		if wipe {
			err := os.Remove(path.Join(targetDir, name))
			if err != nil {
				panic(err)
			}
		}
	}

	if saveDirectory {
		log.Printf("  [handleDirectory] Directory [%s] - some files was backuped\n", targetDir)
	}

	return saveDirectory
}

func backupFile(filename, src string) {
	// Renames given 'src'/'filename' and add '.backup.' suffix before extension
	// e.g. MyFile.txt -> MyFile.backup.txt
	parts := strings.Split(filename, ".")
	os.Rename(
		path.Join(src, filename),
		path.Join(src, fmt.Sprintf("%s.backup.%s", parts[0], parts[1])),
	)
}

func checkFilesDiffer(filename, src, test string) (bool, string) {
	// Tests 2 files ('src/filename' vs 'test/filename') to have differences.
	// If difference found - returns diff data,
	// if there is no 'test/filename' or no differences found - returns empty diff data
	// Returns:
	//   bool   Flag that files differs
	//   string Diff data (HTML formatted) or empty sring

	currentFile := path.Join(src, filename)
	testAgainstFile := path.Join(test, filename)

	// Check that same file exists in fresh version of the Framework
	// if not - exit with FileDiffer=True and empty diff, in order to backup current file
	// otherwise - compare files
	if _, err := os.Stat(testAgainstFile); errors.Is(err, os.ErrNotExist) {
		return true, ""
	}

	if checkFileHashesEquals(currentFile, testAgainstFile) {
		return false, ""
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
	return true, formatDiffHTML(dmp.DiffPrettyHtml(diffs))
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
	// Apply additional formatting to HTML diff data (styles, line numbers)
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
	content = strings.Join(indexed, "<br>")
	return fmt.Sprintf(`<html>
		<head>
			<title>Diff</title>
			<style> 
				body {
					background-color: #494949;
					font-family: monospace;
					font-size: large;
					color: #e1e1e1;
					white-space:pre;
				}
				ins {
					color: #3d3d3d;
					text-decoration: none;
				}
				del {
					color: #3d3d3d;
					background-color: #ffa2a2!important;
					text-decoration: none;
				}
			</style>
		</head>
		<body>%s</body>
		</html>
	`, content)
}
