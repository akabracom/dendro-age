import { useEffect } from "react";
import { PatternCard } from "../components/PatternCard";
import { PatternFilter } from "../components/PatternFilter";
import { useAppDispatch, useAppSelector } from "../store/hooks";
import { setDateFrom, setDateTo, setTitle } from "../store/filterSlice";
import { deleteDendroScalePattern, fetchDendroScalePatterns } from "../store/patternSlice";

export function MosaicPage() {
  const dispatch = useAppDispatch();
  const filter = useAppSelector((state) => state.filter);
  const list = useAppSelector((state) => state.patterns.list);
  const error = useAppSelector((state) => state.patterns.error);
  const source = useAppSelector((state) => state.patterns.source);

  useEffect(() => {
    void dispatch(fetchDendroScalePatterns());
  }, [dispatch]);

  return (
    <main className="page page--mosaic">
      <header className="topbar">
        <h1 className="topbar__title">Аномальные паттерны</h1>
      </header>
      <PatternFilter
        dateFrom={filter.dateFrom}
        dateTo={filter.dateTo}
        titleQuery={filter.title}
        onDateFromChange={(value) => dispatch(setDateFrom(value))}
        onDateToChange={(value) => dispatch(setDateTo(value))}
        onTitleQueryChange={(value) => dispatch(setTitle(value))}
        onApply={() => void dispatch(fetchDendroScalePatterns())}
      />
      {source === "mock" ? <p className="source-note">Сервис недоступен, показаны mock-паттерны гостя.</p> : null}
      {error ? <p className="source-note">{error}</p> : null}
      <section className="mosaic-grid">
        {list.map((pattern) => (
          <PatternCard
            key={pattern.id}
            pattern={pattern}
            onDelete={(id) => void dispatch(deleteDendroScalePattern(id))}
          />
        ))}
        {list.length === 0 ? <p className="mosaic-empty">Нет паттернов за выбранные условия.</p> : null}
      </section>
    </main>
  );
}
