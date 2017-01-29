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

    $ yarn global add typescript webpack webpack-dev-server tslint
    $ yarn install

now open browser to [localhost:3000](http://localhost:3000/) and done! :)



## Contribution (BACKEND)

To test or format code in backend you need to install some additional tools.
To run tests you need to install:

    (ubuntu)$ sudo apt-get install postgresql-client
    (mac)$ brew install postgresql

To be able to use full code formatting you need to install:

    $ go get -u github.com/golang/lint/golint

Then to test or format code run (**tip**: tests require Docker running!):

    $ make test -C ./backend
    $ make format -C ./backend



## Services

You can also start services with using docker:

    $ make core           # basic services
    $ make frontend       # frontend with services
    $ make backend        # backend with services

For more details check `Makefile` and `docker/*.yml` files
