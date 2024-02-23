function logout() {
    $.ajax(
        {
            url: "/logout",
            type: "GET",
            success: function (data) {
                if (data.code === 200) {
                    window.location = "/login"
                    return
                } else {
                    return
                }
            }
        },
    );
}


function createUser() {
    let userName = $("#recipient-userName").val();
    let password = $("#recipient-userPasswd").val();
    let password2 = $("#recipient-userPasswd2").val();

    if (userName === "" || password === "" || password2 === "") {
        alert("不允许为空");
        return;
    }
    if (password !== password2) {
        alert("密码输入不一致");
        return;
    }
    $.ajax({
        url: "/user/create",
        type: "POST",
        contentType: 'application/json',
        data: JSON.stringify({
            "userName": userName,
            "password": password2,
        }),
        success: function (data) {
            if (data.code === 200) {
                document.getElementById("recipient-userName").value = "";
                document.getElementById("recipient-userPasswd").value = "";
                document.getElementById("recipient-userPasswd2").value = "";
                alert("添加成功");
                document.getElementById("cclose").click();
                window.location = "/user/index";
            } else {
                alert("添加失败");
            }
        },
        error: function () {
            alert("请求失败！");
        }
    });
}


function deleteUser() {
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
                    url: "/user/delete",
                    type: "DELETE",
                    contentType: 'application/json',
                    data: JSON.stringify({"userName": arr[2]}),
                    success: function (data) {
                        if (data.code !== 200) {
                            window.location = "/user/index";
                            alert("删除失败");
                        } else {
                            trItems.parentNode.removeChild(trItems);
                            window.location = "/user/index";
                            alert("删除成功");
                        }
                    }
                },
            );
        }
    }
}


function delEditInput() {
    document.getElementById("recipient-userName").value = "";
    // 去掉选择的勾勾
    var checkboxes = document.getElementsByClassName('cb');
    for (var i = 0; i < checkboxes.length; i++) {
        checkboxes[i].checked = false;
    }
}

function editUser() {
    let items = document.getElementsByClassName('cb');
    let len = items.length;
    for (var i = len - 1; i >= 0; i--) {
        let is_checkd = items[i].checked;
        if (is_checkd) {
            let trItems = items[i].parentNode.parentNode;
            let trStr = trItems.innerText;
            const arr = trStr.split('\t');
            let updateObj = {
                "id": arr[1],
                "userName": arr[2],
                "newUserDate": arr[3]
            }
            let id = document.getElementById('recipient-userId3');
            let inputName = document.getElementById('recipient-userName3');
            let newUserDate = document.getElementById('recipient-newUserDate3');
            inputName.value = updateObj.userName;
            newUserDate.value = updateObj.newUserDate;
            id.value = updateObj.id;
            newUserDate.setAttribute('readonly', 'true');
            id.setAttribute('readonly', 'true');
        }
    }
}


function updateUser() {
    let userId = $("#recipient-userId3").val();
    let userName = $("#recipient-userName3").val();
    let newUserDate = $("#recipient-newUserDate3").val();
    $.ajax(
        {
            url: "/user/update/info",
            type: "POST",
            contentType: 'application/json',
            data: JSON.stringify({
                "userId": userId,
                "userName": userName,
                "newUserDate": newUserDate
            }),
            success: function (data) {
                if (data.code !== 200) {
                    alert("修改失败")
                } else {
                    document.getElementById("close3").click();
                    window.location = "/user/index"
                }
            }
        },
    );
}


// function updatePasswd() {
//     let userName = $("#recipient-userName2").val();
//     let password = $("#recipient-userPasswds").val();
//     let password2 = $("#recipient-userPasswds2").val();
//     console.log("-----------",password);
//
//     if (userName === "" || password === "" || password2 === ""){
//         alert("不允许为空");
//         return;
//     }
//     if (password !== password2){
//         alert("密码输入不一致");
//         return;
//     }
//     $.ajax({
//         url: "/user/update/pwd",
//         type: "POST",
//         contentType: 'application/json',
//         data: JSON.stringify({
//             "userName": userName,
//             "password": password2,
//         }),
//         success: function (data) {
//             if (data.code === 200) {
//                 document.getElementById("recipient-userName2").value = "";
//                 document.getElementById("recipient-userPasswds").value = "";
//                 document.getElementById("recipient-userPasswds2").value = "";
//                 alert("修改成功");
//                 document.getElementById("uclose").click();
//                 window.location = "/user/index";
//             } else {
//                 alert("修改失败");
//             }
//         },
//         error: function () {
//             alert("请求失败！");
//         }
//     });
// }


function hostReboot() {
    $.ajax(
        {
            url: "/reboot",
            type: "GET",
            success: function (data) {
                if (data.code !== 200) {
                    alert("机器重启失败");
                    document.getElementById("close4").click();
                }
            }
        },
    );
}