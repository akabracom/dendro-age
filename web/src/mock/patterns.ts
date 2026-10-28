import type { DendroScalePattern } from "../types";

export const DEFAULT_IMAGE = "/defaults/pattern_default.svg";
export const DEFAULT_VIDEO = "/defaults/pattern_default.mp4";

export const mockPatterns: DendroScalePattern[] = [
  {
    id: 1,
    title: "Морозное кольцо 1695",
    shortDescription: "Узкие кольца холодного сезона.",
    fullDescription: "Кластер из четырёх колец фиксирует аномальный холод и служит репером дендрошкалы.",
    status: "published",
    imageUrl: "",
    videoUrl: "",
    eventYear: 1695,
    ringCount: 4,
    likeCount: 3,
    liked: 0,
    creatorId: 1,
  },
  {
    id: 2,
    title: "Пожарные рубцы 1815",
    shortDescription: "След летнего пожара на спиле.",
    fullDescription: "Резкий переход плотности колец после низового пожара в сосняке.",
    status: "published",
    imageUrl: "/defaults/pattern_default.svg",
    videoUrl: "/defaults/pattern_default.mp4",
    eventYear: 1815,
    ringCount: 7,
    likeCount: 1,
    liked: 0,
    creatorId: 2,
  },
  {
    id: 3,
    title: "Засуха 1921",
    shortDescription: "Серия тонких колец без поздней древесины.",
    fullDescription: "Пять последовательных засушливых лет на восточноевропейской шкале.",
    status: "published",
    imageUrl: "",
    videoUrl: "",
    eventYear: 1921,
    ringCount: 5,
    likeCount: 6,
    liked: 0,
    creatorId: 1,
  },
  {
    id: 4,
    title: "Черновик без медиа",
    shortDescription: "",
    fullDescription: "",
    status: "draft",
    imageUrl: "",
    videoUrl: "",
    eventYear: null,
    ringCount: null,
    likeCount: 0,
    liked: 0,
    creatorId: 1,
  },
];

export function resolveMedia(url: string, fallback: string): string {
  if (!url.trim()) {
    return fallback;
  }
  return url;
}

export function publishedPatterns(items: DendroScalePattern[]): DendroScalePattern[] {
  return items.filter((item) => item.status === "published");
}
