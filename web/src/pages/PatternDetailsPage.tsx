import { useEffect } from "react";
import { Link, useParams } from "react-router-dom";
import { DEFAULT_IMAGE, DEFAULT_VIDEO } from "../mock/patterns";
import { proxyMedia } from "../api/media";
import { useAppDispatch, useAppSelector } from "../store/hooks";
import { fetchDendroScalePatterns } from "../store/patternSlice";

export function PatternDetailsPage() {
  const dispatch = useAppDispatch();
  const params = useParams();
  const list = useAppSelector((state) => state.patterns.list);
  const patternId = Number(params.id);
  const pattern = list.find((item) => item.id === patternId);

  useEffect(() => {
    if (list.length === 0) {
      void dispatch(fetchDendroScalePatterns());
    }
  }, [dispatch, list.length]);

  return (
    <main className="page page--details">
      <header className="topbar">
        <h1 className="topbar__title">Подробнее</h1>
      </header>
      {pattern ? (
        <article className="details">
          <img src={proxyMedia(pattern.imageUrl, DEFAULT_IMAGE)} alt={pattern.title} className="details__image" />
          <video
            className="details__video"
            src={proxyMedia(pattern.videoUrl, DEFAULT_VIDEO)}
            poster={proxyMedia(pattern.imageUrl, DEFAULT_IMAGE)}
            controls
            muted
          />
          <h2>{pattern.title}</h2>
          <p>{pattern.shortDescription}</p>
          <p>{pattern.fullDescription}</p>
          <p>
            Год начала {pattern.eventYear ?? "—"} · колец {pattern.ringCount ?? "—"} · лайки {pattern.likeCount}
          </p>
          <Link className="details__back" to="/mosaic">
            К плитке
          </Link>
        </article>
      ) : (
        <p className="feed-empty">Паттерн не найден.</p>
      )}
    </main>
  );
}
