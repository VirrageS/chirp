FROM node:7.6

# Fix bug https://github.com/npm/npm/issues/9863
RUN cd $(npm root -g)/npm \
  && npm install fs-extra \
  && sed -i -e s/graceful-fs/fs-extra/ -e s/fs\.rename/fs.move/ ./lib/utils/rename.js

# install yarn
RUN \
  curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add - && \
  echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list && \
  apt-get update && apt-get install yarn

ENV NODE_ENV=development

# install global dependencies
RUN yarn global add typescript webpack webpack-dev-server tslint typings

# install all packes
RUN mkdir /frontend
COPY package.json /frontend/
WORKDIR /frontend/

RUN yarn install

# copy other files
COPY . /frontend
