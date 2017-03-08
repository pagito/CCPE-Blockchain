#!/bin/sh

FROM s390x/node:6.9.1

# Create app directory
run ["mkdir", "-p", "/usr/src/b-app"]
WORKDIR /usr/src/b-app

# Install app dependencies
COPY package.json /usr/src/b-app/
run npm install

# Bundle app source
COPY . /usr/src/b-app

EXPOSE 8088
CMD [ "npm", "start" ]