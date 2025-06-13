function addBlacklist() {
    let items = document.getElementsByClassName('cb');
    let len = items.length;
    for (var i = len - 1; i >= 0; i--) {
        let is_checkd = items[i].checked;
        if (is_checkd) {
            let trItems = items[i].parentNode.parentNode;
            let trStr = trItems.innerText;
            const arr = trStr.split('\t');
            $.ajax(
                {
                    url: "/security/black/add",
                    type: "POST",
                    contentType: 'application/json',
                    data: JSON.stringify({"blacklistIp": arr[1]}),
                    success: function (data) {
                        if (data.code !== 200) {
                            window.location.href = window.location;
                            alert(data["msg"]);
                        } else {
                            window.location.href = window.location;
                            alert(data["msg"]);
                        }
                    }
                },
            );
        }
    }
}


function moveBlacklist(){
    let items = document.getElementsByClassName('cb');
    let len = items.length;
    for (var i = len - 1; i >= 0; i--) {
        let is_checkd = items[i].checked;
        if (is_checkd) {
            let trItems = items[i].parentNode.parentNode;
            let trStr = trItems.innerText;
            const arr = trStr.split('\t');
            $.ajax(
                {
                    url: "/security/black/mv",
                    type: "DELETE",
                    contentType: 'application/json',
                    data: JSON.stringify({"blacklistIp": arr[1]}),
                    success: function (data) {
                        if (data.code !== 200) {
                            window.location.href = window.location;
                            alert(data["msg"]);
                        } else {
                            window.location.href = window.location;
                            alert(data["msg"]);
                        }
                    }
                },
            );
        }
    }
}


function showBlacklist(){
    $.get(
        {
            "url": "/security/black/show",
            "success": function (data) {
                let htmlData = "";

                for( var i=0; i<data.blackList.length; i++){
                    htmlData += `
                            <tr>
                                <td>${data.blackList[i].NewAddDate}</td>
                                <td style="color:rgb(231, 11, 11)">${data.blackList[i].Ip}</td>
                            </tr>
                    `;

                };

                let list = document.getElementById('blackList');
                list.innerHTML = htmlData;
            },
            "fail": function (error) {
                console.log(error);
            }
        },
    )
}