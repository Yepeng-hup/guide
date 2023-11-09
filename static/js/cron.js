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
                    "url": "/cron/delete",
                    "data": {
                        "cron": divlr,
                    }
                },

            );
            divItems.parentNode.removeChild(divItems);
        }
    }
}


function deleteSvc() {
    let items=document.getElementsByClassName('cb');
    let len=items.length;
    for (var i=len-1; i>=0;i--) {
        let is_checkd = items[i].checked;
        if (is_checkd) {
            let divItems = items[i].parentNode.parentNode;
            let divlr = divItems.innerText;
            $.post(
                {
                    "url": "/svc/delete",
                    "data": {
                        "svc": divlr,
                    }
                },

            );
            divItems.parentNode.removeChild(divItems);
        }
    }
}