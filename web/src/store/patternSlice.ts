import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import axios from "axios";
import { dendroScalePatterns } from "../api/dendroScalePatterns";
import { proxyMedia } from "../api/media";
import { getApiToken } from "../api/session";
import type {
  CreateDendroScalePatternPayload,
  DendroScalePattern as ApiPattern,
} from "../api/generated/data-contracts";
import { mockPatterns, publishedPatterns } from "../mock/patterns";
import type { DendroScalePattern } from "../types";
import type { PatternFilterState } from "./filterSlice";

type PatternLookup = {
  filter: PatternFilterState;
  patterns: {
    source: "api" | "mock";
    list: DendroScalePattern[];
    feed: DendroScalePattern | null;
  };
};

type PatternState = {
  list: DendroScalePattern[];
  feed: DendroScalePattern | null;
  draft: DendroScalePattern | null;
  loading: boolean;
  error: string;
  source: "api" | "mock";
};

const initialState: PatternState = {
  list: [],
  feed: null,
  draft: null,
  loading: false,
  error: "",
  source: "api",
};

function toPattern(raw: ApiPattern): DendroScalePattern {
  return {
    id: raw.id ?? 0,
    title: raw.title ?? "",
    shortDescription: raw.short_description ?? "",
    fullDescription: raw.full_description ?? "",
    status: raw.status === "draft" || raw.status === "deleted" ? raw.status : "published",
    imageUrl: proxyMedia(raw.image_url ?? "", ""),
    videoUrl: proxyMedia(raw.video_url ?? "", ""),
    eventYear: raw.event_year ?? null,
    ringCount: raw.ring_count ?? null,
    likeCount: raw.like_count ?? 0,
    liked: raw.liked ?? 0,
    mine: raw.mine ?? 0,
    creatorId: raw.creator_id ?? 0,
  };
}

function isOffline(error: unknown) {
  if (!axios.isAxiosError(error)) {
    return false;
  }
  if (!error.response) {
    return true;
  }
  return error.response.status >= 500;
}

function errorText(error: unknown) {
  if (axios.isAxiosError(error)) {
    const message = error.response?.data?.message;
    if (typeof message === "string" && message) {
      return message;
    }
  }
  return "Сервис недоступен";
}

function mockList(dateFrom: string, dateTo: string, title: string) {
  const fromYear = yearOf(dateFrom);
  const toYear = yearOf(dateTo);
  const query = title.trim().toLowerCase();
  return publishedPatterns(mockPatterns).filter((item) => {
    if (item.eventYear === null) {
      return false;
    }
    if (fromYear !== null && item.eventYear < fromYear) {
      return false;
    }
    if (toYear !== null && item.eventYear > toYear) {
      return false;
    }
    if (query && !item.title.toLowerCase().includes(query)) {
      return false;
    }
    return true;
  });
}

function yearOf(raw: string) {
  if (!raw) {
    return null;
  }
  const parsed = new Date(raw + "T00:00:00");
  if (Number.isNaN(parsed.getTime())) {
    return null;
  }
  return parsed.getFullYear();
}

export const fetchDendroScalePatterns = createAsyncThunk(
  "patterns/fetchDendroScalePatterns",
  async (_, { getState, rejectWithValue }) => {
    const filter = (getState() as PatternLookup).filter;
    try {
      const response = await dendroScalePatterns.getDendroScalePatterns({
        date_from: filter.dateFrom,
        date_to: filter.dateTo,
        title: filter.title,
      });
      return { list: response.data.map(toPattern), mock: false as const };
    } catch (error) {
      if (!getApiToken() && isOffline(error)) {
        return { list: mockList(filter.dateFrom, filter.dateTo, filter.title), mock: true as const };
      }
      return rejectWithValue(errorText(error));
    }
  },
);

export const fetchDendroScalePatternFeed = createAsyncThunk(
  "patterns/fetchDendroScalePatternFeed",
  async (arg: { id?: number; next?: boolean } | undefined, { getState, rejectWithValue }) => {
    const state = getState() as PatternLookup;
    if (state.patterns.source === "mock" && arg?.next) {
      const items = state.patterns.list;
      const index = items.findIndex((item) => item.id === state.patterns.feed?.id);
      const pattern = items[(index + 1) % items.length];
      return { pattern, mock: true as const };
    }
    try {
      if (arg?.id && arg.next) {
        const response = await dendroScalePatterns.http.instance.get<ApiPattern>(
          `/api/dendro_scale_patterns/feed/${arg.id}`,
          { params: { next: true } },
        );
        return { pattern: toPattern(response.data), mock: false as const };
      }
      const response = await dendroScalePatterns.getDendroScalePatternFeed();
      return { pattern: toPattern(response.data), mock: false as const };
    } catch (error) {
      if (!getApiToken() && isOffline(error)) {
        const list = publishedPatterns(mockPatterns);
        return { pattern: list[0] ?? null, list, mock: true as const };
      }
      return rejectWithValue(errorText(error));
    }
  },
);

export const fetchDendroScalePatternDraft = createAsyncThunk(
  "patterns/fetchDendroScalePatternDraft",
  async (_, { rejectWithValue }) => {
    try {
      const response = await dendroScalePatterns.getDendroScalePatternDraft();
      return toPattern(response.data);
    } catch (error) {
      if (axios.isAxiosError(error) && error.response?.status === 404) {
        return null;
      }
      return rejectWithValue(errorText(error));
    }
  },
);

