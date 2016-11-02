FROM node:6.9

RUN adduser --disabled-password --gecos '' docker

# Fix bug https://github.com/npm/npm/issues/9863
RUN cd $(npm root -g)/npm \
  && npm install fs-extra \
  && sed -i -e s/graceful-fs/fs-extra/ -e s/fs\.rename/fs.move/ ./lib/utils/rename.js

# yarn
RUN \
  apt-key adv --keyserver pgp.mit.edu --recv D101F7899D41F3C3 && \
  echo "deb http://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list && \
  apt-get update && apt-get install yarn

RUN yarn global add typescript webpack webpack-dev-server tslint typings

WORKDIR /frontend
COPY . /frontend

RUN yarn install
RUN typings install
