FROM node:latest

WORKDIR /usr/src/app
COPY ./test_data.json .

WORKDIR /usr/src/app/node_express
COPY ./node_express/package*.json .

RUN npm install

COPY ./node_express/index.js .

EXPOSE 3000

CMD [ "node", "index.js" ]