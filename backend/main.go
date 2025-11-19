package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type App struct {
	db        *sqlx.DB
	jwtSecret []byte
}

type User struct {
	ID           int64     `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

type Tweet struct {
	ID        int64     `db:"id" json:"id"`
	UserID    int64     `db:"user_id" json:"user_id"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	Username  string    `db:"username" json:"username"`
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@db:5432/twitter_dev?sslmode=disable"
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "supersecretchange"
	}

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal("db connect:", err)
	}
	defer db.Close()

	app := &App{db: db, jwtSecret: []byte(jwtSecret)}
	if err := app.ensureSchema(); err != nil {
		log.Fatal("schema:", err)
	}

	r := chi.NewRouter()
	r.Post("/signup", app.handleSignup)
	r.Post("/login", app.handleLogin)

	r.Group(func(r chi.Router) {
		r.Use(app.authMiddleware)
		r.Post("/tweets", app.handleCreateTweet)
		r.Get("/feed", app.handleGetFeed)
	})

	addr := ":8080"
	log.Println("listening", addr)
	http.ListenAndServe(addr, r)
}

/********** Schema **********/
func (a *App) ensureSchema() error {
	schema := `
CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,
  username VARCHAR(30) UNIQUE NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
CREATE TABLE IF NOT EXISTS tweets (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  content VARCHAR(280) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_tweets_created_at ON tweets(created_at DESC);
`
	_, err := a.db.Exec(schema)
	return err
}

/********** Helpers **********/
func (a *App) hashPassword(pw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(b), err
}
func (a *App) checkPassword(hash, pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
}

/********** Auth **********/
func (a *App) createToken(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(a.jwtSecret)
}

func (a *App) parseToken(tokenStr string) (int64, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return a.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}
	idFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("invalid user_id")
	}
	return int64(idFloat), nil
}

type key int

const userIDKey key = 0

func (a *App) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			http.Error(w, "missing auth", http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(h, "Bearer ")
		uid, err := a.parseToken(token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserID(r *http.Request) int64 {
	v := r.Context().Value(userIDKey)
	if v == nil {
		return 0
	}
	return v.(int64)
}

/********** Handlers **********/
func (a *App) handleSignup(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "missing fields", http.StatusBadRequest)
		return
	}
	pwHash, err := a.hashPassword(req.Password)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	var id int64
	err = a.db.QueryRow(
		"INSERT INTO users (username, email, password_hash) VALUES ($1,$2,$3) RETURNING id",
		req.Username, req.Email, pwHash).Scan(&id)
	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			http.Error(w, "username or email exists", http.StatusConflict)
			return
		}
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
}

func (a *App) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	var u User
	err := a.db.Get(&u, "SELECT id, password_hash FROM users WHERE email=$1", req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "invalid creds", http.StatusUnauthorized)
			return
		}
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	if a.checkPassword(u.PasswordHash, req.Password) != nil {
		http.Error(w, "invalid creds", http.StatusUnauthorized)
		return
	}
	token, err := a.createToken(u.ID)
	if err != nil {
		http.Error(w, "token error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (a *App) handleCreateTweet(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if len(req.Content) == 0 || len(req.Content) > 280 {
		http.Error(w, "content length invalid", http.StatusBadRequest)
		return
	}
	userID := getUserID(r)
	var id int64
	err := a.db.QueryRow("INSERT INTO tweets (user_id, content) VALUES ($1,$2) RETURNING id", userID, req.Content).Scan(&id)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
}

func (a *App) handleGetFeed(w http.ResponseWriter, r *http.Request) {
	// Simple feed: return latest 50 tweets across all users (for demo)
	rows, err := a.db.Queryx(`
SELECT t.id, t.user_id, t.content, t.created_at, u.username
FROM tweets t
JOIN users u ON u.id = t.user_id
ORDER BY t.created_at DESC
LIMIT 50
`)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	tweets := []Tweet{}
	for rows.Next() {
		var t Tweet
		if err := rows.StructScan(&t); err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		tweets = append(tweets, t)
	}
	json.NewEncoder(w).Encode(tweets)
}

