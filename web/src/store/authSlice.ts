import { createSlice, type PayloadAction } from "@reduxjs/toolkit";

export type AuthState = {
  status: "guest" | "user";
  login: string;
  token: string;
  loading: boolean;
  error: string;
};

const initialState: AuthState = {
  status: "guest",
  login: "",
  token: "",
  loading: false,
  error: "",
};

const authSlice = createSlice({
  name: "auth",
  initialState,
  reducers: {
    authStarted(state) {
      state.loading = true;
      state.error = "";
    },
    authSucceeded(state, action: PayloadAction<{ login: string; token: string }>) {
      state.loading = false;
      state.status = "user";
      state.login = action.payload.login;
      state.token = action.payload.token;
      state.error = "";
    },
    authFailed(state, action: PayloadAction<string>) {
      state.loading = false;
      state.error = action.payload;
    },
    authCleared() {
      return initialState;
    },
  },
});

export const { authStarted, authSucceeded, authFailed, authCleared } = authSlice.actions;
export const authReducer = authSlice.reducer;
