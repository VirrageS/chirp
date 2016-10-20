FROM node:6.9

# Fix bug https://github.com/npm/npm/issues/9863
RUN cd $(npm root -g)/npm \
  && npm install fs-extra \
  && sed -i -e s/graceful-fs/fs-extra/ -e s/fs\.rename/fs.move/ ./lib/utils/rename.js


WORKDIR /frontend
ADD . /frontend

RUN npm install --global typescript webpack webpack-dev-server tslint typings yarn
RUN yarn install
RUN typings install
RUN npm run build
