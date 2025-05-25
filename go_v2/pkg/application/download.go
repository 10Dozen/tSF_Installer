package application

import (
	"os"
	"slices"
)

func (a *Application) DownloadComponents(targetDir string, components []InstallOption) {
	// Downloads all components and unzip content to temporary dir
	channels := [len(componentsOrder)]chan string{}
	for idx, componentSlug := range componentsOrder {
		componentIdx := slices.IndexFunc(components, func(c InstallOption) bool { return c.Slug == componentSlug })
		if componentIdx < 0 {
			panic("Failed to find component!")
		}
		component := components[componentIdx]
		if !component.Checked {
			continue
		}

		channels[idx] = make(chan string)
		a.Log(LOG_TAG_INFO, "Компонент <b>%s</b>: <a href='%s'>%s</a>", component.Label, component.Value, component.Value)

		go func() {
			a.Log(LOG_TAG_INFO, "Скачиваем компонент <b>%s</b>", component.Label)
			channels[idx] <- downloadComponent(components[componentIdx].Slug, components[componentIdx].Value)
			a.Log(LOG_TAG_INFO, "Компонент <b>%s</b> успешно скачан", component.Label)
		}()
	}

	// Wait for go-routines to finish and sequentially copy unziped data to Temp folder
	for _, ch := range channels {
		if ch == nil {
			continue
		}
		dir := <-ch
		close(ch)

		composeComponentAtTempDirectory(dir, TEMP_DIR)
		os.RemoveAll(dir)
	}

	a.Log(LOG_TAG_INFO, "Компоненты собраны во временной директории.")
}
