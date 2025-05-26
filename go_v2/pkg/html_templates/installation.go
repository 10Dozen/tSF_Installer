package html_templates

const INSTALLATION = `
<!-- Install form -->
<div class="col-md-9 offset-md-1" hx-target="#main" hx-on::before-request="document.getElementById('install-btn').disabled = true;">
    <form>
        <div class="input-group mb-3 ">
            <div class="form-floating {{ if .PathError }} is-invalid {{- end }}">
                <input type="text" class="form-control {{ if .PathError }} is-invalid {{ else if .PathValue }} is-valid {{ end }}" name="pathInput" id="pathInput"
                    hx-post="/validate" 
                    value="{{ .PathValue }}"
                    required>
                <label for="pathInput">Путь до директории миссии</label>
            </div>
            <div class="invalid-feedback">{{ .PathError }}</div>
        </div>
        
        <h5>Компоненты</h5>
        <div>
            <small>Выберите компоненты для установки. Опционально можно указать ссылку на GitHub ветку репозитория из которой хотите получить компонент.</small>
        </div>
        {{range .Components}}
            {{ template "install_component_option" . }}
        {{end}}

        <h5>Опции</h5>
        <div class="mb-3">                    
            <div class="form-check form-switch">
                <input class="form-check-input" type="checkbox" role="switch" name="makeBackup" id="makeBackup" 
                {{ if .BackupEnabled }} checked {{ end }}
                style="cursor: pointer;">
                <label class="form-check-label" for="makeBackup" style="cursor: pointer;">Сохранить настройки</label>
            </div>
            <small>Сохраняет существующие файлы настроек после обновления.</small>
        </div>
        
        <hr>
        <div class="d-grid gap-2 col-12 mb-4 mx-auto">
            <button id="install-btn"  class="btn btn-primary" {{ if not .Verified }} disabled {{ end }}
                    hx-post="/install" hx-target="#main">Установить</button>
        </div>
    </form>
</div>
`

const INSTALL_COMPONENT_OPTION = `
<div class="mb-1">
    <div class="form-check form-switch col-form-label-lg pt-2 pb-0">
        <input class="form-check-input " type="checkbox" role="switch" style="cursor: pointer;" 
                name="{{ .Slug }}CB" id="{{.Slug}}" 
                {{ if .Checked }} checked {{ end }}
                {{ if not .Enabled }} disabled {{ end }}
                hx-post="/validate"
        >
        <label class="form-check-label" for="{{ .Slug }}" style="cursor: pointer;">{{ .Label }}</label>
    </div>
    <small>{{ .Description }}</small>
    {{ if .Checked }}
    <div class="input-group mb-1 {{ if .Error }} is-invalid {{ else }} is-valid {{- end }}">
        <div class="form-floating">
            <input type="text" class="form-control {{ if .Error }} is-invalid {{ else }} is-valid {{- end }}" 
                name="{{ .Slug }}URLInput" id="{{ .Slug }}URLInput" 
                value="{{ .Value }}"
                required
                hx-post="/validate"
                >
            <label for="{{ .Slug }}URLInput">Репозиторий</label>
            <div class="invalid-feedback">{{ .Error }}</div>
        </div>
    </div>
    {{end}}
</div>
`
