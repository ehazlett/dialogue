# Dialogue
This is a communication system, a la forum style messaging.

Dialogue uses RethinkDB for backend datastore.

# API
To build the api, `cd` into the `api` directory and run `godep go build`

You should then have an `api` executable.  Run `./api` to start the api server.

You can curl `http://localhost:3000/setup` to create the admin user.

