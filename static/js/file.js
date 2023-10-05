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
