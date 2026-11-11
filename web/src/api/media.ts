export function proxyMedia(url: string, fallback: string) {
  const value = url.trim();
  if (!value) {
    return fallback;
  }
  return value
    .replace(/^https?:\/\/localhost:9000/i, "/minio")
    .replace(/^https?:\/\/127\.0\.0\.1:9000/i, "/minio");
}
