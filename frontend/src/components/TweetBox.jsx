import { useState } from "react";
import { api } from "../api";

export default function TweetBox({ token, onTweet }) {
  const [content, setContent] = useState("");
  const [err, setErr] = useState("");

  const submit = async (e) => {
    e.preventDefault();
    setErr("");

    const res = await api.post("/tweets", { content }, token);
    if (res.id) {
      setContent("");
      if (onTweet) onTweet();
    } else {
      setErr("Failed to tweet");
    }
  };

  return (
    <form onSubmit={submit}>
      {err && <p style={{ color: "red" }}>{err}</p>}
      <textarea
        placeholder="What's happening?"
        value={content}
        onChange={(e) => setContent(e.target.value)}
      />
      <button>Tweet</button>
    </form>
  );
}
