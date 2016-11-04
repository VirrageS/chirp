# Chirp

Chirp is simplified Twitter written in Angular 2 and Go.



## Delete unused branches

    $ git fetch --all --prune
    $ git branch --merged master | grep -v 'master$' | xargs git branch -d



## Getting started (BACKEND)

Install Go language: https://golang.org/doc/install (don't forget to set your GOPATH).

Now run

    $ go get github.com/VirrageS/chirp
    $ cd $GOPATH/src/github.com/VirrageS/chirp/backend
    $ go get .
    $ go install
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
    $ npm start

now open browser to [localhost:3000](http://localhost:3000/) and done! :)
