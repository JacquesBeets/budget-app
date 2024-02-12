document.body.addEventListener('htmx:afterOnLoad', function (event) {
    var url = "/";
    switch (event.detail.pathInfo.requestPath) {
        case '/templates/dashboard':
            document.title = 'BeetsDeBeer - Dashboard';
            break;
        case '/templates/upload':
            document.title = 'BeetsDeBeer - Upload OFX Data';
            url = '/upload';
            break;
        default:
            document.title = 'BeetsDeBeer';
            break;
    }
    history.pushState({}, '', url);
});