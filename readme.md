# Dialogue
This is a communication system, a la forum style messaging.

Dialogue uses RethinkDB for backend datastore.

# API
To build the api, `cd` into the `api` directory and run `make`

You should then have an `api` executable.  Run `./api` to start the api server.

You can curl `http://localhost:3000/setup` to create the admin user.

# CLI
To build the cli, `cd` into the `cli` directory and run `make`.

## Login
You will need to have the API running and a user account.

`./dialogue login`

## Usage
Once logged in, an authorization token will be stored for easier use.

### Show Topics
`./dialogue topics list`

### Create Topic
`./dialogue topics create --title foo`

### Delete Topic
`./dialogue topics delete --id e67ea2bf-8df2-41ff-b845-b325641c748f`

### Show Posts
`./dialogue posts list --topicId 6ba7c765-fd5e-45e2-bc03-2db969921391`

### Show Posts with IDs
`./dialogue posts list --topicId 6ba7c765-fd5e-45e2-bc03-2db969921391 --ids`

### Create Post
`./dialogue posts create --topicId 6ba7c765-fd5e-45e2-bc03-2db969921391 --content "Foo Content"`

### Delete Post
`./dialogue posts delete --id 1824fcf2-6eac-4edd-9c17-3e92cc6e3c8b`
