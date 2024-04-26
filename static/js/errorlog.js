function getErrorLog(){
    $.ajax({
        url: "/log/r",
        type: "POST",
        contentType: 'application/json',
        data: JSON.stringify({
            "logType": "error",
        }),
        success: function (data) {
            if (data.code === 200) {
                let html= '';
                for (const item of data.logRel) {
                    html += '<tr>';
                    // html += '<td><input class="cb" type="checkbox" id="cb"/></td>'
                    // html += '<td>' + item.ID + '</td>';
                    // html += '<td>' + item.Date + '</td>';
                    html += '<td>' + item.LogType + '</td>';
                    html += '<td><pre>' + item.LogContent + '</pre></td>';
                    html += '</tr>';
                }
                document.getElementById('data').innerHTML = html;
            } else {
                console.log("add log fail.");
            }
        },
        error: function () {
            console.log("get log request fail.");
        }
    });
}


function getWarnLog(){
    $.ajax({
        url: "/log/r",
        type: "POST",
        contentType: 'application/json',
        data: JSON.stringify({
            "logType": "warn",
        }),
        success: function (data) {
            if (data.code === 200) {
                let html= '';
                for (const item of data.logRel) {
                    html += '<tr>';
                    html += '<td>' + item.LogType + '</td>';
                    html += '<td><pre>' + item.LogContent + '</pre></td>';
                    html += '</tr>';
                }
                document.getElementById('data').innerHTML = html;
            } else {
                console.log("add log fail.");
            }
        },
        error: function () {
            console.log("get log request fail.");
        }
    });
}

function deleteLogLimit(){
    let logNum = 50

    $.ajax({
        url: "/log/d",
        type: "POST",
        contentType: 'application/json',
        data: JSON.stringify({
            "logNum": logNum,
        }),
        success: function (data) {
            if (data.code === 200) {
                console.log("ok")
            } else {
                console.log("delete log fail.");
            }
        },
        error: function () {
            console.log("delete log request fail.");
        }
    });
}