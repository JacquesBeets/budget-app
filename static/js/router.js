const CURRENT_ROUTES = {
    '/components/transactions':{
        title: 'BeetsDeBeer - Dashboard',
        url: '/',
        key: '/components/transactions'
    },
    '/templates/transactions':{
        title: 'BeetsDeBeer - Transactions',
        url: '/transactions',
        key: '/templates/transactions'
    },
    '/templates/crypto':{
        title: 'BeetsDeBeer - Crypto',
        url: '/crypto',
        key: '/templates/crypto'
    },
    '/templates/transactionstypes':{
        title: 'BeetsDeBeer - Transaction Types',
        url: '/transactionstypes',
        key: '/templates/transactionstypes'
    },
    '/templates/upload':{
        title: 'BeetsDeBeer - Upload OFX Data',
        url: '/upload',
        key: '/templates/upload'
    }
}

const REVERSE_ROUTES = () =>{
    let reverse = {};
    for (const key in CURRENT_ROUTES) {
        reverse[CURRENT_ROUTES[key].url] = CURRENT_ROUTES[key];
    }
    return reverse;
}

const ROUTES_REVERSED = REVERSE_ROUTES();

const setTitle = (title) => {
    document.title = title;
}

window.onload = () => {
    if(ROUTES_REVERSED[window.location.pathname] && window.location.pathname !== '/'){
        const htmxMethods = document.querySelectorAll('[hx-trigger="click"]');
        htmxMethods.forEach((method) => {
            if(ROUTES_REVERSED[window.location.pathname].key === method.attributes['hx-get'].value){
                method.click();
                setTitle(CURRENT_ROUTES[ROUTES_REVERSED[window.location.pathname].key].title);
                history.replaceState(
                    CURRENT_ROUTES[ROUTES_REVERSED[window.location.pathname].key], 
                    CURRENT_ROUTES[ROUTES_REVERSED[window.location.pathname].key].title, 
                    CURRENT_ROUTES[ROUTES_REVERSED[window.location.pathname].key].url
                );
            }
        });
    }
    
    
    document.body.addEventListener('htmx:afterOnLoad', (event) => {
        let requestPath = event.detail.pathInfo.requestPath;
        if(!CURRENT_ROUTES[requestPath]) return;

        setTitle(CURRENT_ROUTES[requestPath].title);
        history.pushState(
            CURRENT_ROUTES[requestPath], 
            CURRENT_ROUTES[requestPath].title, 
            CURRENT_ROUTES[requestPath].url
        );
    });
}
