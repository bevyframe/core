const get = (context) => {
    return Page({
        title: "BevyFrame Test App",
        themeColor: "blank",
    }, [
        Navbar([
            NavItem("home", "/", "Home", true),
            NavItem("apps", "/env.html", "Demo"),
        ]),
        Root({
            margin: Margin({
                left: Size.pixel(100)
            })
        }, [
            Container({
                id: 'info',
                margin: Margin({
                    bottom: Size.pixel(10)
                }),
            }, [
                Button('mini', {
                    onclick: "load_info()",
                    innertext: 'Load Info'
                })
            ]),
            Box({
                width: Size.max_content,
                text_align: Align.center,
            }, [
                Line([ Textbox('', {type: "text", placeholder: 'textbox', value: context.ip}) ]),
                Line([ Button({innertext: 'Button'}) ]),
                Line([ Button('', {selector: 'small', innertext: 'Button'}) ]),
                Line([ Button('', {selector: 'mini', innertext: 'Button'}) ])
            ])
        ])
    ])
}