# ☢️ AtomBot ☢️

Golang bot to watch atom feeds and post new content to an API.
## Usage
Access the project's directory and build the executable.
```bash
go build -tags netgo -a -v
```

Build a docker container through the Dockerfile
```bash
docker build .
```

Make sure to provide the necessary environment variables when starting
```bash
  docker run -d -e ROOM=roomId -e AUTH=authKey -e CHATAPI=apiUrl \
  -e RSS=feedLink oborek/atombot:latest
```

## Roadmap

Currently, this bot is adapted to read atom feeds and post to the Webex API.

- WebApp for Config and Deployment

- Multi RSS/Atom feed support

- Support for more chat APIs