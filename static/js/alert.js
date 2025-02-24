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


function textpj(fileName) {
    // let htmlData = `
    //     <pre style="color: white; padding: 10px">${text}</pre>
    // `;
    let htmlTitle = `
        <h5 class="modal-title">${fileName}</h5>
    `;

    // let list = document.getElementById('texts');
    // list.innerHTML = htmlData;
    let title = document.getElementById('fileName');
    title.innerHTML = htmlTitle;
}