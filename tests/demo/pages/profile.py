from bevyframe import *
from TheProtocols import *

from bevyframe.Features.ContextManager import set_to_context_manager


def activity(context: Context) -> Activity:
    u = User(context.query['email'])
    return Activity(
        name=f"{u.name} {u.surname}",
        prop2="demo"
    )


def get(context: Context) -> Page:
    u = User(context.query['email'])
    # context.last_visited_profile = context.query['email']
    set_to_context_manager(context.tp.package_name, context.email, "last_visited_profile", context.query['email'])
    return Page(
        title='',
        description='',
        color=context.user.id.settings.theme_color,
        childs=[
            Title(f"{u.name} {u.surname}")
        ]
    )
