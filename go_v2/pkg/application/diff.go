package application

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

func (a *Application) ResolveDiffViaVSCode(idx int) {
	diff := a.Diffs[idx]

	_, err := exec.Command("cmd", "/C", "code", "-w", "-d", diff.OriginalFile, diff.NewFile).Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
	}

	diff.Resolution = DIFF_RESOLUTION_NEW
	a.State = APP_STATE_DIFF
}

func (a *Application) ResolveDiff(idx int, resolution string) {
	diff := a.Diffs[idx]
	diff.Resolution = resolution
}

func (a *Application) ResolveAllDiffs() {
	for _, diff := range a.Diffs {
		resolution := diff.Resolution
		if resolution == "" {
			resolution = DIFF_RESOLUTION_BACKUP
		}
		a.Log(LOG_TAG_INFO, "&nbsp;&nbsp;Разрешаем конфликт: <span class='tSF-MoveFrom'>%s</span> => %s", diff.OriginalFile, resolution)

		fmt.Println("[ResolveAllDiffs] Diff=", diff.NewFile, "Resolution=", resolution)

		if resolution == DIFF_RESOLUTION_BACKUP {
			filename := path.Base(diff.OriginalFile)
			dir := path.Dir(diff.OriginalFile)
			fmt.Println("Resolve", diff.OriginalFile, "as [BACKUP]: Going to backup: ", filename, dir)
			backupFile(filename, dir)
		} else if resolution == DIFF_RESOLUTION_NEW {
			fmt.Println("Resolve", diff.OriginalFile, "as [USE NEW]: Going to delete: ", diff.OriginalFile)
			os.RemoveAll(diff.OriginalFile)
		} else if resolution == DIFF_RESOLUTION_OLD {
			fmt.Println("Resolve", diff.OriginalFile, "as [USE OLD]: Going to delete: ", diff.NewFile)
			os.RemoveAll(diff.NewFile)
		} else {
			panic("Invalid diff resolution " + resolution)
		}
	}
}
