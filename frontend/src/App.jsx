import { useState } from "react";
import Login from "./components/Login";
import Signup from "./components/Signup";
import Feed from "./components/Feed";
import TweetBox from "./components/TweetBox";

export default function App() {
  const [token, setToken] = useState(localStorage.getItem("token") || "");
  const [page, setPage] = useState("feed");

  const logout = () => {
    setToken("");
    localStorage.removeItem("token");
  };

  const onLogin = (t) => {
    setToken(t);
    localStorage.setItem("token", t);
    setPage("feed");
  };

  return (
    <div className="wrapper">
      <header>
        <h1>Mini Twitter</h1>
        <nav>
          {token ? (
            <>
              <button onClick={() => setPage("feed")}>Feed</button>
              <button onClick={logout}>Logout</button>
            </>
          ) : (
            <>
              <button onClick={() => setPage("login")}>Login</button>
              <button onClick={() => setPage("signup")}>Signup</button>
            </>
          )}
        </nav>
      </header>

      <main>
        {!token && page === "login" && <Login onLogin={onLogin} />}
        {!token && page === "signup" && <Signup onSignup={onLogin} />}

        {token && (
          <>
            <TweetBox token={token} onTweet={() => setPage("feed")} />
            <Feed token={token} />
          </>
        )}

        {!token && page === "feed" && <p>Login to view tweets.</p>}
      </main>

      <footer>
        <small>React + Go + PG Demo</small>
      </footer>
    </div>
  );
}
