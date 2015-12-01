# wedirect

Open up control of the A record of a domain to the masses.

## Setup
* Install using `go get github.com/ethernetdan/wedirect`
* Create a new [Firebase](https://firebase.com/). Get the URL and authentication token for it from Secrets.
* Setup the domain you wish to use with [CloudFlare](https://cloudflare.com). Get an authentication token.
* Modify `config.json.example` with the proper configuration and save it as `config.json`

## Docker
A minimal Docker container can be built by running `./build-docker.sh <image name>`

## UI

Make sure you have the following installed

* node >= v4.0.0
* npm
* gulp

Then run the following commands in the UI directory

* npm install
* gulp
