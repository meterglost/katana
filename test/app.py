import logging
from flask import Flask, render_template, request
from urllib.parse import unquote

app = Flask(__name__)
logging.getLogger('werkzeug').disabled = True

@app.route("/", methods=['GET', 'POST'])
def main():
    if request.method == 'POST':
        app.logger.info('%s [%s]', request.url, unquote(request.get_data()))
        return "OK"
    else:
        return render_template('index.html')