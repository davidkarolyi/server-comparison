FROM node:latest

WORKDIR /usr/src/app
COPY ./test_data.json .

WORKDIR /usr/src/app/node_fastify
COPY ./node_fastify/ .

RUN npm install

EXPOSE 3000

CMD [ "node", "index.js" ]