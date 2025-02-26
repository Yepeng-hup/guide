function deleteUP() {
    let items=document.getElementsByClassName('cb');
    let len=items.length;
    for (var i=len-1; i>=0;i--) {
        let is_checkd = items[i].checked;
        if (is_checkd) {
            let divItems = items[i].parentNode.parentNode;
            let divlr = divItems.innerText;
            $.post(
                {
                    "url": "/pwd/delete",
                    "data": {
                        "pwd": divlr,
                    }
                },

            );
            divItems.parentNode.removeChild(divItems);
        }
    }
}


function catdPwd() {
    let items=document.getElementsByClassName('cb');
    let len=items.length;
    for (var i=len-1; i>=0;i--) {
        let is_checkd = items[i].checked;
        if (is_checkd) {
            let divItems = items[i].parentNode.parentNode;
            let divlr = divItems.innerText;
            $.post(
                {
                    "url": "/pwd/cat",
                    "data": {
                        "pwd": divlr,
                    },
                    "success": function (data) {
                        if (data["code"] === 200) {
                            passwdpj(data["dPassword"]);
                            return
                        } else {
                            window.location = location.pathname;
                        }
                    },
                    "fail": function (error) {
                        console.log(error);
                    }

                },

            );
        }
    }
}


function backupData() {
            $.get({
                "url": "/pwd/bak",
                "success": function (data) {
                    if (data["code"] === 200) {
                            alert(data["msg"]);
                    } else {
                            alert(data["msg"]);
                    }
                },
                "fail": function (error) {
                    console.log(error);
                }
            });
}