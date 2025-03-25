from bevyframe import *
from widgets.CustomButton import CustomButton

blacklist = lambda: ["demo@hereus.net"]


@login_required
def get(context: Context) -> (Page, Response):
    # last_ip_query = context.db.query(Test).filter_by(email=context.email).all()
    # if len(last_ip_query) == 0:
    #     raise Error404
    # last_ip = last_ip_query[-1].ip
    return Page(
        title='BevyFrame Test App',
        description='BevyFrame Test App',
        color=context.user.id.settings.theme_color,
        childs=[
            Navbar([
                NavItem('home', '/', 'Home', True),
                NavItem('apps', '/env.html', 'Demo'),
            ]),
            Root([
                Container(
                    id='info',
                    margin=Margin(bottom=Size.pixel(10)),
                    childs=[
                        Button(
                            'mini',
                            onclick=context.execute.load_info(),
                            innertext='Load Info'
                        )
                    ]
                ),
                Box(
                    width=Size.max_content,
                    text_align=Align.center,
                    childs=[
                        Line([Textbox('', type="text", placeholder='textbox', value=context.ip)]),
                        Line([CustomButton(innertext='Button')]),
                        Line([Button(selector='small', innertext='Button')]),
                        Line([Button(selector='mini', innertext='Button')])
                    ]
                )
            ], margin=Margin(left=Size.pixel(100)))
        ]
    )
