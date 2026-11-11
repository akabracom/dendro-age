import { Navigate, Route, Routes } from "react-router-dom";
import { AppNavbar } from "./components/AppNavbar";
import { BusyOverlay } from "./components/BusyOverlay";
import { AddPage } from "./pages/AddPage";
import { FeedPage } from "./pages/FeedPage";
import { LoginPage } from "./pages/LoginPage";
import { MosaicPage } from "./pages/MosaicPage";

function App() {
  return (
    <div className="phone">
      <Routes>
        <Route path="/" element={<Navigate to="/feed" replace />} />
        <Route path="/feed" element={<FeedPage />} />
        <Route path="/feed/:id" element={<FeedPage />} />
        <Route path="/mosaic" element={<MosaicPage />} />
        <Route path="/add" element={<AddPage />} />
        <Route path="/login" element={<LoginPage />} />
      </Routes>
      <AppNavbar />
      <BusyOverlay />
    </div>
  );
}

export default App;
