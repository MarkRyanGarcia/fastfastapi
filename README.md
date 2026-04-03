# fastfastapi

A CLI tool that scaffolds a production-ready FastAPI project with an interactive TUI. Choose your database, ORM, auth provider, package manager, and optional Docker/Redis support, and get a fully structured project in seconds.

## Installation

### via pip (recommended)

```bash
pip install fastfastapi
```

### via pipx

```bash
pipx install fastfastapi
```

### via Go

```bash
go install github.com/markryangarcia/fastfastapi@latest
go install github.com/markryangarcia/fastfastapi/cmd/ffa@latest
```

Make sure `$GOPATH/bin` (or `$HOME/go/bin`) is in your `PATH`.

## Prerequisites

- Python 3.8+ (for pip/pipx install)
- [Go](https://go.dev/dl/) 1.21+ (only if installing via Go)
- `pipenv` (optional, only if you choose it during setup)
- Docker (optional, only if you choose Docker support)

## Usage

```bash
# Create a new project in a new directory
fastfastapi

# Scaffold into the current directory
fastfastapi .

# Pass a project name directly (skips the name prompt)
fastfastapi my-api

# ffa is an alias for fastfastapi
ffa my-api
```

The TUI walks you through:

1. Project name
2. Database — `PostgreSQL (SQLAlchemy)` or `MongoDB (PyMongo)`
3. ORM — `SQLAlchemy`, `SQLModel`, or `FastCRUD` (PostgreSQL only)
4. Auth provider — `None`, `Clerk`, or `AWS Cognito`
5. Package manager — `Pipenv` or `requirements.txt`
6. Docker support — generates `Dockerfile`, `docker-compose.yml`, and `.dockerignore`
7. Redis caching — generates `app/core/cache.py` with Redis integration
8. Install & start — run setup automatically or skip

## Generated Project Structure

```
my-api/
├── app/
│   ├── api/v1/
│   │   └── routers/
│   │       ├── users.py
│   │       └── items.py
│   ├── core/
│   │   ├── config.py
│   │   ├── security.py
│   │   └── cache.py          # Redis only
│   ├── db/
│   │   ├── base.py           # SQLAlchemy only
│   │   └── session.py
│   ├── models/
│   │   ├── user.py
│   │   └── item.py
│   ├── schemas/
│   │   ├── user.py
│   │   └── item.py
│   ├── services/
│   │   ├── user_service.py
│   │   └── item_service.py
│   ├── utils/
│   │   ├── pagination.py
│   │   ├── responses.py
│   │   └── exceptions.py
│   └── main.py
├── migrations/               # PostgreSQL only (Alembic)
│   └── versions/
├── tests/
│   ├── test_users.py
│   └── test_items.py
├── conftest.py
├── .env
├── .gitignore
├── alembic.ini               # PostgreSQL only
├── Dockerfile                # Docker only
├── docker-compose.yml        # Docker only
├── .dockerignore             # Docker only
├── requirements.txt          # if not using pipenv
└── Pipfile                   # if using pipenv
```

## ORM Options (PostgreSQL only)

| Option | Description |
|---|---|
| `SQLAlchemy` | Classic SQLAlchemy Core + ORM with Alembic migrations |
| `SQLModel` | SQLModel (SQLAlchemy + Pydantic) — models double as schemas |
| `FastCRUD` | FastCRUD on top of SQLModel for auto-generated CRUD endpoints |

MongoDB always uses PyMongo directly — no ORM prompt.

## Auth Providers

| Option | Description |
|---|---|
| `None` | Custom JWT auth via `app/core/security.py` |
| `Clerk` | Clerk JWT verification wired into the security module |
| `AWS Cognito` | AWS Cognito JWT verification wired into the security module |

## Setting Up the Generated Project

### 1. Configure environment variables

Edit the generated `.env` file:

**PostgreSQL:**
```env
APP_NAME="my-api"
DATABASE_URL="postgresql://user:password@localhost:5432/my-api"
```

**MongoDB:**
```env
APP_NAME="my-api"
MONGODB_URL="mongodb://localhost:27017"
MONGODB_DB="my-api"
```

**Redis (if enabled):**
```env
REDIS_URL="redis://localhost:6379"
```

### 2. Install dependencies

**With Docker:**
```bash
cd my-api
docker compose up --build
```

**With pipenv:**
```bash
cd my-api
pipenv install --dev
pipenv shell
fastapi dev app
```

**With requirements.txt:**
```bash
cd my-api
python3 -m venv .venv
source .venv/bin/activate  # Windows: .venv\Scripts\activate
pip install -r requirements.txt
fastapi dev app
```

The API will be available at `http://localhost:8000`.
Interactive docs at `http://localhost:8000/docs`.

### 3. Run tests

```bash
pytest
```

## Notes

- `q` or `Ctrl+C` at any point in the TUI cancels generation without writing any files.
- MongoDB skips Alembic entirely — no `alembic.ini` or `migrations/` folder is generated.
- `SQLModel` and `FastCRUD` skip `app/db/base.py` since SQLModel handles that internally.
- Redis adds `app/core/cache.py` with a ready-to-use cache client.
- If Docker is selected but the daemon isn't running, the tool warns you and exits cleanly.
- If you opt into "Install & start" with Docker, `docker compose up --build` runs automatically after scaffolding.
