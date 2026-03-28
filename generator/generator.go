package generator

import (
	"embed"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

//go:embed all:templates
var templateFS embed.FS

type ProjectConfig struct {
	ProjectName       string
	OutputDir         string
	Database          string
	IncludeSQLAlchemy bool
	IncludeMongoDB    bool
	UsePipenv         bool
	SetupVenv         bool
}

// fileMap maps destination path -> template path (executed as Go templates)
func fileMap() map[string]string {
	return map[string]string{
		".gitignore":                        "templates/.gitignore.tmpl",
		"conftest.py":                       "templates/conftest.py.tmpl",
		"alembic.ini":                       "templates/alembic.ini.tmpl",
		"migrations/env.py":                 "templates/migrations/env.py.tmpl",
		"migrations/versions/.gitkeep":      "templates/migrations/versions/.gitkeep.tmpl",
		"app/main.py":                       "templates/app/main.py.tmpl",
		"app/__init__.py":                   "templates/app/__init__.py.tmpl",
		"app/api/v1/__init__.py":            "templates/app/api/v1/__init__.py.tmpl",
		"app/api/v1/routers/users.py":       "templates/app/api/v1/routers/users.py.tmpl",
		"app/api/v1/routers/items.py":       "templates/app/api/v1/routers/items.py.tmpl",
		"app/core/config.py":                "templates/app/core/config.py.tmpl",
		"app/core/security.py":              "templates/app/core/security.py.tmpl",
		"app/db/session.py":                 "templates/app/db/session.py.tmpl",
		"app/db/base.py":                    "templates/app/db/base.py.tmpl",
		"app/models/user.py":                "templates/app/models/user.py.tmpl",
		"app/models/item.py":                "templates/app/models/item.py.tmpl",
		"app/schemas/user.py":               "templates/app/schemas/user.py.tmpl",
		"app/schemas/item.py":               "templates/app/schemas/item.py.tmpl",
		"app/services/user_service.py":      "templates/app/services/user_service.py.tmpl",
		"app/services/item_service.py":      "templates/app/services/item_service.py.tmpl",
		"tests/test_users.py":               "templates/tests/test_users.py.tmpl",
		"tests/test_items.py":               "templates/tests/test_items.py.tmpl",
		"requirements.txt":                  "templates/requirements.txt.tmpl",
		"Pipfile":                            "templates/Pipfile.tmpl",
		".env":                              "templates/.env.tmpl",
		"README.md":                         "templates/README.md.tmpl",
	}
}

// rawFileMap maps destination path -> embed path for files copied verbatim (no template execution)
func rawFileMap() map[string]string {
	return map[string]string{
		"migrations/script.py.mako": "templates/migrations/script.py.mako.tmpl",
	}
}

var alembicFiles = map[string]bool{
	"alembic.ini":              true,
	"migrations/env.py":        true,
	"migrations/versions/.gitkeep": true,
}

var alembicRawFiles = map[string]bool{
	"migrations/script.py.mako": true,
}

func CreateProject(cfg ProjectConfig) error {
	if cfg.OutputDir == "" {
		cfg.OutputDir = cfg.ProjectName
	}
	for dest, tmplPath := range fileMap() {
		if cfg.IncludeMongoDB && alembicFiles[dest] {
			continue
		}
		// Skip requirements.txt if using pipenv, skip Pipfile if not
		if dest == "requirements.txt" && cfg.UsePipenv {
			continue
		}
		if dest == "Pipfile" && !cfg.UsePipenv {
			continue
		}
		if err := writeTemplate(cfg, dest, tmplPath); err != nil {
			return err
		}
	}
	for dest, src := range rawFileMap() {
		if cfg.IncludeMongoDB && alembicRawFiles[dest] {
			continue
		}
		if err := copyRaw(cfg.OutputDir, dest, src); err != nil {
			return err
		}
	}
	if cfg.SetupVenv {
		if cfg.UsePipenv {
			if err := runPipenvInstall(cfg.OutputDir); err != nil {
				return err
			}
		} else {
			if err := runVenvCreate(cfg.OutputDir); err != nil {
				return err
			}
		}
	}
	return nil
}

func runPipenvInstall(dir string) error {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	cmd := exec.Command("pipenv", "install", "--dev")
	cmd.Dir = absDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runVenvCreate(dir string) error {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	// Create the virtual environment
	create := exec.Command("python3", "-m", "venv", ".venv")
	create.Dir = absDir
	create.Stdout = os.Stdout
	create.Stderr = os.Stderr
	if err := create.Run(); err != nil {
		return err
	}
	// Install packages into the venv
	pip := exec.Command(".venv/bin/pip", "install", "-r", "requirements.txt")
	pip.Dir = absDir
	pip.Stdout = os.Stdout
	pip.Stderr = os.Stderr
	return pip.Run()
}

func writeTemplate(cfg ProjectConfig, dest, tmplPath string) error {
	outPath := filepath.Join(cfg.OutputDir, filepath.FromSlash(dest))
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return err
	}
	tmpl, err := template.ParseFS(templateFS, tmplPath)
	if err != nil {
		return err
	}
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return tmpl.Execute(f, cfg)
}

func copyRaw(outputDir, dest, src string) error {
	outPath := filepath.Join(outputDir, filepath.FromSlash(dest))
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return err
	}
	in, err := templateFS.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
