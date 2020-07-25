const { join } = require("path");
const express = require("express");
const app = express();
const WebSocket = require("ws");
const http = require("http");

const port = 3000;

app.get("/", (req, res) => res.sendFile(join(__dirname, "bike.html")));
app.get("/mock", (req, res) => res.sendFile(join(__dirname, "bike-mock.html")));
app.get("/pubsub.js", (req, res) => res.sendFile(join(__dirname, "pubsub.js")));

const server = http.createServer(app);

const wss = new WebSocket.Server({ server });
wss.on("connection", (ws) => {
  ws.on("message", (message) => {
    console.log("received: %s", message);

    wss.clients.forEach((client) => {
      if (client !== ws) {
        client.send(message);
      }
    });
  });

  ws.send("something");
});

server.listen(port, () =>
  console.log(`Bike app listening at http://localhost:${port}`)
);
