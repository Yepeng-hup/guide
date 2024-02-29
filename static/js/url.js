let temporaryStorage = [];

function editUrlInfo() {
    let items = document.getElementsByClassName('cb');
    let len = items.length;
    for (var i = len - 1; i >= 0; i--) {
        let is_checkd = items[i].checked;
        if (is_checkd) {
            let trItems = items[i].parentNode.parentNode;
            let trStr = trItems.innerText;
            const arr = trStr.split('\t');
            let updateObj = {
                "urlName": arr[1],
                "url": arr[2],
                "notes": arr[3]
            }
            let urlName = document.getElementById('recipient-urlName');
            let url = document.getElementById('recipient-url');
            let notes = document.getElementById('recipient-notes');
            urlName.value = updateObj.urlName;
            url.value = updateObj.url;
            notes.value = updateObj.notes;
            temporaryStorage.push(updateObj.urlName);
            temporaryStorage.push(updateObj.url);
            temporaryStorage.push(updateObj.notes);
        }
    }
}

function delUrlEditInput() {
    document.getElementById("recipient-urlName").value = "";
    document.getElementById("recipient-url").value = "";
    document.getElementById("recipient-notes").value = "";
    // 去掉选择的勾勾
    var checkboxes = document.getElementsByClassName('cb');
    for (var i = 0; i < checkboxes.length; i++) {
        checkboxes[i].checked = false;
    }
}

function updateUrlInfo() {
    let urlName = $("#recipient-urlName").val();
    let url = $("#recipient-url").val();
    let notes = $("#recipient-notes").val();
    $.ajax(
        {
            url: "/url/update",
            type: "POST",
            contentType: 'application/json',
            data: JSON.stringify({
                "srcUrlName": temporaryStorage[0],
                "srcUrl": temporaryStorage[1],
                "srcNotes": temporaryStorage[2],
                "urlName": urlName,
                "url": url,
                "notes": notes
            }),
            success: function (data) {
                if (data.code !== 200) {
                    temporaryStorage.length = 0;
                    alert(data.message);
                } else {
                    temporaryStorage.length = 0;
                    document.getElementById("close").click();
                    window.location = "/url/index"
                }
            }
        },
    );
}