export const createDendroScalePattern = createAsyncThunk(
  "patterns/createDendroScalePattern",
  async (payload: { title: string; image?: File; video?: File }, { rejectWithValue }) => {
    try {
      const body: CreateDendroScalePatternPayload = { title: payload.title };
      if (payload.image) {
        body.image = payload.image;
      }
      if (payload.video) {
        body.video = payload.video;
      }
      const response = await dendroScalePatterns.createDendroScalePattern(body);
      return toPattern(response.data);
    } catch (error) {
      return rejectWithValue(errorText(error));
    }
  },
);

export const publishDendroScalePattern = createAsyncThunk(
  "patterns/publishDendroScalePattern",
  async (
    payload: { shortDescription: string; fullDescription: string; eventYear: number; ringCount: number },
    { rejectWithValue },
  ) => {
    try {
      const response = await dendroScalePatterns.publishDendroScalePattern({
        short_description: payload.shortDescription,
        full_description: payload.fullDescription,
        event_year: payload.eventYear,
        ring_count: payload.ringCount,
      });
      return toPattern(response.data);
    } catch (error) {
      return rejectWithValue(errorText(error));
    }
  },
);

export const deleteDendroScalePattern = createAsyncThunk(
  "patterns/deleteDendroScalePattern",
  async (id: number, { dispatch, rejectWithValue }) => {
    try {
      await dendroScalePatterns.deleteDendroScalePattern(id);
      await dispatch(fetchDendroScalePatterns());
      return id;
    } catch (error) {
      return rejectWithValue(errorText(error));
    }
  },
);

export const likeDendroScalePattern = createAsyncThunk(
  "patterns/likeDendroScalePattern",
  async (payload: { id: number; like: 0 | 1 }, { rejectWithValue }) => {
    try {
      const response = await dendroScalePatterns.likeDendroScalePattern(payload.id, { like: payload.like });
      return toPattern(response.data);
    } catch (error) {
      return rejectWithValue(errorText(error));
    }
  },
);

const patternSlice = createSlice({
  name: "patterns",
  initialState,
  reducers: {
    clearDraft(state) {
      state.draft = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchDendroScalePatterns.pending, (state) => {
        state.loading = true;
        state.error = "";
      })
      .addCase(fetchDendroScalePatterns.fulfilled, (state, action) => {
        state.loading = false;
        state.list = action.payload.list;
        state.source = action.payload.mock ? "mock" : "api";
      })
      .addCase(fetchDendroScalePatterns.rejected, (state, action) => {
        state.loading = false;
        state.error = String(action.payload ?? "Ошибка списка");
      })
      .addCase(fetchDendroScalePatternFeed.pending, (state) => {
        state.loading = true;
        state.error = "";
      })
      .addCase(fetchDendroScalePatternFeed.fulfilled, (state, action) => {
        state.loading = false;
        state.feed = action.payload.pattern;
        if (action.payload.mock) {
          state.source = "mock";
          if ("list" in action.payload && action.payload.list) {
            state.list = action.payload.list;
          }
        } else {
          state.source = "api";
        }
      })
      .addCase(fetchDendroScalePatternFeed.rejected, (state, action) => {
        state.loading = false;
        state.error = String(action.payload ?? "Ошибка ленты");
      })
      .addCase(fetchDendroScalePatternDraft.pending, (state) => {
        state.loading = true;
        state.error = "";
      })
      .addCase(fetchDendroScalePatternDraft.fulfilled, (state, action) => {
        state.loading = false;
        state.draft = action.payload;
      })
      .addCase(fetchDendroScalePatternDraft.rejected, (state, action) => {
        state.loading = false;
        state.error = String(action.payload ?? "Ошибка черновика");
      })
      .addCase(createDendroScalePattern.pending, (state) => {
        state.loading = true;
        state.error = "";
      })
      .addCase(createDendroScalePattern.fulfilled, (state, action) => {
        state.loading = false;
        state.draft = action.payload;
      })
      .addCase(createDendroScalePattern.rejected, (state, action) => {
        state.loading = false;
        state.error = String(action.payload ?? "Ошибка создания");
      })
      .addCase(publishDendroScalePattern.pending, (state) => {
        state.loading = true;
        state.error = "";
      })
      .addCase(publishDendroScalePattern.fulfilled, (state) => {
        state.loading = false;
        state.draft = null;
      })
      .addCase(publishDendroScalePattern.rejected, (state, action) => {
        state.loading = false;
        state.error = String(action.payload ?? "Ошибка публикации");
      })
      .addCase(deleteDendroScalePattern.pending, (state) => {
        state.loading = true;
        state.error = "";
      })
      .addCase(deleteDendroScalePattern.fulfilled, (state) => {
        state.loading = false;
      })
      .addCase(deleteDendroScalePattern.rejected, (state, action) => {
        state.loading = false;
        state.error = String(action.payload ?? "Ошибка удаления");
      })
      .addCase(likeDendroScalePattern.pending, (state) => {
        state.loading = true;
        state.error = "";
      })
      .addCase(likeDendroScalePattern.fulfilled, (state, action) => {
        state.loading = false;
        state.feed = action.payload;
        state.list = state.list.map((item) => (item.id === action.payload.id ? action.payload : item));
      })
      .addCase(likeDendroScalePattern.rejected, (state, action) => {
        state.loading = false;
        state.error = String(action.payload ?? "Ошибка лайка");
      });
  },
});

export const { clearDraft } = patternSlice.actions;
export const patternsReducer = patternSlice.reducer;
