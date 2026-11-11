import { configureStore } from "@reduxjs/toolkit";
import { authReducer } from "./authSlice";
import { filterReducer } from "./filterSlice";
import { patternsReducer } from "./patternSlice";

export const store = configureStore({
  reducer: {
    auth: authReducer,
    filter: filterReducer,
    patterns: patternsReducer,
  },
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
