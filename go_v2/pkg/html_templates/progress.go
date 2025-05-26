package html_templates

const PROGRESS = `
<!-- Progress form -->
<div class="col-md-9 offset-md-1">
  <div class="d-flex justify-content-center" hx-trigger="every 1s" hx-get="/status" hx-target="#main">
    <div class="spinner-border" role="status">
      <span class="visually-hidden">Загрузка...</span>
    </div>
  </div>

  <ul class="list-unstyled" id="logs">
    {{range .}}
      <li><span class="badge bg-{{.Badge}}">{{ .Tag }}</span> {{ .Message }}</li>
    {{end}}
  </ul>
</div>
`
