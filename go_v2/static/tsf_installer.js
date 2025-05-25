const App = new (function () {
    this.page = 0;
})();
var AX = 0;

document.addEventListener('DOMContentLoaded', () => {
    document.body.addEventListener('htmx:beforeSwap', function(evt) {
        if(evt.detail.xhr.status === 404){
            console.error("Error: Could Not Find Resource");
        } else if(evt.detail.xhr.status === 422){
            evt.detail.shouldSwap = true;
            evt.detail.isError = true;
        }
    });
});