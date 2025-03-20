function createRole(){
    let roleName = $("#recipient-roleName").val();
    let items = document.getElementsByClassName('cb');
    let len = items.length;
    let permissionList = [];

    if (roleName === "") {
        alert("不允许为空");
        return;
    }

    for (var i = len - 1; i >= 0; i--) {
        let is_checkd = items[i].checked;
        if (is_checkd) {
            let trItems = items[i].parentNode.parentNode;
            let trStr = trItems.innerText;
            const arr = trStr.split('\t');
            const arr1 = arr[0].split(':')
            permissionList.push(arr1[1])
        }
    }

    $.ajax(
        {
            url: "/user/role/create",
            type: "POST",
            contentType: 'application/json',
            data: JSON.stringify({
                "roleName": roleName,
                "permission": permissionList
            }),
            success: function (data) {
                if (data.code !== 200) {
                    alert(data["msg"]);
                    return;
                } else {
                    window.location = "/user/role/index";
                }
            }
        },
    );

}

function selectRole(){
    $.get(
        {
            "url": "/user/role/select",
            "success": function (data) {
                let htmlData = "";

                for( var i=0; i<data.roleList.length; i++){
                    htmlData += `
                            <option>${data.roleList[i].RoleName}</option>
                    `;

                };

                let list = document.getElementById('roles');
                list.innerHTML = htmlData;
            },
            "fail": function (error) {
                console.log(error);
            }
        },
    )
}

function selectRolePermission(){
    let sRole = $("#roles").val();
    $.post(
        {
            "url": "/user/role/per/select",
            "data": {
                "roleName": sRole
            },
            "success": function (data) {
                let htmlData = "";

                for( var i=0; i<data.permission.length; i++){
                    htmlData += `
                             <tr>
                                <td><input class="cb" type="checkbox"/></td>
                                <td>${data.permission[i].RoleName}</td>
                                <td style="color: #0b87e7">${data.permission[i].Permission}</td>
                            </tr>
                    `;

                };

                let list = document.getElementById('per');
                list.innerHTML = htmlData;
            },
            "fail": function (error) {
                console.log(error);
            }
        },
    )
}