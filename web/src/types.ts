export type PatternStatus = "draft" | "published" | "deleted";

export type DendroScalePattern = {
  id: number;
  title: string;
  shortDescription: string;
  fullDescription: string;
  status: PatternStatus;
  imageUrl: string;
  videoUrl: string;
  eventYear: number | null;
  ringCount: number | null;
  likeCount: number;
  liked: number;
  creatorId: number;
};
