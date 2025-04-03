function removeHtml(){
    const htmlPerID = ["user", "url", "file", "service", "passwd", "log", "security", "cron", "other"]
    const storedData = localStorage.getItem("userPermissions");
    const permissions = storedData ? JSON.parse(storedData) : [];
    const result = htmlPerID.filter(item => !permissions.includes(item));

    document.addEventListener('DOMContentLoaded', () => {
        result.forEach((item) => {
            const element = document.getElementById(item);
            if (element) {
                element.remove();
            } else {
                console.warn(`${item} nil`);
            }
        });
    });
}
removeHtml()