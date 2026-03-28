# fastapi-gen

A CLI tool that scaffolds a production-ready FastAPI project with an interactive TUI. Choose your database backend (PostgreSQL via SQLAlchemy or MongoDB via PyMongo) and package manager (pipenv or requirements.txt) and get a fully structured project in seconds.

## Prerequisites

- [Go](https://go.dev/dl/) 1.21+
- Python 3.11+
- `pipenv` (optional, only if you choose it during setup)

## Installation

### Install with `go install`

```bash
go install github.com/markryangarcia/fastapi-gen@latest
```

Make sure `$GOPATH/bin` (or `$HOME/go/bin`) is in your `PATH`.

## Usage

### Create a new project in a new directory

```bash
fastapi-gen
```

The TUI will walk you through three steps:

1. Project name
2. Database — `PostgreSQL (SQLAlchemy)` or `MongoDB (PyMongo)`
3. Package manager — `pipenv` or `requirements.txt`
4. Virtual environment — set it up automatically or skip

A new folder named after your project will be created in the current directory.

### Scaffold into the current directory

```bash
fastapi-gen .
```

Skips the name prompt and uses the current directory name as the project name. Files are written in-place.

### Pass a project name directly

```bash
fastapi-gen my-api
```

Skips the name prompt and goes straight to database selection.

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
│   │   └── security.py
│   ├── db/
│   │   ├── base.py
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
│   └── main.py
├── migrations/          # PostgreSQL only (Alembic)
│   └── versions/
├── tests/
│   ├── test_users.py
│   └── test_items.py
├── conftest.py
├── .env
├── .gitignore
├── alembic.ini          # PostgreSQL only
├── requirements.txt     # if not using pipenv
└── Pipfile              # if using pipenv
```

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

### 2. Install dependencies

**With pipenv + auto venv setup** (runs `pipenv install --dev` automatically):
```bash
cd my-api && pipenv shell
```

**With pipenv, manual setup:**
```bash
cd my-api
pipenv install --dev
pipenv shell
```

**With requirements.txt + auto venv setup** (`.venv` is created automatically):
```bash
source my-api/.venv/bin/activate  # Windows: my-api\.venv\Scripts\activate
pip install -r requirements.txt
```

**With requirements.txt, manual setup:**
```bash
python3 -m venv .venv
source .venv/bin/activate  # Windows: .venv\Scripts\activate
pip install -r requirements.txt
```

### 3. Start the development server

```bash
fastapi dev app
```

The API will be available at `http://localhost:8000`.  
Interactive docs at `http://localhost:8000/docs`.

### 5. Run tests

```bash
pytest
```

## Notes

- If you opt into automatic venv setup with pipenv, `pipenv install --dev` runs right after scaffolding.
- If you opt into automatic venv setup with requirements.txt, a `.venv` folder is created via `python3 -m venv .venv` — you still need to activate it and run `pip install`.
- The MongoDB option skips Alembic entirely — no `alembic.ini` or `migrations/` folder is generated.
- `q` or `Ctrl+C` at any point in the TUI cancels generation without writing any files.
