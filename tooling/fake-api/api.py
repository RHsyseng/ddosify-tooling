#!/usr/bin/env python3
# coding=utf-8

import os
import flask
import requests
from requests.packages.urllib3.exceptions import InsecureRequestWarning
requests.packages.urllib3.disable_warnings(InsecureRequestWarning)

app = flask.Flask(__name__)
app.config['TEMPLATES_AUTO_RELOAD'] = True

@app.route('/v1/balance')
def balance_resource():
    return flask.render_template('balance-200.json')

@app.route('/v1/latency/test', methods=['POST'])
def latency_test_resource():
    return flask.render_template('latency-test-result-200.json')

def run(port):
    app.run(host='::', port=port, debug=False, ssl_context=("./cert.pem","./cert.key"))

if __name__ == '__main__':
    port = os.getenv('API_PORT', 443)
    run(port)
