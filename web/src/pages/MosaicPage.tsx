import { useEffect, useState } from "react";
import { PatternCard } from "../components/PatternCard";
import { PatternFilter } from "../components/PatternFilter";
import { mockPatterns, publishedPatterns } from "../mock/patterns";
import type { DendroScalePattern } from "../types";

function yearFromDate(raw: string): number | null {
  if (!raw) {
    return null;
  }
  const parsed = new Date(raw + "T00:00:00");
  if (Number.isNaN(parsed.getTime())) {
    return null;
  }
  return parsed.getFullYear();
}

function applyFilter(
  items: DendroScalePattern[],
  dateFrom: string,
  dateTo: string,
  titleQuery: string,
): DendroScalePattern[] {
  const fromYear = yearFromDate(dateFrom);
  const toYear = yearFromDate(dateTo);
  const title = titleQuery.trim().toLowerCase();
  return items.filter((item) => {
    if (item.eventYear === null) {
      return false;
    }
    if (fromYear !== null && item.eventYear < fromYear) {
      return false;
    }
    if (toYear !== null && item.eventYear > toYear) {
      return false;
    }
    if (title && !item.title.toLowerCase().includes(title)) {
      return false;
    }
    return true;
  });
}

export function MosaicPage() {
  const [allPublished, setAllPublished] = useState<DendroScalePattern[]>([]);
  const [visible, setVisible] = useState<DendroScalePattern[]>([]);
  const [dateFrom, setDateFrom] = useState("1700-01-01");
  const [dateTo, setDateTo] = useState("2000-12-31");
  const [titleQuery, setTitleQuery] = useState("");

  useEffect(() => {
    const loaded = publishedPatterns(mockPatterns);
    setAllPublished(loaded);
    setVisible(applyFilter(loaded, "1700-01-01", "2000-12-31", ""));
  }, []);

  function handleApply() {
    setVisible(applyFilter(allPublished, dateFrom, dateTo, titleQuery));
  }

  return (
    <main className="page page--mosaic">
      <header className="topbar">
        <h1 className="topbar__title">Аномальные паттерны</h1>
      </header>
      <PatternFilter
        dateFrom={dateFrom}
        dateTo={dateTo}
        titleQuery={titleQuery}
        onDateFromChange={setDateFrom}
        onDateToChange={setDateTo}
        onTitleQueryChange={setTitleQuery}
        onApply={handleApply}
      />
      <section className="mosaic-grid">
        {visible.map((pattern) => (
          <PatternCard key={pattern.id} pattern={pattern} />
        ))}
        {visible.length === 0 ? <p className="mosaic-empty">Нет паттернов за выбранные условия.</p> : null}
      </section>
    </main>
  );
}
