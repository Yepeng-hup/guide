function alertMontage(info, alertClass){
    let htmlData = `      
        <div id="myAlert" class="alert ${alertClass} alert-dismissible fade show" role="alert">
          <strong>${info}</strong>
          <button type="button" class="close" data-dismiss="alert" aria-label="Close">
            <span aria-hidden="true">&times;</span>
          </button>
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