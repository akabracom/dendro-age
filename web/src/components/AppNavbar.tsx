import { Nav } from "react-bootstrap";
import { NavLink, useNavigate } from "react-router-dom";
import { logoutResearcher } from "../api/auth";
import { useAppDispatch, useAppSelector } from "../store/hooks";
import { authCleared, authFailed, authStarted } from "../store/authSlice";
import { resetFilter } from "../store/filterSlice";

export function AppNavbar() {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const isUser = useAppSelector((state) => state.auth.status === "user");

  async function handleLogout() {
    dispatch(authStarted());
    try {
      await logoutResearcher();
    } catch {
      dispatch(authFailed("Не удалось выйти"));
    }
    dispatch(resetFilter());
    dispatch(authCleared());
    navigate("/feed");
  }

  return (
    <nav className="tabbar">
      <Nav className="tabbar__nav">
        <Nav.Item>
          <NavLink className={({ isActive }) => "tabbar__item" + (isActive ? " is-active" : "")} to="/feed">
            Лента
          </NavLink>
        </Nav.Item>
        <Nav.Item>
          {isUser ? (
            <NavLink className={({ isActive }) => "tabbar__item" + (isActive ? " is-active" : "")} to="/add">
              Добавить
            </NavLink>
          ) : (
            <button className="tabbar__item tabbar__button" type="button" onClick={() => navigate("/login")}>
              Добавить
            </button>
          )}
        </Nav.Item>
        <Nav.Item>
          <NavLink className={({ isActive }) => "tabbar__item" + (isActive ? " is-active" : "")} to="/mosaic">
            Плитка
          </NavLink>
        </Nav.Item>
        <Nav.Item>
          {isUser ? (
            <button className="tabbar__item tabbar__button" type="button" onClick={() => void handleLogout()}>
              Выход
            </button>
          ) : (
            <NavLink className={({ isActive }) => "tabbar__item" + (isActive ? " is-active" : "")} to="/login">
              Вход
            </NavLink>
          )}
        </Nav.Item>
      </Nav>
    </nav>
  );
}
