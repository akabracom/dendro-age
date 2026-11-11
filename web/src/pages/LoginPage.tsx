import { useState, type FormEvent } from "react";
import { Button, Form } from "react-bootstrap";
import { useNavigate } from "react-router-dom";
import axios from "axios";
import { loginResearcher } from "../api/auth";
import { useAppDispatch, useAppSelector } from "../store/hooks";
import { authFailed, authStarted, authSucceeded } from "../store/authSlice";

export function LoginPage() {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const error = useAppSelector((state) => state.auth.error);
  const [login, setLogin] = useState("shebe");
  const [password, setPassword] = useState("");

  async function handleSubmit(event: FormEvent) {
    event.preventDefault();
    dispatch(authStarted());
    try {
      const data = await loginResearcher(login.trim(), password);
      dispatch(authSucceeded({ login: data.login, token: data.access_token }));
      navigate("/mosaic");
    } catch (caught) {
      let message = "Не удалось войти";
      if (axios.isAxiosError(caught)) {
        const body = caught.response?.data?.message;
        if (typeof body === "string" && body) {
          message = body;
        }
      }
      dispatch(authFailed(message));
    }
  }

  return (
    <main className="page page--form">
      <header className="topbar">
        <h1 className="topbar__title">Вход</h1>
      </header>
      <Form className="stack-form" onSubmit={(event) => void handleSubmit(event)}>
        <Form.Group>
          <Form.Label>Логин</Form.Label>
          <Form.Control value={login} onChange={(event) => setLogin(event.target.value)} required />
        </Form.Group>
        <Form.Group>
          <Form.Label>Пароль</Form.Label>
          <Form.Control
            type="password"
            value={password}
            onChange={(event) => setPassword(event.target.value)}
            required
          />
        </Form.Group>
        {error ? <p className="form-error">{error}</p> : null}
        <Button className="filter-bar__submit" type="submit">
          Войти
        </Button>
      </Form>
    </main>
  );
}
