$(function () {
    $("#createdir").click(function () {
        let dirName = $("#dir").val();
        if (dirName === '') {
            alertMontage("不允许为空.","alert-danger");
            setTimeout(function() {
                window.location = "/";
            }, 2000);
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
                        }, 2000);
                    } else {
                        alertMontage(data["message"],"alert-danger");
                        setTimeout(function() {
                            window.location = location.pathname;
                        }, 2000);
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
