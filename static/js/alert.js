function alertMontage(info, alertClass){
    let htmlData = `
        <div class="container mt-5">
            <div id="myAlert" class="alert ${alertClass}" role="alert">
                ${info}
            </div>
        </div>
    `;

    let list = document.getElementById('alert');
    list.innerHTML = htmlData;
}


function textpj(text) {
    let htmlData = `
        <pre style="color: white">${text}</pre>
    `;

    let list = document.getElementById('texts');
    list.innerHTML = htmlData;
}