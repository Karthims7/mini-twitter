Sure Karthi — here is the **complete high-level README.md** you asked for, formatted as a **copy-paste code block**.

Create the file inside:

```
~/mini-twitter/README.md
```

---

# ✅ **COPY THIS EXACT README.md**

```markdown
# 🐦 Mini Twitter – Full Stack App (Go + Postgres + React + Vite + Docker)

A complete full-stack “Mini Twitter” clone built with:

- **Go** (backend API)
- **PostgreSQL** (database)
- **React + Vite** (frontend)
- **Docker Compose** (production-ready orchestration)
- **EC2 Deployment** (real cloud hosting)

This README explains **how the whole system works**, **how to develop**, and **how to deploy** it anywhere.

---

# 🚀 Architecture Overview

```

mini-twitter/
│
├── backend/          → Go REST API (Chi router, JWT auth, SQLX ORM)
│   ├── main.go
│   ├── go.mod
│   ├── go.sum
│   └── Dockerfile
│
├── frontend/         → React + Vite UI
│   ├── index.html
│   ├── package.json
│   └── src/
│
└── docker-compose.yml → Runs API + DB services together

````

---

# ⚙️ Backend – Go API

### Features
- JWT-based login/signup
- Tweet creation
- Feed listing (with username)
- Secure password hashing (bcrypt)
- PostgreSQL integration via `sqlx`
- Dockerized build (multi-stage)

### Development (Locally)
```bash
cd backend
go run main.go
````

### Build inside Docker

```bash
docker compose build api
```

---

# 🗄️ Database – PostgreSQL

### Location

Handled automatically by `docker-compose.yml`.

### Credentials (dev defaults)

```
user: postgres
password: postgres
database: twitter_dev
host (from backend): db
```

---

# 🎨 Frontend – React + Vite

### Development Mode

```bash
cd frontend
VITE_API_URL=http://localhost:8080 npm run dev
```

### Build for Production

```bash
npm run build
```

---

# 🐳 Docker Compose – Full Stack (API + DB)

### Start services:

```bash
docker compose up -d
```

### Stop services:

```bash
docker compose down
```

### View running containers:

```bash
docker ps
```

---

# 🌍 EC2 Deployment Guide (Working Setup)

1. Launch Ubuntu EC2.
2. Install Docker + Docker Compose.
3. Clone or SCP project folder.
4. Run:

   ```bash
   docker compose up -d --build
   ```
5. API will run on:

   ```
   http://<EC2-public-IP>:8080
   ```

Frontend runs separately using Vite (dev mode), or you deploy it via Nginx if desired.

---

# 🔑 API Examples

### Signup

```bash
curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","email":"alice@example.com","password":"password"}'
```

### Login (returns token)

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"password"}'
```

### Create Tweet

```bash
curl -X POST http://localhost:8080/tweets \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{"content":"Hello world!"}'
```

### View Feed

```bash
curl -H "Authorization: Bearer <TOKEN>" http://localhost:8080/feed
```

---

# 📁 .gitignore Included

Automatically ignores:

* node_modules
* dist
* docker volumes
* build outputs (`server`)
* .env files
* OS junk files
* caches

---

# ✔️ Final Notes

This project:

* Runs fully on Docker.
* Works locally and on EC2.
* Backend and DB are fully production-ready.
* Frontend can run via Vite dev or built + hosted via Nginx.
