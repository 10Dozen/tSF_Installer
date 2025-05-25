package application

import (
	"fmt"
	"log"
	"os"
	"path"
)

type UpdateFileStrategy struct {
	filename string
	newName  string
}

type UpdateStrategy struct {
	dir    []string
	target string
	files  []UpdateFileStrategy
}

func NewtSFModuleUpdateStrategy(dir string, ext string, extraFiles ...UpdateFileStrategy) UpdateStrategy {
	files := []UpdateFileStrategy{
		{fmt.Sprintf("Settings.%s", ext), fmt.Sprintf("%s.%s", dir, ext)},
	}

	if len(extraFiles) > 0 {
		files = append(files, extraFiles...)
	}

	s := UpdateStrategy{
		dir:    []string{"dzn_tSFramework", "Modules", dir},
		target: "",
		files:  files,
	}

	return s
}

var UpdateStrategies = []UpdateStrategy{
	// -- Gear, before detached config update
	{
		dir:    []string{"dzn_gear"},
		target: "Gear",
		files: []UpdateFileStrategy{
			{"GearAssignementTable.sqf", ""}, // Copy with the same name
			{"Kits.sqf", ""},
			{"Settings.sqf", ""},
		},
	},

	// -- Dynai, before detached config update
	{
		dir:    []string{"dzn_dynai"},
		target: "Dynai",
		files: []UpdateFileStrategy{
			{"Settings.sqf", ""},
			{"Zones.sqf", ""},
		},
	},

	// -- tSF in 2.0.13 version, before detached config update
	{
		dir:    []string{"dzn_tSFramework"},
		target: "",
		files:  []UpdateFileStrategy{{"Settings.yaml", "_Modules.yaml"}},
	},
	NewtSFModuleUpdateStrategy("ACEActions", "sqf"),
	NewtSFModuleUpdateStrategy("AirborneSupport", "sqf"),
	NewtSFModuleUpdateStrategy("ArtillerySupport", "sqf"),
	NewtSFModuleUpdateStrategy("Authorization", "sqf"),
	NewtSFModuleUpdateStrategy("Briefing", "sqf", UpdateFileStrategy{"tSF_Briefing.sqf", "_Briefing.sqf"}),
	NewtSFModuleUpdateStrategy("CCP", "sqf"),
	NewtSFModuleUpdateStrategy("Chatter", "yaml"),
	NewtSFModuleUpdateStrategy("Conversations", "sqf"),
	NewtSFModuleUpdateStrategy("CrewOptions", "yaml"),
	NewtSFModuleUpdateStrategy("EditorRadioSettings", "yaml"),
	NewtSFModuleUpdateStrategy("EditorUnitBehavior", "sqf"),
	NewtSFModuleUpdateStrategy("EditorVehicleCrew", "yaml"),
	NewtSFModuleUpdateStrategy("FARP", "sqf"),
	NewtSFModuleUpdateStrategy("Interactives", "sqf"),
	NewtSFModuleUpdateStrategy("IntroText", "yaml"),
	NewtSFModuleUpdateStrategy("JIPTeleport", "yaml"),
	NewtSFModuleUpdateStrategy("MissionConditions", "yaml", UpdateFileStrategy{"Endings.hpp", "_MissionEndings.hpp"}),
	NewtSFModuleUpdateStrategy("MissionDefaults", "yaml"),
	NewtSFModuleUpdateStrategy("POM", "sqf"),
	NewtSFModuleUpdateStrategy("Respawn", "yaml"),
	NewtSFModuleUpdateStrategy("tSAdminTools", "sqf"),
	NewtSFModuleUpdateStrategy("tSNotes", "yaml"),
	NewtSFModuleUpdateStrategy("tSSettings", "yaml"),
}

func (a *Application) UpdateOriginalMission(rootDir string) {
	a.Log(LOG_TAG_INFO, "Обновляем целевую директорию под формат \"открепленных\" конфигов...")

	dirs := []string{
		path.Join(rootDir, "Config"),
		path.Join(rootDir, "Config", "Gear"),
		path.Join(rootDir, "Config", "Dynai"),
	}

	// -- Exit if Config already exists in root
	exists := PathExists(dirs[0])
	a.Log(LOG_TAG_INFO, "&nbsp;&nbsp;Конфиги уже откреплены")
	if exists {
		return
	}

	// -- Create Config dir
	a.Log(LOG_TAG_INFO, "&nbsp;&nbsp;Создаем директорию для открепленных конфигов")
	for _, dirToCreate := range dirs {
		if err := os.Mkdir(dirToCreate, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	// -- Rename settings files and move to Config
	a.Log(LOG_TAG_INFO, "&nbsp;&nbsp;Открепляем конфиги миссии")
	const numWorkers = 4
	numJobs := len(UpdateStrategies)
	jobs := make(chan UpdateStrategy, numJobs)
	results := make(chan int, numJobs)

	for w := 1; w <= numWorkers; w++ {
		go updateWorker(a, rootDir, dirs[0], jobs, results)
	}
	for _, strategy := range UpdateStrategies {
		jobs <- strategy
	}
	close(jobs)

	for a := 1; a <= numJobs; a++ {
		<-results
	}
	close(results)

	// -- Delete old component versions
	a.Log(LOG_TAG_INFO, "&nbsp;&nbsp;Удаляем старые версии компонентов:")
	for _, dirname := range []string{SLUG_COMMON_FUNCTION, SLUG_GEAR, SLUG_DYNAI, SLUG_TSFRAMEWORK, SLUG_BRV} {
		dirToRemove := path.Join(rootDir, dirname)
		a.Log(
			LOG_TAG_INFO,
			"&nbsp;&nbsp;&nbsp;Удаляем директорию "+
				"<span class='tSF-Remove'>%s</span>",
			dirToRemove,
		)
		os.RemoveAll(dirToRemove)
	}

	a.Log(LOG_TAG_INFO, "Целевая директория успешно обновлена")
}

func updateWorker(a *Application, rootDir, targetDir string, jobs <-chan UpdateStrategy, results chan<- int) {
	for strategy := range jobs {
		source := path.Join(rootDir, path.Join(strategy.dir...))
		target := path.Join(targetDir, strategy.target)

		for _, fileInfo := range strategy.files {
			fileCurrent := path.Join(source, fileInfo.filename)
			if fileInfo.newName == "" {
				fileInfo.newName = fileInfo.filename
			}
			fileNew := path.Join(target, fileInfo.newName)

			a.Log(
				LOG_TAG_INFO,
				"&nbsp;&nbsp;Перемещаем конфиг "+
					"<span class='tSF-MoveFrom'>%s</span> "+
					"в <span class='tSF-MoveTo'>%s</span>",
				path.Join(path.Join(strategy.dir...), fileInfo.filename),
				path.Join(strategy.target, fileInfo.newName),
			)
			copyFile(fileCurrent, fileNew)
		}

		results <- 1
	}
}
