import { DendroScalePatterns } from "./generated/DendroScalePatterns";
import { HttpClient } from "./generated/http-client";
import { getApiToken } from "./session";

export const http = new HttpClient({
  baseURL: "/",
  securityWorker: async () => {
    const token = getApiToken();
    if (!token) {
      return {};
    }
    return { headers: { Authorization: `Bearer ${token}` } };
  },
});

http.instance.interceptors.request.use((config) => {
  const token = getApiToken();
  if (token && !config.headers.get("Authorization")) {
    config.headers.set("Authorization", `Bearer ${token}`);
  }
  if (typeof FormData !== "undefined" && config.data instanceof FormData) {
    config.headers.set("Content-Type", false as unknown as string);
  }
  return config;
});

export const dendroScalePatterns = new DendroScalePatterns(http);
