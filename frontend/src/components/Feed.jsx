import { useEffect, useState } from "react";
import { api } from "../api";

export default function Feed({ token }) {
  const [tweets, setTweets] = useState([]);

  const load = async () => {
    const res = await api.get("/feed", token);
    if (Array.isArray(res)) setTweets(res);
  };

  useEffect(() => {
    load();
  }, []);

  return (
    <div>
      <h2>Latest Tweets</h2>
      {tweets.length === 0 && <p>No tweets yet.</p>}
      {tweets.map((t) => (
        <div key={t.id} style={{ padding: "12px 0", borderBottom: "1px solid #eee" }}>
          <strong>@{t.username}</strong> —{" "}
          <small>{new Date(t.created_at).toLocaleString()}</small>
          <p>{t.content}</p>
        </div>
      ))}
    </div>
  );
}
