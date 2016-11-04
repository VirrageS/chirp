# Chirp

Chirp is simplified Twitter written in Angular 2 and Go.


## Delete unused branches

    $ git fetch --all --prune
    $ git branch --merged master | grep -v 'master$' | xargs git branch -d


## Getting started (BACKEND)


## Getting started (FRONTEND)

You should get `Node > 6.x` and `npm > 3.x`.


Now run

    $ npm install --global typescript webpack webpack-dev-server tslint
    $ npm install
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
