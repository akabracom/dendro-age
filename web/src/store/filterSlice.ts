import { createSlice, type PayloadAction } from "@reduxjs/toolkit";

export type PatternFilterState = {
  dateFrom: string;
  dateTo: string;
  title: string;
};

const initialState: PatternFilterState = {
  dateFrom: "1700-01-01",
  dateTo: "2000-12-31",
  title: "",
};

const filterSlice = createSlice({
  name: "filter",
  initialState,
  reducers: {
    setDateFrom(state, action: PayloadAction<string>) {
      state.dateFrom = action.payload;
    },
    setDateTo(state, action: PayloadAction<string>) {
      state.dateTo = action.payload;
    },
    setTitle(state, action: PayloadAction<string>) {
      state.title = action.payload;
    },
    resetFilter() {
      return initialState;
    },
  },
});

export const { setDateFrom, setDateTo, setTitle, resetFilter } = filterSlice.actions;
export const filterReducer = filterSlice.reducer;
