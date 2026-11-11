let accessToken: string | null = null;

export const ACCESS_TOKEN_KEY = "dendro_access_token";

export function setApiToken(token: string | null) {
  accessToken = token;
}

export function getApiToken() {
  return accessToken;
}
