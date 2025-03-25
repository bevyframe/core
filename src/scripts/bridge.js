const _bridge = async (func, ...args) => {
    let resp = await fetch(`${location.protocol}//${location.host}/.well-known/bevyframe/proxy`,{
        method:'POST',
        headers: {
            'Content-Type':'application/json'
        },
        body: JSON.stringify({
            func: func,
            args: JSON.stringify(args),
            path: window.location.pathname,
            cookie: document.cookie
        })
    });
    let data = await resp.json();
    if (data.error)
        throw new Error(data.error);
    else if (data.type === 'return')
        return data.value;
    else if (data.type === 'script')
        eval(data.value);
    else if (data.type === 'view') {
        const target = data.element === 'body' ?
            document.body :
            document.querySelector(data.element);
        target.innerHTML = '';
        if (Array.isArray(data.value))
            for (let element of data.value)
                renderWidget(element, target);
        else
            target.innerHTML = data.value;
    }
};