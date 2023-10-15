# Chess CLI - command line chess game

> Multiplayer chess game played using command line interface with dedicated server.

## Configuration

To override default configuration, copy file `./config.yaml` to `./config.local.yaml` and edit all necessary
configuration properties like database connection and server port.

Use this new `config.local.yaml` file path as a first argument when executing migrations or server command.

## Migrations

New migrations should be added in `./migrations` directory with following naming convention:
`{year}{month}{day}{hour}{minute}_new_migration_name.{up|down}.sql`

Migrations are executed on application startup (if configured) or by running following command:
`go run ./migrations/main.go ./config.local.yaml`

## Swagger

Generate swagger documentation when making changes to comments and structs:

```shell
swag init -g ./cmd/chess-server/main.go -o ./docs
```

Swagger page is available on following URL: `http://{host}:{port}/swagger/index.html`

## Docker support

The game server service has full docker support provided by [Dockerfile](Dockerfile).

There is also [docker-compose.yaml](docker-compose.yaml) file with PostgreSQL database along the chess server.

The easiest way to pull & run docker image is to use already built public image
from [official DockerHub repository](https://hub.docker.com/repository/docker/lukamatosevic/chess-server):

```sh
docker pull lukamatosevic/chess-server:latest

docker run --env-file .env lukamatosevic/chess-server
```

Or, you can build the image yourself with docker command:

```sh
docker image build --rm -t chess-server .
```

Then you can start the chess server service with any of the following commands:

```sh
# provide .env file with environment variables to override default config
docker run --env-file .env chess-server

# also, you can set environment variables as parameters to override default config
docker run --env "DATABASE_HOST=127.0.0.1" \
           --env "DATABASE_USERNAME=chess-cli" \
           --env "DATABASE_PASSWORD=..." chess-server
```

## CMD

### Chess Server

This command starts the chess server API. Accepts one optional argument: the path for the config.yaml file.

Example:

```shell
go run ./cmd/chess-server ./config.local.yaml
```

### Chess CLI

This command starts command line interface client for playing the chess on desired server. You can play in interactive
mode or by using separate commands to make perform actions (login, create game, join game, play move, etc.)

The interactive mode is started by not providing any sub-commands defined in usage (you can only use global flags).

```shell
go run ./cmd/chess-cli --server http://localhost:64355
````

Usage:

```shell
NAME:
   Chess CLI - Play a game of chess using command line interface

USAGE:
   Chess CLI [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
   info, i    show server info
   events, e
   help, h    Shows a list of commands or help for one command
   auth:
     register, r
     login, l
     logout, o
     changePassword, c
     whoami, w
   games:
     game, g, games
   players:
     player, p, players

GLOBAL OPTIONS:
   --server value, -s value    chess server base URL
   --username value, -u value  players username or email
   --password value, -p value  players password
   --token value, -t value     players access token
   --stateless, -l             do not use and set default configs from home directory (default: false)
   --help, -h                  show help
   --version, -v               print the version
```

#### Players

```shell
NAME:
   Chess CLI player

USAGE:
   Chess CLI player command [command options] [arguments...]

COMMANDS:
   list     list all players
   info     show information about the player
   help, h  Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

#### Games

```shell
NAME:
   Chess CLI game

USAGE:
   Chess CLI game command [command options] [arguments...]

COMMANDS:
   list     list all games
   info     show information about the game
   create   create new game
   join     join existing game
   quit     quit currently active game
   play     play move in currently active game
   manual   Shows the instructions for all types of available moves
   help, h  Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```
