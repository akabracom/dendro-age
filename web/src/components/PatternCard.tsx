import { Card } from "react-bootstrap";
import { Link } from "react-router-dom";
import { DEFAULT_IMAGE, resolveMedia } from "../mock/patterns";
import type { DendroScalePattern } from "../types";

type PatternCardProps = {
  pattern: DendroScalePattern;
};

export function PatternCard({ pattern }: PatternCardProps) {
  const imageSrc = resolveMedia(pattern.imageUrl, DEFAULT_IMAGE);
  const year = pattern.eventYear === null ? "—" : String(pattern.eventYear);

  return (
    <Card className="pattern-card">
      <Link to={"/feed/" + String(pattern.id)} className="pattern-card__link">
        <Card.Img src={imageSrc} alt={pattern.title} className="pattern-card__image" />
        <Card.ImgOverlay className="pattern-card__shade">
          <Card.Title>{pattern.title}</Card.Title>
          <Card.Text>
            ♥ {pattern.likeCount} · {year}
          </Card.Text>
        </Card.ImgOverlay>
      </Link>
    </Card>
  );
}
