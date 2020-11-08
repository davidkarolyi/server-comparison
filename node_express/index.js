const app = require("express")();
const jwt = require("jsonwebtoken");
const fs = require("fs");
const path = require("path");
const { secret } = require("../test_data.json");

app.get("/", async (req, res) => {
  const token = await readToken();
  const payload = verifyToken(token);
  res.json(payload);
});

app.listen(3000, () => {
  console.log("node_express is listening on localhost:3000");
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
