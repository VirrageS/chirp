# Chirp

[![Build Status](https://travis-ci.org/VirrageS/chirp.svg?branch=master)](https://travis-ci.org/VirrageS/chirp)
[![Go Report Card](https://goreportcard.com/badge/github.com/VirrageS/chirp)](https://goreportcard.com/report/github.com/VirrageS/chirp)
[![GoDoc](https://godoc.org/github.com/VirrageS/chirp?status.svg)](https://godoc.org/github.com/VirrageS/chirp)
[![CircleCI](https://circleci.com/gh/VirrageS/chirp/tree/master.svg?style=svg)](https://circleci.com/gh/VirrageS/chirp/tree/master)


Chirp is simplified Twitter written in Angular 2 and Go.



## Delete unused branches

    $ git fetch --all --prune
    $ git branch --merged master | grep -v 'master$' | xargs git branch -d



## Getting started (BACKEND)

Install Go language: https://golang.org/doc/install (don't forget to set your GOPATH).

Install postgresql client tools - will be required to run tests:

    $ sudo apt-get install postgresql-client

Install Go-lint - will be required for code formatting:

    $ go get -u github.com/golang/lint/golint

In order for backed to work you also need to have docker and docker-compose, see [Docker](https://github.com/VirrageS/chirp#docker)

Now run

    $ go get github.com/VirrageS/chirp
    $ cd $GOPATH/src/github.com/VirrageS/chirp/backend
    $ make install
    $ docker-compose -f docker/backend.yml build && docker-compose -f docker/backend.yml up
    $ $GOPATH/bin/backend

Now you've got your chirp backend running on [localhost:8080](http://localhost:8080/)!



### Running backend easier

You can add `$GOPATH/bin` to your `$PATH` and run `backend` easier.

    $ export PATH=$PATH:$GOPATH/bin
    $ backend



## Getting started (FRONTEND)

You should get `Node > 6.x` and `npm > 3.x`.


Now run

    $ npm install --global typescript webpack webpack-dev-server tslint
    $ npm install
    $ npm rebuild node-sass
    $ npm start

now open browser to [localhost:3000](http://localhost:3000/) and done! :)


## Docker

Before we begin we have to install `docker-compose` command [Install](https://docs.docker.com/compose/install/)

Then, depending on the services we want to start we have to type:


### Backend with services

    $ docker-compose -f docker/backend.yml build && docker-compose -f docker/backend.yml up

### Frontend with services

    $ docker-compose -f docker/frontend.yml build && docker-compose -f docker/frontend.yml up

### Basic services

    $ docker-compose -f docker/core.yml build && docker-compose -f docker/core.yml up

### Production

    $ docker-compose -f docker/production.yml build && docker-compose -f docker/production.yml up

### Backend testing

    /backend$ make test
