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