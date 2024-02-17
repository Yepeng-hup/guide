function logout() {
    $.ajax(
        {
            url: "/logout",
            type: "GET",
            // contentType: 'application/json',
            // data: JSON.stringify({"id": 1}),
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