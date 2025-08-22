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
    let htmlTitle = `
        <h5 class="modal-title">${fileName}</h5>
    `;

    let title = document.getElementById('fileName');
    title.innerHTML = htmlTitle;
}

function SStextpj(fileName) {
    let htmlTitle = `
        <h5 class="modal-title">${fileName}</h5>
    `;

    let title = document.getElementById('fileName1');
    title.innerHTML = htmlTitle;
}

function passwdpj(pwd){
    let htmlData = `
        <pre style="color: white; padding: 10px">${pwd}</pre>
    `;

    let list = document.getElementById('pwd-txt');
    list.innerHTML = htmlData;
}
