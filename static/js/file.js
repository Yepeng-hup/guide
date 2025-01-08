$(function () {
    $("#createdir").click(function () {
        let dirName = $("#recipient-dir").val();
        if (dirName === '') {
            alertMontage("不允许为空.","alert-danger");
            setTimeout(function() {
                window.location = "/";
            }, 1000);
        } else {
            $.post({
                "url": "/file/create",
                "data": {
                    "name": dirName,
                    "path": decodeURIComponent(location.pathname),
                },
                "success": function (data) {
                    if (data["code"] === 200) {
                        window.location = location.pathname;
                    } else {
                        alert(data["message"]);
                    }
                },
                "fail": function (error) {
                    console.log(error);
                }
            });
        }

    });
});


$(function () {
    $("#createfile").click(function () {
        let fileName = $("#recipient-file").val();
        if (fileName === '') {
            alertMontage("不允许为空.","alert-danger");
            setTimeout(function() {
                window.location = "/";
            }, 1000);
        } else {
            $.post({
                "url": "/file/file/create",
                "data": {
                    "name": fileName,
                    "path": decodeURIComponent(location.pathname),
                },
                "success": function (data) {
                    if (data["code"] === 200) {
                        window.location = location.pathname;
                    } else {
                        alert(data["message"]);
                        // alertMontage(data["message"],"alert-danger");
                        // setTimeout(function() {
                            // window.location = location.pathname;
                        // }, 1000);
                    }
                },
                "fail": function (error) {
                    console.log(error);
                }
            });
        }

    });
});


function pushFile() {

    const forms = document.getElementById('sForm');

    forms.addEventListener('submit', function (e) {
        e.preventDefault();
        let htmlData = `
                  <div class="progress" style="width: 70%; margin: 20px auto;">
                    <div class="progress-bar progress-bar-striped" role="progressbar" style="width: 0%;" aria-valuenow="0" aria-valuemin="0" aria-valuemax="100" id="progressBar">0%</div>
                  </div>
            `

        let list = document.getElementById('jdt');
        list.innerHTML = htmlData;

        const progressBar = document.getElementById('progressBar');
        const formData = new FormData(forms);
        formData.append('path', decodeURIComponent(location.pathname));
        const xhr = new XMLHttpRequest();

        xhr.open('POST', '/file/upload', true);

        xhr.upload.addEventListener('progress', function (event) {
            if (event.lengthComputable) {
                const percentComplete = (event.loaded / event.total) * 100;
                progressBar.style.width = percentComplete + '%';
                progressBar.textContent = percentComplete.toFixed(2) + '%';
            }
        });

        xhr.onload = function () {
            if (xhr.status === 200) {
                setTimeout(function () {
                    window.location = decodeURIComponent(location.pathname);
                }, 500);
            } else {
                console.error('push file fail.');
                setTimeout(function () {
                    window.location = decodeURIComponent(location.pathname);
                }, 500);
            }
        };

        xhr.send(formData);
    });

};


