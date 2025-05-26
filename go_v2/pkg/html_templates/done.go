package html_templates

const DONE = `
<div class="alert alert-success" role="alert">
    <h4 class="alert-heading">Восхитительный успех!</h4>
    <p>О да, Вы успешно установили все нобходимые скрипты для создания миссии. Успехов в вашем творчестве!</p>
    <hr>
    <p><b>Директория миссии:</b></p>
    <div class="card" onClick="navigator.clipboard.writeText(document.getElementById('toCopy').textContent);" style="cursor: pointer">
        <div class="card-body">
            <div id="toCopy" class="card-text font-monospace">{{ . }}</div>
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
`
