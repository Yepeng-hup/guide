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


function deleteCheckbox() {
    let items=document.getElementsByClassName('cb');
    let len=items.length;
    for (var i=len-1; i>=0;i--) {
        let is_checkd = items[i].checked;
        if (is_checkd) {
            let divItems = items[i].parentNode.parentNode;
            let divlr = divItems.innerText;
            $.post(
                {
                    "url": "/file/delete",
                    "data": {
                        "FDname": divlr,
                        "FDpath": location.pathname
                    }
                },

            );
            divItems.parentNode.removeChild(divItems);
        }
    }
}


function catFileCheckbox() {
    let items=document.getElementsByClassName('cb');
    let len=items.length;
    for (var i=len-1; i>=0;i--) {
        let is_checkd = items[i].checked;
        if (is_checkd) {
            let divItems = items[i].parentNode.parentNode;
            let divlr = divItems.innerText;
            $.get(
                {
                    "url": "/file/cat",
                    "data": {
                        "fileName": divlr,
                        "filePath": location.pathname
                    },
                    "success": function (data) {
                        if (data["code"] === 200) {
                            textpj(data["fileText"]);
                            return
                        } else {
                            alertMontage("这不是文件格式.","alert-danger");
                            setTimeout(function() {
                                window.location = location.pathname;
                            }, 1000);
                        }
                    },
                    "fail": function (error) {
                        console.log(error);
                    }
                },
            )
        }
    }
}




