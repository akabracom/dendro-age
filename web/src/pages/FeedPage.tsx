import { useEffect, useState } from "react";
import { Button } from "react-bootstrap";
import { useNavigate, useParams } from "react-router-dom";
import { DEFAULT_IMAGE, DEFAULT_VIDEO, mockPatterns, publishedPatterns, resolveMedia } from "../mock/patterns";
import type { DendroScalePattern } from "../types";

export function FeedPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [items, setItems] = useState<DendroScalePattern[]>([]);
  const [index, setIndex] = useState(0);

  useEffect(() => {
    setItems(publishedPatterns(mockPatterns));
  }, []);

  useEffect(() => {
    if (items.length === 0) {
      return;
    }
    if (!id) {
      setIndex(0);
      return;
    }
    const found = items.findIndex((item) => item.id === Number(id));
    setIndex(found >= 0 ? found : 0);
  }, [id, items]);

  if (items.length === 0) {
    return (
      <main className="page page--feed">
        <p className="feed-empty">Нет опубликованных паттернов.</p>
      </main>
    );
  }

  const pattern = items[index];
  const imageSrc = resolveMedia(pattern.imageUrl, DEFAULT_IMAGE);
  const videoSrc = resolveMedia(pattern.videoUrl, DEFAULT_VIDEO);
  const year = pattern.eventYear === null ? "—" : String(pattern.eventYear);
  const rings = pattern.ringCount === null ? "—" : String(pattern.ringCount);

  function handleNext() {
    const nextIndex = (index + 1) % items.length;
    navigate("/feed/" + String(items[nextIndex].id));
  }

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
            <strong>{pattern.likeCount}</strong>
          </div>
          <div>
            <span className="feed__label">Мой лайк</span>
            <strong>{pattern.liked}</strong>
          </div>
          <Button className="feed__next" onClick={handleNext}>
            Следующий
          </Button>
        </aside>
        <div className="feed__bottom">
          <h1>{pattern.title}</h1>
          <p>{pattern.shortDescription}</p>
          <p className="feed__full">{pattern.fullDescription}</p>
        </div>
      </div>
    </main>
  );
}
