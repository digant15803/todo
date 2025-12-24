# 📝 TODO API – Golang

A minimal REST API built in **Go** for managing TODO items.  
This project demonstrates clean API design, concurrency-safe in-memory storage, proper separation of concerns, and testable business logic.

---

## 🚀 Features

- Create TODO items
- List all TODO items
- Mark a TODO as completed (without deleting)
- Delete a TODO
- In-memory storage (no database)
- Thread-safe operations using mutex
- Request validation
- Unit tests (including negative cases)

---

## 📌 API Endpoints

### ➕ Create a TODO
```
POST /todos
POST /todos/
```

**Request Body**
```json
{
  "title": "Buy groceries"
}
```

**Response**
```json
{
  "id": {"uuid"},
  "title": "Buy groceries",
  "completed": false
}
```

---

### 📋 List All TODOs
```
GET /todos
GET /todos/
```

**Response**
```json
[
  {
    "id": {"uuid"},
    "title": "Buy groceries",
    "completed": false
  }
]
```

---

### ✅ Mark TODO as Completed
```
PUT /todos/{id}/complete
```

- No request body required
- Sets `completed = true`

---

### ❌ Delete a TODO
```
DELETE /todos/{id}
```

---

## 🗂 Project Structure

```
todo/
├── main.go
├── internal/
│   ├── handler/
│   └── todo/
├── go.mod
├── go.sum
└── README.md
```

---

## ⚙️ Setup Instructions

### Clone the Repository
```bash
git clone git@github.com:digant15803/todo.git
cd todo
```

### Install Dependencies
```bash
go mod tidy
```

---

## ▶️ Run the Server

```bash
go run main.go
```

Server starts at:
```
http://localhost:8080
```

---

## 🧪 Running Tests

```bash
go test ./...
```
---

## 🧠 Design Decisions & Notes

### In-Memory Storage
- No persistence across restarts
- Simplifies the architecture

### Mutex Usage
The service uses a `sync.Mutex` to protect the shared in-memory map.
This avoids race conditions since HTTP handlers in Go run concurrently.

---

## 📌 Assumptions

- No authentication
- Use of library to generate UUID (id) for any todo Task
- Completed TODOs cannot be reverted

