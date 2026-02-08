# Memory Quest

Memory Quest is a full-stack web application for creating and playing memory-based quiz decks. It combines a Go backend API with a modern React frontend, allowing users to build decks, answer quiz questions, track results, and compete via a leaderboard.

The backend follows a modular, layered structure to keep business logic clean, testable, and easy to extend. The frontend is built as a separate SPA that communicates with the API.

From a user perspective, Memory Quest is designed around typical learning and quiz flows: users can create an account, log in, create decks, play quizzes, manage their progress, and interact with features that update dynamically. The UI is modern and responsive, backed by API-driven data.

---

## What Users Can Do

- Register and log in to a personal account  
- Create and manage quiz decks  
- Add, edit, and delete cards inside a deck  
- Play quizzes generated from their decks  
- Receive immediate feedback on answers  
- View quiz results and history  
- Compete via a leaderboard  
- Use a responsive UI with dynamic updates  

---

## Table of Contents

- Features  
- Project Structure  
- Tech Stack  
- Getting Started  
- Prerequisites  
- Installation  
- Configuration (Environment Variables)  
- Run  
- Development  
- Troubleshooting  
- Roadmap  

---

## Features

Server-rendered and API-driven UI components with a modern frontend  
Organized backend structure: api, middleware, auth, models, db  
Database layer (MySQL)  
Authentication and authorization (JWT)  
Common app flows (examples):

authentication + JWT handling  
user actions (deck management, quiz attempts, results)  
leaderboard and score tracking  
quiz generation logic  

---

## Project Structure

.
├── internal/
│   ├── api/           # Route handlers / controllers  
│   ├── auth/          # JWT + password logic  
│   ├── db/            # Database connection + queries  
│   ├── middleware/    # Auth / logging / etc.  
│   ├── models/        # Domain models / DTOs  
│   └── main.go        # App entrypoint  
│
├── memory-quest-frontend/
│   ├── src/           # React application source  
│   ├── public/  
│   ├── package.json  
│   └── config files  
│
├── go.mod  
└── go.sum  

---

## Tech Stack

Backend: Go (Golang), Gin  
Database: MySQL  
Authentication: JWT  
Frontend: React, Vite, JavaScript, HTML, CSS  

---

## Getting Started

### Prerequisites

Go installed (recommended: latest stable)  
Node.js (18+ recommended)  
A MySQL database instance  

---

## Installation

git clone <your-repo-url>  
cd memory-quest-backend  
go mod download  

cd memory-quest-frontend  
npm install  

---

## Configuration (Environment Variables)

Create a .env file (or export env vars in your shell). Rename these to match what your code expects (look in main.go and db/):

# Server  
PORT=8080  

# Database  
DB_HOST=localhost  
DB_PORT=3306  
DB_USER=root  
DB_PASSWORD=your_password  
DB_NAME=memoryquest  

# Auth  
JWT_SECRET=replace_me  

---

## Run

Backend:

go run internal/main.go  

Then open:

http://localhost:8080  

Frontend:

cd memory-quest-frontend  
npm run dev  

---

## Development

Common commands:

go fmt ./...  
go test ./...  
go vet ./...  

Frontend:

npm run lint  
npm run build  

---

## Troubleshooting

1) Login always fails even with correct password  
Make sure passwords are being hashed during registration and compared using the same hashing function during login.

2) API returns 401 Unauthorized on protected routes  
Verify the JWT is included in the Authorization header as `Bearer <token>` and that the token has not expired.

3) Decks or cards not saving  
Check database migrations / table structure and ensure the DB user has INSERT and UPDATE permissions.

4) Slow responses from quiz endpoints  
Look for missing database indexes on foreign keys like user_id or deck_id.

5) Frontend shows network error  
Confirm the frontend API base URL matches the backend port and protocol (http vs https).

6) Web app works locally but not on deployment  
Check environment variables in the hosting environment and ensure the database host is reachable from the server.

---

## Roadmap

Public deck sharing  
Multiplayer quiz mode  
Timer-based challenges  
Improved statistics dashboard  
Card images support  
Admin tools  
