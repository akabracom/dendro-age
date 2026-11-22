const fs = require("fs");
const https = require("https");
const os = require("os");
const path = require("path");
const express = require("express");
const { createProxyMiddleware } = require("http-proxy-middleware");
const { createCertificate, keyPath, certPath } = require("./generate-cert");

const port = 3000;
const dist = path.resolve(__dirname, "../web/dist");
const apiTarget = "http://127.0.0.1:8088";
const minioTarget = "http://127.0.0.1:9000";

if (!fs.existsSync(path.join(dist, "index.html"))) {
  console.error("Нет сборки фронтенда. В папке web выполни: npm run build");
  process.exit(1);
}

if (!fs.existsSync(certPath) || !fs.existsSync(keyPath)) {
  createCertificate();
}

const app = express();

app.use(
  "/api",
  createProxyMiddleware({
    target: apiTarget,
    changeOrigin: true,
  }),
);
app.use(
  "/static",
  createProxyMiddleware({
    target: apiTarget,
    changeOrigin: true,
  }),
);
app.use(
  "/minio",
  createProxyMiddleware({
    target: minioTarget,
    changeOrigin: true,
    pathRewrite: { "^/minio": "" },
  }),
);

app.get("/", (_req, res) => {
  res.redirect("/dendro-age/");
});

app.use("/dendro-age", express.static(dist));
app.get("/dendro-age/*", (_req, res) => {
  res.sendFile(path.join(dist, "index.html"));
});

const server = https.createServer(
  {
    key: fs.readFileSync(keyPath),
    cert: fs.readFileSync(certPath),
  },
  app,
);

server.listen(port, "0.0.0.0", () => {
  console.log("HTTPS-прокси Dendro Age слушает порт", port);
  console.log("API", apiTarget);
  console.log("MinIO", minioTarget);
  const addresses = ["127.0.0.1"];
  for (const entries of Object.values(os.networkInterfaces())) {
    for (const entry of entries || []) {
      if (entry.family === "IPv4" && !entry.internal) {
        addresses.push(entry.address);
      }
    }
  }
  for (const address of addresses) {
    console.log(`https://${address}:${port}/dendro-age/`);
  }
});
