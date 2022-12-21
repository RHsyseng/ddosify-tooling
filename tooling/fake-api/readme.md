# Fake DDosify API

This is a simple fake API that can be used for demo purposes or to help with the development of new features.

Official API documentation can be found [here](https://docs.ddosify.com/cloud/api).

## Run it locally

> **INFO**: The api requires sudo as it needs to run on port 443, other options might be adding NET_BIND_SERVICE capability or tune sysctl conf. 

~~~sh
sudo python3.10 ./api.py
~~~

## Configure clients connecting to the fake API

Since the cli/operator code rely on valid TLS configurations, we need to do a few tunings in the server that will run the operator/cli.

1. Add the certificate to the trust-store:

    ~~~sh
    sudo cp cert.pem /etc/pki/ca-trust/source/anchors/
    sudo update-ca-trust
    ~~~

2. Fake the `api.ddosify.com` address

    ~~~
    echo 127.0.0.1 api.ddosify.com | sudo tee -a /etc/hosts
    ~~~

## Change results

You can edit [this file](./templates/latency-test-result-200.json) and put whatever results you want.