function updateContent() {
    let content = document.getElementById('texts').innerText;
    let items=document.getElementsByClassName('cb');
    var checkboxes = document.getElementsByClassName('cb');
    let len=items.length;
    for (var i=len-1; i>=0;i--) {
        let is_checkd = items[i].checked;
        if (is_checkd) {
            let divItems = items[i].parentNode.parentNode;
            let divlr = divItems.innerText;

            $.post({
                "url": "/file/edit",
                "data": {
                    "content": content,
                    "file": divlr,
                    "path": decodeURIComponent(location.pathname)
                },
                "success": function (data) {
                    if (data["code"] === 200) {
                        alert("update success")
                        for (let i = 0; i < checkboxes.length; i++) {
                            checkboxes[i].checked = false;
                        }
                        document.getElementById("down").click();
                    } else {
                        alert("err: update fail")
                        for (let i = 0; i < checkboxes.length; i++) {
                            checkboxes[i].checked = false;
                        }
                    }
                },
                "fail": function (error) {
                    console.log(error);
                }
            });
        }
    }
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
                        "FDpath": decodeURIComponent(location.pathname)
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
    var checkboxes = document.getElementsByClassName('cb');
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
                        "filePath": decodeURIComponent(location.pathname)
                    },
                    "success": function (data) {
                        if (data["code"] === 200) {
                            textpj(data["fileText"], data["fileName"]);
                            return
                        } else {
                            alertMontage("不支持的文件格式.","alert-danger");
                            for (let i = 0; i < checkboxes.length; i++) {
                                checkboxes[i].checked = false;
                            }
                            setTimeout(function() {
                                window.location = decodeURIComponent(location.pathname);
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

function down() {
    var checkboxes = document.getElementsByClassName('cb');
    for (let i = 0; i < checkboxes.length; i++) {
        checkboxes[i].checked = false;
    }
}


function ysCheckbox() {
    let items=document.getElementsByClassName('cb');
    let len=items.length;
    for (var i=len-1; i>=0;i--) {
        let is_checkd = items[i].checked;
        if (is_checkd) {
            let divItems = items[i].parentNode.parentNode;
            let divlr = divItems.innerText;
            $.post(
                {
                    "url": "/file/ys",
                    "data": {
                        "fileName": divlr,
                        "filePath": decodeURIComponent(location.pathname)
                    },
                    "success": function (data) {
                        if (data["code"] === 200) {
                            alertMontage(data["message"],"alert-success");
                            setTimeout(function() {
                                window.location = decodeURIComponent(location.pathname);
                            }, 1000);
                        } else {
                            alertMontage("文件或目录压缩失败.","alert-danger");
                            setTimeout(function() {
                                window.location = decodeURIComponent(location.pathname);
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

function jyCheckbox() {
    let items=document.getElementsByClassName('cb');
    let len=items.length;
    for (var i=len-1; i>=0;i--) {
        let is_checkd = items[i].checked;
        if (is_checkd) {
            let divItems = items[i].parentNode.parentNode;
            let divlr = divItems.innerText;
            $.post(
                {
                    "url": "/file/jy",
                    "data": {
                        "fileName": divlr,
                        "filePath": decodeURIComponent(location.pathname)
                    },
                    "success": function (data) {
                        if (data["code"] === 200) {
                            alertMontage(data["message"],"alert-success");
                            setTimeout(function() {
                                window.location = decodeURIComponent(location.pathname);
                            }, 1000);
                        } else {
                            alertMontage("解压失败.","alert-danger");
                            setTimeout(function() {
                                window.location = decodeURIComponent(location.pathname);
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

function addDelData(data){
    var tbody_rm= document.getElementById('hs-data').querySelector('tbody'); 
    tbody_rm.innerHTML = '';
    const tbody = document.getElementById('hs-data').getElementsByTagName('tbody')[0];
    data.forEach(t => {
        const row = document.createElement('tr');
        const typeNameCell = document.createElement('td');  
        typeNameCell.textContent = t[0];  
        row.appendChild(typeNameCell);

        tbody.appendChild(row);
    })    
}

function showRecycle(){
    $.get(
        {
            "url": "/file/hs",
            "success": function (data) {
                if (data["code"] === 200) {
                    let htmlData = ""
                    for (let i=0; i<data.data.length; i++){
                        htmlData += `<tr>
                        <td><input class="cb" type="checkbox"/></td>
                        <td>${data.data[i]}</td>
                        <td><button class="btn btn-danger" onclick="deleteRecycleFile()">删除</button></td>
                        </tr>`
                    }
                    let h = document.getElementById('hs-data');
                    h.innerHTML=htmlData;
                } else {
                    alert(data["data"]);
                }
            },
            "fail": function (error) {
                console.log(error);
            }
        },
    )
}

function deleteRecycleFile(){
    let items=document.getElementsByClassName('cb');
    let len=items.length;
    for (var i=len-1; i>=0;i--) {
        let is_checkd = items[i].checked;
        if (is_checkd) {
            let divItems = items[i].parentNode.parentNode;
            let divlr = divItems.innerText;
            $.post(
                {
                    "url": "/file/hs/delete",
                    "data": {
                        "fileName": divlr
                    },
                    "success": function (data) {
                        if (data["code"] !== 200) {
                            alert(data.msg)
                        } 
                    },
                },
                

            );
            divItems.parentNode.removeChild(divItems);
        }
    }
}

function unit32ToStrPermissions(octal) {
    // 将数字转换为八进制字符串
    const permissions = (octal).toString(8).padStart(3, '0'); // 确保有三位数

    let result = '';

    for (let i = 0; i < 3; i++) {
        const digit = parseInt(permissions[i], 8); // 将每个八进制数字转换为十进制

        // 计算权限字符
        result += (digit & 4) ? 'r' : '-';
        result += (digit & 2) ? 'w' : '-';
        result += (digit & 1) ? 'x' : '-';
    }

    return result;
}

function fileSearch(){
    let ssText = $("#context").val();

    $.post(
        {
            "url": "/file/ss",
            "data": {
                "fileName": ssText,
                "filePath": decodeURIComponent(location.pathname)
            },
            "success": function (data) {
                let htmlData = "";

                for( var i=0; i<data.fileList.length; i++){
                    const permissionString = unit32ToStrPermissions(data.fileList[i].Power);
                    htmlData += `
                        <tr>
                            <td><input class="cb" type="checkbox"/></td>
                            <td><a href="${data.fileList[i].Href}" style="text-decoration: none;font-size: 18px;"><img src="/sta/img/file.png" style="width: 20px; height: 20px"> ${data.fileList[i].FileName}</a></td>
                            <td>${data.fileList[i].Size}MB</td>
                            <td>${data.fileList[i].Time}</td>
                            <td>${permissionString}</td>
                            <td>
                                <button data-toggle="modal" data-target="#mtk" style="margin-left: 15px" class="btn btn-success" onclick="catFileCheckbox()">查看</button>
                            </td>
                        </tr>
                    `;

                };
                
                let list = document.getElementById('ss');
                list.innerHTML = htmlData;
            },
            "fail": function (error) {
                console.log(error);
            }
        },
    )

}
