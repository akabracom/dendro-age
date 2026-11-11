import { http } from "./dendroScalePatterns";
import { ACCESS_TOKEN_KEY, getApiToken, setApiToken } from "./session";

export type LoginResponse = {
  access_token: string;
  token_type: string;
  expires_in: number;
  login: string;
};

export async function loginResearcher(login: string, password: string) {
  const response = await http.instance.post<LoginResponse>("/api/dendro_researchers/login", {
    login,
    password,
  });
  localStorage.setItem(ACCESS_TOKEN_KEY, response.data.access_token);
  setApiToken(response.data.access_token);
  return response.data;
}

export async function logoutResearcher() {
  const token = getApiToken();
  try {
    await http.instance.post("/api/dendro_researchers/logout", null, {
      headers: token ? { Authorization: `Bearer ${token}` } : {},
    });
  } finally {
    localStorage.removeItem(ACCESS_TOKEN_KEY);
    setApiToken(null);
  }
}
