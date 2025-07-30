const TheProtocols = {
    getCurrentUser: async (context) => fetch(
        `https://${window.location.host}/.well-known/theprotocols`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                network: context.network,
                endpoint: "current_user_info",
                body: JSON.stringify({
                    cred: context.token
                })
            })
        }
    )
        .then(response => response.json())
        .catch(_ => {}),
};

const buildContext = async (inp) => {
    let context = {
        headers: {},
        user: {
            id: {}
        }
    };
    for (let line of inp.split('\n')) {
        let key = line.split(': ', 2)[0];
        let val = line.split(': ', 2)[1];
        if (val === undefined)
            val = '';
        switch (key) {
            case "Package":
                context.package = val;
                break
            case "Cred.Email":
                context.email = val;
                break
            case "Cred.Username":
                context.username = val;
                break
            case "Cred.Network":
                context.username = val;
                break
            case "Cred.Password":
                context.password = val;
                break
            case "Cred.Token":
                context.token = val;
                break
            case "Path":
                context.path = val;
                break
            case "Method":
                context.method = val;
                break
            case "IP":
                context.ip = val;
                break
            case "Permissions":
                context.permissions = val.split(',');
                break
            case "LoginView":
                context.loginView = val;
                break
            default:
                if (key.split('.')[0] === 'Header')
                    context.headers[key.split('.')[1]] = val;
                break
        }
    }
    if (!context.network)
        context.network = context.email.split('@')[1];
    context.user.id = await TheProtocols.getCurrentUser(context);
    document.context = context;
    return context;
};