import React from "react";
import { createRoot } from "react-dom/client";
import App from "./App";
import "./styles.css";

const rootEl = document.getElementById("root");

function hideSplash() {
  const s = document.getElementById("splash");
  if (!s) return;
  s.classList.add("splash-hidden");
  setTimeout(()=> s.remove(), 420);
}

if (rootEl) {
  createRoot(rootEl).render(
    <React.StrictMode>
      <App />
    </React.StrictMode>
  );
  setTimeout(hideSplash, 80);
}

