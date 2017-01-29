# Chirp

[![Build Status](https://travis-ci.org/VirrageS/chirp.svg?branch=master)](https://travis-ci.org/VirrageS/chirp)
[![Go Report Card](https://goreportcard.com/badge/github.com/VirrageS/chirp)](https://goreportcard.com/report/github.com/VirrageS/chirp)
[![GoDoc](https://godoc.org/github.com/VirrageS/chirp?status.svg)](https://godoc.org/github.com/VirrageS/chirp)
[![CircleCI](https://circleci.com/gh/VirrageS/chirp/tree/master.svg?style=svg)](https://circleci.com/gh/VirrageS/chirp/tree/master)


Chirp is simplified Twitter written in Angular 2 and Go. You can start
fully working website with just one line.


## Setup

Before we begin we have to install `docker-compose` command [Install](https://docs.docker.com/compose/install/).
Important is that `docker-compose` should have `>1.6` version!

    $ make production

To be able to use system you need to add to `/etc/hosts`:

```
127.0.0.1   backend.show frontend.show
```

It is because we are not using any external domains yet. Then you can just hit
`frontend.show/` and now you are able to access fully working project ^^.
