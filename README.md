# Chirp

[![Build Status](https://travis-ci.org/VirrageS/chirp.svg?branch=master)](https://travis-ci.org/VirrageS/chirp)
[![Go Report Card](https://goreportcard.com/badge/github.com/VirrageS/chirp)](https://goreportcard.com/report/github.com/VirrageS/chirp)
[![GoDoc](https://godoc.org/github.com/VirrageS/chirp?status.svg)](https://godoc.org/github.com/VirrageS/chirp)
[![CircleCI](https://circleci.com/gh/VirrageS/chirp/tree/master.svg?style=svg)](https://circleci.com/gh/VirrageS/chirp/tree/master)


Chirp is simplified Twitter written in Angular 2 and Go.


## Getting started (BACKEND)

Install Go language: https://golang.org/doc/install (don't forget to set your GOPATH).
In order for backed to work you also need to have docker and docker-compose, see [Docker](https://github.com/VirrageS/chirp#docker)

Now run

    $ go get github.com/VirrageS/chirp
    $ cd $GOPATH/src/github.com/VirrageS/chirp/backend
    $ make install
    $ docker-compose -f docker/core.yml up --build
    $ $GOPATH/bin/backend

Now you've got your chirp backend running on [localhost:8080](http://localhost:8080/)!


### Running backend easier

You can add `$GOPATH/bin` to your `$PATH` and run `backend` easier.

    $ export PATH=$PATH:$GOPATH/bin
    $ backend



## Getting started (FRONTEND)

You should get `Node > 6.x`, `npm > 3.x` and `yarn`.

Now run

    $ yarn global add typescript webpack webpack-dev-server tslint typings
    $ yarn install

now open browser to [localhost:3000](http://localhost:3000/) and done! :)



## Docker

Before we begin we have to install `docker-compose` command [Install](https://docs.docker.com/compose/install/)

Then, depending on the services we want to start we have to type:


### Backend with services

    $ docker-compose -f docker/backend.yml up --build


### Frontend with services

    $ docker-compose -f docker/frontend.yml up --build


### Basic services

    $ docker-compose -f docker/core.yml up --build


### Production

    $ docker-compose -f docker/production.yml up --build

If you want use production Docker you have to add this line to `/etc/hosts`:

```
127.0.0.1   backend.show frontend.show
```

It is because we are not using any external domains yet. Then you can just hit
`frontend.show/` and now you are able to access fully working project ^^.



## Contribution (BACKEND)

To test or format code in backend you need to install some additional tools.
To run tests you need to install:

    (ubuntu)$ sudo apt-get install postgresql-client
    (mac)$ brew install postgresql

To be able to use full code formatting you need to install:

    $ go get -u github.com/golang/lint/golint

Then to test or format code run (**tip**: tests require Docker running!):

    /backend$ make test
    /backend$ make format
