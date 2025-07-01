# Attila
Attila is a meta application for the Nomad workload orchestrator. It provides high level
abstractions and operator tooling which should not exist in Nomad, and in particular, aims to help operate non-federate, cell based cluster deployments.

Note: the project is under development and should be considered unstable and not production ready.

## Building
The Attila binary can be built using the included makefile. You will need [Go](https://go.dev/) installed on your machine along with `make` which is provided by
[Command Line Tools](https://developer.apple.com/xcode/resources/) on macOS or the
[build-essential](https://launchpad.net/ubuntu/+source/build-essential) package on Linux.
```console
$ make build
==> Building Attila...
==> Done
```

The compiled binary will be available within `./bin/`.

## Running
The [docs](./docs/) directory contains reference material for running an Attila server and using the binary. To quickly start, you can run `./bin/attila server run` and use the binary to discover the available commands and options.

The [demo](./docs/demo) directory contains an end-to-end example of running, configuring, and using Attila to run a Nomad job.

## Roadmap
The list below are some initial ideas of what will be added to Attila. If you have thoughts or other ideas, please open an issue, providing as much detail as possible about use-cases, expected UX, and
such.

* Region Capacity CLI: list and info commands which provide capacity details on Nomad regions.

* Advanced Meta Scheduling Capabilities: allow operators to define meta scheduling rules, such as "register job to region closest to these coordinates" which can be used to dictate job
registrations.

* Server High Availability: allow multiple instances of the Attila server to run and share state, so
that failures can be tolerated.

* Server Persistent State: currently state is persisted within memory on the Attila server which is
not persisted across process interruptions.
