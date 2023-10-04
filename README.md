# realtime-quiz

Implements a server which hosts and quiz and a client which can join the quiz.

## Implementation

Two golang binaries: server and client.  They communicate via RabbitMQ queues.

## Test

To run the flow E2E with one player, you can run:

> NB: This will start a docker detached container of RMQ and background the server process. If you don't want this, you may prefer running it manually via [Running](#running)

```bash
$ make test
```

To run the (small amount of) tests:

```bash
$ go test ./...
```

## Running

To run the server binary with N number of people for a quiz:

```bash
$ go run cmd/server/main.go -n <N>
```

To run the client binary with client name NAME for a quiz:

```bash
$ go run cmd/client/main.go --name <NAME>
```

## TODO

- TESTS
    - `quizsteps` need tests
    - `quiz`  need error tests
    - `repositories` needs tests, possibly with the spun up rabbitmq
- Enforce unique client names
- Respond to client on registration to confirm their entry
- Create manual user quizstep
- Think if quizstep is a good name for that business logic layer