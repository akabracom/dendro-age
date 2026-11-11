import { useEffect } from "react";
import { Button } from "react-bootstrap";
import { DEFAULT_IMAGE, DEFAULT_VIDEO } from "../mock/patterns";
import { proxyMedia } from "../api/media";
import { useAppDispatch, useAppSelector } from "../store/hooks";
import { fetchDendroScalePatternFeed, likeDendroScalePattern } from "../store/patternSlice";

export function FeedPage() {
  const dispatch = useAppDispatch();
  const feed = useAppSelector((state) => state.patterns.feed);
  const error = useAppSelector((state) => state.patterns.error);
  const isUser = useAppSelector((state) => state.auth.status === "user");

  useEffect(() => {
    void dispatch(fetchDendroScalePatternFeed());
  }, [dispatch]);

  if (!feed) {
    return (
      <main className="page page--feed">
        <p className="feed-empty">{error || "Нет опубликованных паттернов."}</p>
      </main>
    );
  }

  const imageSrc = proxyMedia(feed.imageUrl, DEFAULT_IMAGE);
  const videoSrc = proxyMedia(feed.videoUrl, DEFAULT_VIDEO);
  const year = feed.eventYear === null ? "—" : String(feed.eventYear);
  const rings = feed.ringCount === null ? "—" : String(feed.ringCount);

  return (
    <main className="page page--feed">
      <video className="feed__video" src={videoSrc} poster={imageSrc} autoPlay muted loop playsInline />
      <div className="feed__overlay">
        <p className="feed__logo">Dendro Age</p>
        <aside className="feed__rail">
          <div>
            <span className="feed__label">Год начала</span>
            <strong>{year}</strong>
          </div>
          <div>
            <span className="feed__label">Кольца</span>
            <strong>{rings}</strong>
          </div>
          <div>
            <span className="feed__label">Лайки</span>
            <strong>{feed.likeCount}</strong>
          </div>
          <div>
            <span className="feed__label">Мой лайк</span>
            <strong>{feed.liked}</strong>
          </div>
          {isUser ? (
            <Button
              className="feed__next"
              type="button"
              onClick={() => void dispatch(likeDendroScalePattern({ id: feed.id, like: feed.liked === 1 ? 0 : 1 }))}
            >
              Лайк
            </Button>
          ) : null}
          <Button
            className="feed__next"
            type="button"
            onClick={() => void dispatch(fetchDendroScalePatternFeed({ id: feed.id, next: true }))}
          >
            Следующий
          </Button>
        </aside>
        <div className="feed__bottom">
          <h1>{feed.title}</h1>
          <p>{feed.shortDescription}</p>
          <p className="feed__full">{feed.fullDescription}</p>
        </div>
      </div>
    </main>
  );
}
