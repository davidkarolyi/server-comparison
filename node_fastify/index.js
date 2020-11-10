const fastify = require("fastify")();
const jwt = require("jsonwebtoken");
const fs = require("fs");
const path = require("path");
const { secret } = require("../test_data.json");

fastify.get("/", async (request, reply) => {
  const token = await readToken();
  const payload = verifyToken(token);
  return payload;
});

fastify.listen(3000, "0.0.0.0", () => {
  console.log(`node_fastify is listening on localhost:3000`);
});

async function readToken() {
  const tokenString = await fs.promises.readFile(
    path.resolve(__dirname, "../test_data.json"),
    { encoding: "utf8" }
  );
  return JSON.parse(tokenString).token;
}

function verifyToken(token) {
  return jwt.verify(token, secret);
}
