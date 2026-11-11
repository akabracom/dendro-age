import { Button, Card } from "react-bootstrap";
import { Link } from "react-router-dom";
import { DEFAULT_IMAGE } from "../mock/patterns";
import { proxyMedia } from "../api/media";
import type { DendroScalePattern } from "../types";

type PatternCardProps = {
  pattern: DendroScalePattern;
  onDelete?: (id: number) => void;
};

export function PatternCard({ pattern, onDelete }: PatternCardProps) {
  const imageSrc = proxyMedia(pattern.imageUrl, DEFAULT_IMAGE);
  const year = pattern.eventYear === null ? "—" : String(pattern.eventYear);

  return (
    <Card className="pattern-card">
      <Link to={"/feed/" + String(pattern.id)} className="pattern-card__link">
        <Card.Img src={imageSrc} alt={pattern.title} className="pattern-card__image" />
        <Card.ImgOverlay className="pattern-card__shade">
          <Card.Title>{pattern.title}</Card.Title>
          <Card.Text>
            ♥ {pattern.likeCount} · {year}
            {pattern.mine === 1 ? " · моя" : ""}
          </Card.Text>
        </Card.ImgOverlay>
      </Link>
      {pattern.mine === 1 && onDelete ? (
        <Button className="pattern-card__delete" type="button" onClick={() => onDelete(pattern.id)}>
          Удалить
        </Button>
      ) : null}
    </Card>
  );
}
