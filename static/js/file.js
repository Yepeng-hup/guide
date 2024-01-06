$(function () {
    $("#createdir").click(function () {
        let dirName = $("#dir").val();
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
                    "path": location.pathname
                },
                "success": function (data) {
                    if (data["code"] === 200) {
                        alertMontage(data["message"],"alert-success");
                        setTimeout(function() {
                            window.location = location.pathname;
                        }, 1000);
                    } else {
                        alertMontage(data["message"],"alert-danger");
                        setTimeout(function() {
                            window.location = location.pathname;
                        }, 1000);
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
        let fileName = $("#file").val();
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
                    "path": location.pathname
                },
                "success": function (data) {
                    if (data["code"] === 200) {
                        alertMontage(data["message"],"alert-success");
                        setTimeout(function() {
                            window.location = location.pathname;
                        }, 1000);
                    } else {
                        alertMontage(data["message"],"alert-danger");
                        setTimeout(function() {
                            window.location = location.pathname;
                        }, 1000);
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
        formData.append('path', location.pathname);
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
                    window.location = location.pathname;
                }, 500);
            } else {
                console.error('push file fail.');
                setTimeout(function () {
                    window.location = location.pathname;
                }, 500);
            }
        };

        xhr.send(formData);
    });

};


function updateContent() {
    let content = document.getElementById('texts').innerText;
    let items=document.getElementsByClassName('cb');
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
                    "path": location.pathname
                },
                "success": function (data) {
                    if (data["code"] === 200) {
                        alert("update success")
                    } else {
                        alert("err: update fail")
                    }
                },
                "fail": function (error) {
                    console.log(error);
                }
            });
        }
    }
}