const { join } = require("path");
const express = require("express");
const app = express();

const port = 3000;

const content = "Hello World!";

app.get("/", (req, res) => res.sendFile(join(__dirname, "index.html")));
app.get("/data", (req, res) => res.send(content));

app.listen(port, () =>
  console.log(`Example app listening at http://localhost:${port}`)
);
