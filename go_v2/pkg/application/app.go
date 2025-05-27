package application

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

type Application struct {
	TargetDir string
	State     string
	Logs      []*LogMessage
	LogsCh    chan LogMessage
	Diffs     []*Diff
	HasVSC    bool
}

func NewApplication() *Application {
	a := Application{
		State:  APP_STATE_NEW,
		Logs:   make([]*LogMessage, 0),
		LogsCh: make(chan LogMessage, 6),
	}

	go func() {
		for l := range a.LogsCh {
			a.Logs = append(a.Logs, &l)
			if l.Level == LOG_TAG_ERROR {
				a.State = APP_STATE_ERROR
			}
		}
	}()

	return &a
}

func (a *Application) Reset() {
	fmt.Println("[APP] Reset")
	a.State = APP_STATE_NEW
	a.TargetDir = ""
	a.Diffs = []*Diff{}
	a.Logs = []*LogMessage{}
}

func (a *Application) Log(level string, msg string, args ...any) {
	message := fmt.Sprintf("%s %s", time.Now().Format("15:04:05"), fmt.Sprintf(msg, args...))
	fmt.Println("[APP] LOG:", message)

	a.LogsCh <- LogMessage{message, level}
}

func (a *Application) StartInstallation(targetDir string, components []InstallOption, needBackup bool) {
	// Downloads actual components from GitHub repo, unzips and copies files to temporary directory.
	// Then backups data in 'dir' if it differs from actual framework, generated diff files for manual analysys.
	// Then copies actual framework into the 'dir'

	a.TargetDir = targetDir

	zipDirectory(targetDir)

	a.Log(LOG_TAG_INFO, "Установка началась")
	a.State = APP_STATE_INSTALL

	a.Log(LOG_TAG_INFO, "Создаем временную директорию")
	os.RemoveAll(TEMP_DIR)
	os.Mkdir(TEMP_DIR, os.ModeAppend)

	a.Log(LOG_TAG_INFO, "<b>Целевая директория</b>: %s", targetDir)
	a.Log(LOG_TAG_INFO, "<b>Нужен бекап</b>: %t", needBackup)
	a.Log(LOG_TAG_INFO, "--------------------")

	a.DownloadComponents(TEMP_DIR, components)
	a.Log(LOG_TAG_INFO, "--------------------")
	a.State = APP_STATE_DOWNLOADED

	a.UpdateOriginalMission(targetDir)
	a.Log(LOG_TAG_INFO, "--------------------")

	if needBackup {
		// -- Backup files and gather diff info
		a.BackupCustomFiles(targetDir, TEMP_DIR)

		// -- If there are diffs to resolve - stop process and wait for resolution via UI
		if len(a.Diffs) != 0 {
			a.HasVSC = checkVsCodeExists()
			a.State = APP_STATE_DIFF
			return
		}
	}

	// -- Copy tmp folder to main
	a.FinishInstallation()

	a.State = APP_STATE_DONE
}

func (a *Application) FinishInstallation() {
	a.ResolveAllDiffs()

	a.Log(LOG_TAG_INFO, "Копируем новые файлы в директорию миисии!")
	err := copyDir(TEMP_DIR, a.TargetDir)
	if err != nil {
		panic(err)
	}
}

func checkVsCodeExists() bool {
	output, err := exec.Command("cmd", "/C", "where", "code").Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return false
	}
	fmt.Println(string(output))
	fmt.Println(string(output) != "")

	return string(output) != ""
}
