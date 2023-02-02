# subscribe

Service that processes subscribe operations requests

## Hooks

This repository is configured with client-side Git hooks which you need to install by running the following command:

```bash
./hooks/INSTALL
```

## Usage

To properly run this service, you will need to a setup a `.env` file. Start by creating a copy of the `.env.tpl` file and fill the variables with values appropriate for the execution context.

Then, all you need to do is to run the service with the following command:

```bash
go run cmd/subscribe/subscribe.go
```

## Docker

To run the service with Docker, you will first need to setup the `.git-local-credentials` file. This credentials file shall contain the git credentials config to access `damn` and `palavrapasse` private modules.

To build the service image:

```bash
docker_tag=subscribe:latest
docker build \
    -f ./deployments/Dockerfile \
    --secret id=git-credentials,src=.local-git-credentials \
    . -t $docker_tag
```

To run the service container:

```bash
export $(grep -v '^#' .env | xargs)
docker run \
    -p $server_port:$server_port \
    --mount "type=bind,src=$leaksdb_fp,dst=$leaksdb_fp" \
    --mount "type=bind,src=$subscriptionsdb_fp,dst=$subscriptionsdb_fp" \
    --env-file .env \
    -t $docker_tag
```