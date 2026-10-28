type LockedPageProps = {
  title: string;
};

export function LockedPage({ title }: LockedPageProps) {
  return (
    <main className="page page--locked">
      <header className="topbar">
        <h1 className="topbar__title">{title}</h1>
      </header>
      <p className="locked-text">Страница недоступна гостю. В лабораторной 5 открыты только лента и плитка.</p>
    </main>
  );
}
