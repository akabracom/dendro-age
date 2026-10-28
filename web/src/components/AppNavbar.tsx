import { Nav } from "react-bootstrap";
import { NavLink } from "react-router-dom";

export function AppNavbar() {
  return (
    <nav className="tabbar">
      <Nav className="tabbar__nav">
        <Nav.Item>
          <NavLink className={({ isActive }) => "tabbar__item" + (isActive ? " is-active" : "")} to="/feed">
            Лента
          </NavLink>
        </Nav.Item>
        <Nav.Item>
          <span className="tabbar__item is-locked">Добавить</span>
        </Nav.Item>
        <Nav.Item>
          <NavLink className={({ isActive }) => "tabbar__item" + (isActive ? " is-active" : "")} to="/mosaic">
            Плитка
          </NavLink>
        </Nav.Item>
        <Nav.Item>
          <span className="tabbar__item is-locked">Вход</span>
        </Nav.Item>
      </Nav>
    </nav>
  );
}
