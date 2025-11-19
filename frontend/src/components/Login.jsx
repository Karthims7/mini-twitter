import { useState } from "react";
import { api } from "../api";

export default function Login({ onLogin }) {
  const [form, setForm] = useState({ email: "", password: "" });
  const [err, setErr] = useState("");

  const submit = async (e) => {
    e.preventDefault();
    setErr("");

    const res = await api.post("/login", form);
    if (res.token) onLogin(res.token);
    else setErr("Invalid credentials");
  };

  return (
    <form onSubmit={submit}>
      <h2>Login</h2>
      {err && <p style={{ color: "red" }}>{err}</p>}
      <input
        placeholder="Email"
        value={form.email}
        onChange={(e) => setForm({ ...form, email: e.target.value })}
      />
      <input
        placeholder="Password"
        type="password"
        value={form.password}
        onChange={(e) => setForm({ ...form, password: e.target.value })}
      />
      <button>Login</button>
    </form>
  );
}
