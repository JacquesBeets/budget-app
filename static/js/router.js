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

window.QUERY_PARAMS = {
    searchURL: () => {
        return new URLSearchParams(window.location.search);
    },
    all: () => {
        const urlSearchParams = QUERY_PARAMS.searchURL();
        const params = Object.fromEntries(urlSearchParams.entries());
        return params;
    },
    allNames: () => {
        const urlSearchParams = QUERY_PARAMS.searchURL();
        const params = urlSearchParams.keys();
        return params;
    },
    get: (param) => {
        const urlSearchParams = QUERY_PARAMS.searchURL();
        return urlSearchParams.get(param) || null;
    },
    set: (param, value) => {
        const url = new URL(window.location.href);
        const urlSearchParams = QUERY_PARAMS.searchURL();
        const params = {
            ...QUERY_PARAMS.all(),
            [param]: value
        }

        for(const key in params){
            urlSearchParams.set(key, params[key]);
        }

        url.search = urlSearchParams.toString();
        history.pushState(url, document.title, url);
    },
    remove: (param) => {
        const url = new URL(window.location.href);
        const urlSearchParams = QUERY_PARAMS.searchURL();
        urlSearchParams.delete(param);
        url.search = urlSearchParams.toString();
        history.pushState(url, document.title, url);
    }
};

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


    document.body.addEventListener('htmx:beforeRequest', (event) => {
        // Set the title and push the state
        const requestPath = event.detail.pathInfo.requestPath;
        
        // Set Page Title
        if(CURRENT_ROUTES[requestPath]) setTitle(CURRENT_ROUTES[requestPath].title);
    
        // Update the URL and query params
        const url = CURRENT_ROUTES[requestPath]?.url && CURRENT_ROUTES[requestPath]?.url !== '/' ? new URL(`${window.location.origin}${CURRENT_ROUTES[requestPath]?.url}`) : new URL(window.location.origin);
        const urlSearchParams = QUERY_PARAMS.searchURL();
        const params = {
            ...QUERY_PARAMS.all(),
            ...event.detail.requestConfig.parameters
        }

        for(const key in params){
            urlSearchParams.set(key, params[key]);
        }

        url.search = urlSearchParams.toString();

        history.replaceState(
            CURRENT_ROUTES[requestPath], 
            CURRENT_ROUTES[requestPath]?.title ?? "BeetsDeBeer - Dashboard", 
            url
        );
    });

}
