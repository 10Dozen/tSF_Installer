package html_templates

const DONE = `
<div class="col-md-9 offset-md-1">
    <div class="alert alert-success" role="alert">
        <h4 class="alert-heading">Восхитительный успех!</h4>
        <p>О да, Вы успешно установили все нобходимые скрипты для создания миссии. Успехов в вашем творчестве!</p>
        <hr>
        <p><b>Директория миссии:</b></p>
        <div class="card" onClick="navigator.clipboard.writeText(document.getElementById('toCopy').textContent);" style="cursor: pointer">
            <div class="card-body">
                <div id="toCopy" class="card-text font-monospace">{{ .TargetDir }}</div>
            </div>
        </div>
        <hr>
        <p class="mb-0">Полезные материалы:</p>
        <ul>
            <li><a href="https://tacticalshift.ru/docs/MMO/editor_requirements.html">Требования к оформлению миссии</a></li>
            <li><a href="https://tacticalshift.ru/docs/MMO/editor_framework.html">Обзор и знакомство Tactical Shift Mission Framework</a></li>
            <li><a href="https://tacticalshift.ru/docs/MMO/editor_concept.html">Концепт и структура миссии</a></li>
            <li><a href="https://tacticalshift.ru/docs/MMO/editor_production.html">Сборка миссии</a></li>
            <li><a href="https://tacticalshift.ru/docs/MMO/editor_misc.html">Ревью и поддержка миссии</a></li>
        </ul>
    </div>

    <div class="accordion accordion-flush border border-primary-subtle" id="accordionFlushExample">
        <div class="accordion-item">
            <h2 class="accordion-header">
                <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse"
                data-bs-target="#flush-collapseOne" aria-expanded="false" aria-controls="flush-collapseOne">
                    Журнал установки
                </button>
            </h2>
            <div id="flush-collapseOne" class="accordion-collapse collapse" data-bs-parent="#accordionFlushExample">
                <div class="accordion-body">
                    <ul class="list-unstyled" id="logs">
                        {{range .Logs}}
                            <li><span class="badge bg-{{.Badge}}">{{ .Tag }}</span> {{ .Message }}</li>
                        {{end}}
                    </ul>
                </div>
            </div>
        </div>
    </div>
</div>
`
