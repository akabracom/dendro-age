import { useEffect, useState, type FormEvent } from "react";
import { Button, Form } from "react-bootstrap";
import { Navigate, useNavigate } from "react-router-dom";
import { DEFAULT_IMAGE, DEFAULT_VIDEO } from "../mock/patterns";
import { proxyMedia } from "../api/media";
import { useAppDispatch, useAppSelector } from "../store/hooks";
import {
  createDendroScalePattern,
  fetchDendroScalePatternDraft,
  publishDendroScalePattern,
} from "../store/patternSlice";

export function AddPage() {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const isUser = useAppSelector((state) => state.auth.status === "user");
  const draft = useAppSelector((state) => state.patterns.draft);
  const error = useAppSelector((state) => state.patterns.error);
  const [title, setTitle] = useState("");
  const [image, setImage] = useState<File | undefined>();
  const [video, setVideo] = useState<File | undefined>();
  const [shortDescription, setShortDescription] = useState("");
  const [fullDescription, setFullDescription] = useState("");
  const [eventYear, setEventYear] = useState("1695");
  const [ringCount, setRingCount] = useState("4");

  useEffect(() => {
    if (isUser) {
      void dispatch(fetchDendroScalePatternDraft());
    }
  }, [dispatch, isUser]);

  if (!isUser) {
    return <Navigate to="/login" replace />;
  }

  async function handleCreate(event: FormEvent) {
    event.preventDefault();
    const result = await dispatch(createDendroScalePattern({ title: title.trim(), image, video }));
    if (createDendroScalePattern.fulfilled.match(result)) {
      setTitle("");
      setImage(undefined);
      setVideo(undefined);
    }
  }

  async function handlePublish(event: FormEvent) {
    event.preventDefault();
    const result = await dispatch(
      publishDendroScalePattern({
        shortDescription: shortDescription.trim(),
        fullDescription: fullDescription.trim(),
        eventYear: Number(eventYear),
        ringCount: Number(ringCount),
      }),
    );
    if (publishDendroScalePattern.fulfilled.match(result)) {
      navigate("/mosaic");
    }
  }

  return (
    <main className="page page--form">
      <header className="topbar">
        <h1 className="topbar__title">Добавление</h1>
      </header>
      {error ? <p className="form-error">{error}</p> : null}
      {draft ? (
        <Form className="stack-form" onSubmit={(event) => void handlePublish(event)}>
          <p className="draft-title">{draft.title}</p>
          <video
            className="draft-video"
            src={proxyMedia(draft.videoUrl, DEFAULT_VIDEO)}
            poster={proxyMedia(draft.imageUrl, DEFAULT_IMAGE)}
            controls
            muted
          />
          <Form.Group>
            <Form.Label>Краткое описание</Form.Label>
            <Form.Control value={shortDescription} onChange={(event) => setShortDescription(event.target.value)} required />
          </Form.Group>
          <Form.Group>
            <Form.Label>Полное описание</Form.Label>
            <Form.Control
              as="textarea"
              rows={3}
              value={fullDescription}
              onChange={(event) => setFullDescription(event.target.value)}
            />
          </Form.Group>
          <Form.Group>
            <Form.Label>Год начала</Form.Label>
            <Form.Control value={eventYear} onChange={(event) => setEventYear(event.target.value)} required />
          </Form.Group>
          <Form.Group>
            <Form.Label>Число колец</Form.Label>
            <Form.Control value={ringCount} onChange={(event) => setRingCount(event.target.value)} required />
          </Form.Group>
          <Button className="filter-bar__submit" type="submit">
            Опубликовать
          </Button>
        </Form>
      ) : (
        <Form className="stack-form" onSubmit={(event) => void handleCreate(event)}>
          <Form.Group>
            <Form.Label>Название</Form.Label>
            <Form.Control value={title} onChange={(event) => setTitle(event.target.value)} required />
          </Form.Group>
          <Form.Group>
            <Form.Label>Видео</Form.Label>
            <Form.Control
              type="file"
              accept="video/*"
              onChange={(event) => setVideo((event.target as HTMLInputElement).files?.[0])}
            />
          </Form.Group>
          <Form.Group>
            <Form.Label>Фото</Form.Label>
            <Form.Control
              type="file"
              accept="image/*"
              onChange={(event) => setImage((event.target as HTMLInputElement).files?.[0])}
            />
          </Form.Group>
          <Button className="filter-bar__submit" type="submit">
            Далее
          </Button>
        </Form>
      )}
    </main>
  );
}
