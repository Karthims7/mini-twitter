import { useState } from "react";
import { api } from "../api";

export default function Signup({ onSignup }) {
  const [form, setForm] = useState({
    username: "",
    email: "",
    password: "",
  });
  const [err, setErr] = useState("");

  const submit = async (e) => {
    e.preventDefault();
    setErr("");

    const res = await api.post("/signup", form);

    if (res.id) {
      const login = await api.post("/login", {
        email: form.email,
        password: form.password,
      });
      if (login.token) onSignup(login.token);
    } else {
      setErr("Signup failed");
    }
  };

  return (
    <form onSubmit={submit}>
      <h2>Signup</h2>
      {err && <p style={{ color: "red" }}>{err}</p>}
      <input
        placeholder="Username"
        value={form.username}
        onChange={(e) => setForm({ ...form, username: e.target.value })}
      />
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
      <button>Signup</button>
    </form>
  );
}
