const fs = require("fs");
const os = require("os");
const path = require("path");
const selfsigned = require("selfsigned");

const certDir = path.join(__dirname, "certs");
const keyPath = path.join(certDir, "key.pem");
const certPath = path.join(certDir, "cert.pem");

function laptopAddresses() {
  const addresses = ["127.0.0.1"];
  for (const entries of Object.values(os.networkInterfaces())) {
    for (const entry of entries || []) {
      if (entry.family === "IPv4" && !entry.internal) {
        addresses.push(entry.address);
      }
    }
  }
  return [...new Set(addresses)];
}

function createCertificate() {
  const ips = laptopAddresses();
  const altNames = [{ type: 2, value: "localhost" }, ...ips.map((ip) => ({ type: 7, ip }))];
  const pems = selfsigned.generate([{ name: "commonName", value: "dendro-age" }], {
    days: 365,
    keySize: 2048,
    algorithm: "sha256",
    extensions: [{ name: "subjectAltName", altNames }],
  });
  fs.mkdirSync(certDir, { recursive: true });
  fs.writeFileSync(keyPath, pems.private);
  fs.writeFileSync(certPath, pems.cert);
  console.log("Сертификат записан:");
  console.log(certPath);
  console.log(keyPath);
  console.log("Адреса в сертификате:", ips.join(", "));
  return { keyPath, certPath, ips };
}

if (require.main === module) {
  createCertificate();
}

module.exports = { createCertificate, keyPath, certPath };
