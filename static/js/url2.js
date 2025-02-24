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
                "type": arr[3],
                "notes": arr[4],
            }
            let urlName = document.getElementById('recipient-urlName');
            let url = document.getElementById('recipient-url');
            let notes = document.getElementById('recipient-notes');
            let type = document.getElementById('recipient-type');
            urlName.value = updateObj.urlName;
            url.value = updateObj.url;
            type.value = updateObj.type;
            notes.value = updateObj.notes;
            temporaryStorage.push(updateObj.urlName);
            temporaryStorage.push(updateObj.url);
            temporaryStorage.push(updateObj.type);
            temporaryStorage.push(updateObj.notes);
            urlName.setAttribute('readonly', 'true');
            type.setAttribute('readonly', 'true');
        }
    }
}

function delUrlEditInput() {
    document.getElementById("recipient-urlName").value = "";
    document.getElementById("recipient-url").value = "";
    document.getElementById("recipient-notes").value = "";
    document.getElementById("recipient-type").value = "";
    // 去掉选择的勾勾
    var checkboxes = document.getElementsByClassName('cb');
    for (var i = 0; i < checkboxes.length; i++) {
        checkboxes[i].checked = false;
    }
}

function updateUrlInfo() {
    let urlName = $("#recipient-urlName").val();
    let url = $("#recipient-url").val();
    let type = $("#recipient-type").val();
    let notes = $("#recipient-notes").val();
    $.ajax(
        {
            url: "/url/update",
            type: "POST",
            contentType: 'application/json',
            data: JSON.stringify({
                "urlName": urlName,
                "url": url,
                "type": type,
                "notes": notes
            }),
            success: function (data) {
                if (data.code !== 200) {
                    temporaryStorage.length = 0;
                    alert(data.data);
                } else {
                    temporaryStorage.length = 0;
                    document.getElementById("close").click();
                    window.location = location.pathname;
                }
            }
        },
    );
}


$(function () {
    $("#createtype").click(function () {
        let typeName = $("#recipient-file").val();
        if (typeName === '') {
            alert("不允许为空");
        } else {
            $.post({
                "url": "/url/type/create",
                "data": {
                    "type-name": typeName
                },
                success: function (data) {
                    if (data["code"] === 200) {
                        window.location = location.pathname;
                    } else {
                        alert(data["message"]);
                        window.location = location.pathname;
                    }
                },
                fail: function (error) {
                    console.log(error);
                }
            });
        }

    });
});


function updateWebUrlTypeData(data){
    var tbody_rm= document.getElementById('type-data').querySelector('tbody'); 
    tbody_rm.innerHTML = '';
    const tbody = document.getElementById('type-data').getElementsByTagName('tbody')[0];
    data.forEach(t => {
        const row = document.createElement('tr');
        const useCell = document.createElement('td');  
        useCell.innerHTML = `
            <input class="cb" type="checkbox"/>
        `;  
        row.appendChild(useCell);

        const typeNameCell = document.createElement('td');  
        typeNameCell.textContent = t.TypeName;  
        row.appendChild(typeNameCell);

        const linkButtonCell = document.createElement('td');
        linkButtonCell.innerHTML = `
            <button class="btn btn-danger" onclick="deleteTypeCheckbox()">删除</button>
        `;
        row.appendChild(linkButtonCell);

        tbody.appendChild(row);
    })    
}


function showUrlType() {
    $.ajax(
        {
            url: "/url/type/list",
            type: "GET",
            success: function (data) {
                if (data.code !== 200) {
                    alert(data.data);
                } else {
                    updateWebUrlTypeData(data.data);
                }
            },
            fail: function (error) {
                console.log(error);
            }
        }
    );
}


function deleteTypeCheckbox(){
    var items=document.getElementsByClassName('cb');
    var len=items.length;
    for (var i=len-1; i>=0;i--) {
        var is_checkd = items[i].checked;
        if (is_checkd) {
            var divItems = items[i].parentNode.parentNode;
            var typeName = divItems.innerText;
            $.ajax(
                {
                    url: "/url/type/del",
                    type: "POST",
                    data: {
                        "type-name": typeName,
                    },
                    success: function (data) {
                        if (data.code !== 200) {
                            alert(data.data);
                        } else {
                            divItems.parentNode.removeChild(divItems);
                            window.location = location.pathname;
                        }
                    },
                    fail: function (error) {
                        console.log(error);
                    }
                },

            );
        }
    }
}


function updateWebUrlData(data){
        var tbody_rm= document.querySelector('tbody'); 
        tbody_rm.innerHTML = '';
        const tbody = document.getElementById('url-data').getElementsByTagName('tbody')[0];
        data.forEach(d => {
            const row = document.createElement('tr');

            // 创建并添加选择单元格  
            const useCell = document.createElement('td');  
            useCell.innerHTML = `
                <input class="cb" type="checkbox"/>
            `;  
            row.appendChild(useCell); 

            // 创建并添加名称单元格  
            const urlNameCell = document.createElement('td');  
            urlNameCell.textContent = d.UrlName;  
            row.appendChild(urlNameCell);  
  
            // 创建并添加网址单元格  
            const urlAddressCell = document.createElement('td');  
            urlAddressCell.textContent = d.UrlAddress;  
            row.appendChild(urlAddressCell);  
  
            // 创建并添加类名单元格  
            const urlTypeCell = document.createElement('td');  
            urlTypeCell.textContent = d.UrlType;  
            row.appendChild(urlTypeCell);  
  
            // 创建并添加备注单元格  
            const urlNotesCell = document.createElement('td');  
            urlNotesCell.textContent = d.UrlNotes;  
            row.appendChild(urlNotesCell);

            // 创建添加操作单元格
            const linkButtonCell = document.createElement('td');
            linkButtonCell.innerHTML = `
                <a href="${d.UrlAddress}" target=_blank><input class="btn btn-success" type="button" value="访问"></a>
                <input  class="btn btn-primary" type="button" value="编辑" data-toggle="modal" data-target="#exampleModal" data-whatever="@mdo" onclick="editUrlInfo()">
            `;
            row.appendChild(linkButtonCell);

            // 将行添加到tbody中  
            tbody.appendChild(row);
        });

}

function getTypeUrlData(type){
    $.ajax(
        {
            url: "/url/show",
            type: "POST",
            contentType: 'application/json',
            data: JSON.stringify({
                "urltype": type
            }),
            success: function (data) {
                if (data.code !== 200) {
                    alert(data.data);
                    window.location = location.pathname;;
                }
                const urlData = data.data;
                updateWebUrlData(urlData);
            },

            fail: function (error) {
                    console.log(error);
                }
        },
    );
}