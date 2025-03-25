from flask import Flask

application = Flask(__name__)


@application.route("/<path:path>")
def demo_flask(path) -> str:
    return "Hello Flask!"
