import { Navigate, Route, Routes } from "react-router-dom";
import { AppNavbar } from "./components/AppNavbar";
import { FeedPage } from "./pages/FeedPage";
import { LockedPage } from "./pages/LockedPage";
import { MosaicPage } from "./pages/MosaicPage";

function App() {
  return (
    <div className="phone">
      <Routes>
        <Route path="/" element={<Navigate to="/feed" replace />} />
        <Route path="/feed" element={<FeedPage />} />
        <Route path="/feed/:id" element={<FeedPage />} />
        <Route path="/mosaic" element={<MosaicPage />} />
        <Route path="/add" element={<LockedPage title="Добавление" />} />
        <Route path="/login" element={<LockedPage title="Вход" />} />
      </Routes>
      <AppNavbar />
    </div>
  );
}

export default App;
