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

### RethinkDB Setup

This project requires a RethinkDB instance to be running. If you don't have it
installed yet, you can see the [RethinkDB installation guide](https://rethinkdb.com/docs/install/)
for your system. If you have Homebrew on a mac, install with `brew update && brew install rethinkdb`.

Once you have it installed, you must set the "RdbAddress" value within
the "Network" variable in `Goberon/config/config.go` to the port that rethinkdb is
open on.

### Request Setup

Go into `Goberon/config/secret.go`, edit out the comment and place the data for the
POST request to the registrar (set the cookie, body, headers, and URL). If you
don't know how to get this data, fear not.

You can alternatively navigate to the registrar HTML page with all of the courses, save
the HTML page with all of the courses into the Goberon directory with the name `courses.html`.

## Getting Started

Once you've run `go get github.com/ahermida/Goberon` and you're in the Goberon
directory, you can run `go install`.

This will setup the command: `goberon`

### Fetching Data

Running `$ goberon` will load a help menu.

### Fetching Data

Running `$ goberon fetch` will download the course data. It usually takes ~4mins.

### Indexing Data

Using `$ goberon index` will build course data, and write it into appropriate
RethinkDB tables.

###  Data

Using `$ goberon drop` will drop all course data tables

##Acknowledgments

Can't build this without a shoutout to the original Oberon authors:
* David Siah [@dsiah](https://github.com/dsiah)
* Kent Wu [@kentkwu](https://github.com/kentkwu)
