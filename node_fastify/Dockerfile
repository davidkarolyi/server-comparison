FROM node:latest

WORKDIR /usr/src/app
COPY ./test_data.json .

WORKDIR /usr/src/app/node_fastify
COPY ./node_fastify/package*.json .

RUN npm install

COPY ./node_fastify/index.js .

EXPOSE 3000

CMD [ "node", "index.js" ]