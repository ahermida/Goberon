# Goberon
Villanova Course Indexing CLI. Rewrite of Oberon (https://github.com/Space-Cadets/Oberon) in Golang to work with RethinkDB. Built over a long weekend.

<!-- toc -->

- [Overview](#overview)
- [Installation](#installation)
  * [Supported platforms](#supported-platforms)
  * [RethinkDB Setup](#request-setup)
  * [Request Setup](#request-setup)
- [Getting Started](#getting-started)

<!-- tocstop -->


## Overview

The Villanova course registry is a mess and they've taken away Schedulr. Things like
finding interesting electives to take shouldn't be as difficult as they are.
This is the first step to organizing course data and making registration less of a
hassle.

## Installation

Make sure you have a working Go environment.  Go version 1.8+ is supported.  [See
the install instructions for Go](http://golang.org/doc/install.html).

To install cli, simply run:
```
$ go get github.com/ahermida/Goberon
```

Make sure your `PATH` includes the `$GOPATH/bin` directory so your commands can
be easily used:
```
export PATH=$PATH:$GOPATH/bin
```

### Supported platforms

This CLI is currently only tested against MacOS and the Linux Subsystem for Windows.

### RethinkDB setup

This project requires a RethinkDB instance to be running. If you don't have it
installed yet, you can see the [RethinkDB installation guide](https://rethinkdb.com/docs/install/)
for your system.

Once you have it installed, you must set the "RdbAddress" value within
the "Network" variable in `Goberon/config/config.go` to the port that rethinkdb is
open on.

### Request setup

Go into `Goberon/config/config.go`, set
