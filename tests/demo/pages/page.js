const get = (context) => {
    return Page({
        title: "BevyFrame Test App",
        color: context.user.id.settings.theme_color,
    }, [
        Navbar([
            NavItem("home", "/", "Home", true),
            NavItem("apps", "/env.html", "Demo"),
        ]),
        Root({
            margin: Margin({
                left: Size.pixel(80)
            })
        }, [
            Container({
                id: 'info',
                margin: Margin({
                    bottom: Size.pixel(10)
                }),
            }, [
                MiniButton({
                    onClick: "load_info()",
                    innerText: 'Load Info'
                })
            ]),
            Box({
                width: Size.maxContent,
                textAlign: Align.center,
            }, [
                Line([ Textbox('', {type: "text", placeholder: 'textbox', value: context.ip}) ]),
                Line([ Button({innerText: 'Button'}) ]),
                Line([ SmallButton({innerText: 'Button'}) ]),
                Line([ MiniButton({innerText: 'Button'}) ])
            ])
        ])
    ])
}