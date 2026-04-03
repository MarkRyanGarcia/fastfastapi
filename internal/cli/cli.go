package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/markryangarcia/fastfastapi/generator"
	"github.com/markryangarcia/fastfastapi/tui"
)

var (
	green  = lipgloss.NewStyle().Foreground(lipgloss.Color("78"))
	cyan   = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	muted  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	border = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	check  = lipgloss.NewStyle().Foreground(lipgloss.Color("78")).Bold(true)
	errSty = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
)

func pipe() string { return border.Render("│") }

func Run() {
	fs := flag.NewFlagSet("fastfastapi", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: fastfastapi [name|.] [flags]

Flags:
  --db        Database: postgres, mongo (default: interactive)
  --orm       ORM: sqlalchemy, sqlmodel, fastcrud (postgres only)
  --auth      Auth provider: none, clerk, cognito (default: none)
  --pkg       Package manager: pip, pipenv (default: pip)
  --docker    Enable Docker support
  --redis     Enable Redis caching
  --install   Install dependencies and start after scaffolding
  --no-tui    Skip TUI even if flags are incomplete (use defaults)

Examples:
  fastfastapi my-api --db postgres --orm sqlalchemy --auth none --pkg pip
  ffa my-api --db postgres --orm sqlalchemy --auth none --pkg pip
  fastfastapi . --db mongo --docker --redis
  fastfastapi my-api --db postgres --orm fastcrud --docker --install
`)
	}

	dbFlag := fs.String("db", "", "Database: postgres, mongo")
	ormFlag := fs.String("orm", "", "ORM: sqlalchemy, sqlmodel, fastcrud")
	authFlag := fs.String("auth", "", "Auth: none, clerk, cognito")
	pkgFlag := fs.String("pkg", "", "Package manager: pip, pipenv")
	dockerFlag := fs.Bool("docker", false, "Enable Docker support")
	redisFlag := fs.Bool("redis", false, "Enable Redis caching")
	installFlag := fs.Bool("install", false, "Install and start after scaffolding")
	noTUIFlag := fs.Bool("no-tui", false, "Skip TUI, use defaults for missing flags")

	args := os.Args[1:]
	var initialName string
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		arg := args[0]
		if arg == "." {
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Println(errSty.Render("❌ Could not get current directory: " + err.Error()))
				os.Exit(1)
			}
			initialName = filepath.Base(cwd)
		} else {
			initialName = arg
		}
		args = args[1:]
	}
	fs.Parse(args)

	if *dbFlag != "" || *noTUIFlag {
		runNonInteractive(initialName, *dbFlag, *ormFlag, *authFlag, *pkgFlag, *dockerFlag, *redisFlag, *installFlag)
		return
	}

	p := tea.NewProgram(tui.InitialModelWithName(initialName))
	finalModel, err := p.Run()
	if err != nil {
		fmt.Println(errSty.Render("Error: " + err.Error()))
		os.Exit(1)
	}

	m := finalModel.(tui.Model)

	if m.Selected == "" || m.Quitting {
		fmt.Println(muted.Render("\n  Generation cancelled."))
		return
	}

	fmt.Println(m.Summary())

	buildAndRun(
		m.ProjectName,
		initialName,
		m.Selected,
		m.ORMChoice,
		m.AuthProvider,
		m.UsePipenv,
		m.UseDocker,
		m.UseRedis,
		m.SetupVenv,
	)
}

func runNonInteractive(name, db, orm, auth, pkg string, docker, redis, install bool) {
	if name == "" {
		name = "my-fastapi-app"
	}

	var selected string
	switch strings.ToLower(db) {
	case "mongo", "mongodb":
		selected = "MongoDB (PyMongo)"
	default:
		selected = "PostgreSQL (SQLAlchemy)"
	}

	var ormChoice string
	if !strings.Contains(selected, "MongoDB") {
		switch strings.ToLower(orm) {
		case "sqlmodel":
			ormChoice = "SQLModel"
		case "fastcrud":
			ormChoice = "FastCRUD"
		default:
			ormChoice = "SQLAlchemy"
		}
	}

	var authProvider string
	switch strings.ToLower(auth) {
	case "clerk":
		authProvider = "Clerk"
	case "cognito", "awscognito", "aws_cognito":
		authProvider = "AWS Cognito"
	default:
		authProvider = "None"
	}

	usePipenv := strings.ToLower(pkg) == "pipenv"

	fmt.Println(pipe())
	fmt.Printf("%s  %s %s\n", pipe(), cyan.Render("Project:       "), name)
	fmt.Printf("%s  %s %s\n", pipe(), cyan.Render("Database:      "), selected)
	if ormChoice != "" {
		fmt.Printf("%s  %s %s\n", pipe(), cyan.Render("ORM:           "), ormChoice)
	}
	fmt.Printf("%s  %s %s\n", pipe(), cyan.Render("Auth:          "), authProvider)
	pkgLabel := "pip"
	if usePipenv {
		pkgLabel = "pipenv"
	}
	fmt.Printf("%s  %s %s\n", pipe(), cyan.Render("Pkg manager:   "), pkgLabel)
	fmt.Printf("%s  %s %v\n", pipe(), cyan.Render("Docker:        "), docker)
	fmt.Printf("%s  %s %v\n", pipe(), cyan.Render("Redis:         "), redis)
	fmt.Println(pipe())
	fmt.Println(check.Render("◇  ") + green.Render("Scaffolding project in ./"+name+"..."))

	buildAndRun(name, name, selected, ormChoice, authProvider, usePipenv, docker, redis, install)
}

func buildAndRun(projectName, outArg, selected, ormChoice, authProvider string, usePipenv, useDocker, useRedis, setupVenv bool) {
	isSQL := strings.Contains(selected, "SQL")
	isMongo := strings.Contains(selected, "MongoDB")
	isSQLModel := ormChoice == "SQLModel"
	isFastCRUD := ormChoice == "FastCRUD"
	isPlainSQL := isSQL && !isSQLModel && !isFastCRUD

	outDir := projectName
	if outArg == "." {
		outDir = "."
	}

	config := generator.ProjectConfig{
		ProjectName:       projectName,
		OutputDir:         outDir,
		Database:          selected,
		IncludeSQLAlchemy: isPlainSQL,
		IncludeMongoDB:    isMongo,
		UseSQLModel:       isSQLModel,
		UseFastCRUD:       isFastCRUD,
		AuthProvider:      authProvider,
		UseClerk:          authProvider == "Clerk",
		UseCognito:        authProvider == "AWS Cognito",
		UsePipenv:         usePipenv,
		SetupVenv:         setupVenv,
		UseDocker:         useDocker,
		UseRedis:          useRedis,
	}

	if useDocker && !generator.IsDockerRunning() {
		fmt.Println(errSty.Render("❌ Docker doesn't appear to be running."))
		fmt.Println(pipe() + "  " + muted.Render("Start Docker, then re-run with the same flags."))
		fmt.Println(pipe())
		os.Exit(1)
	}

	if err := generator.CreateProject(config); err != nil {
		fmt.Println(errSty.Render("❌ Failed to create project: " + err.Error()))
		os.Exit(1)
	}

	fmt.Println(check.Render("◇  ") + green.Render("Done! Project generated in ./"+outDir))
	fmt.Println(pipe())

	if setupVenv {
		if useDocker {
			fmt.Println(pipe() + "  " + cyan.Render("Running docker compose up --build..."))
			fmt.Println(pipe())
			if err := generator.RunDockerCompose(outDir); err != nil {
				fmt.Println(errSty.Render("❌ Failed to start Docker: " + err.Error()))
				os.Exit(1)
			}
		} else {
			fmt.Println(pipe() + "  " + cyan.Render("Starting dev server..."))
			fmt.Println(pipe())
			if err := generator.RunDevServer(outDir, usePipenv); err != nil {
				fmt.Println(errSty.Render("❌ Failed to start dev server: " + err.Error()))
				os.Exit(1)
			}
		}
	} else {
		fmt.Println(pipe() + "  " + cyan.Render("Next steps:"))
		if outDir != "." {
			fmt.Println(pipe() + "  " + muted.Render("cd "+outDir))
		}
		if useDocker {
			fmt.Println(pipe() + "  " + muted.Render("docker compose up --build"))
		} else if usePipenv {
			fmt.Println(pipe() + "  " + muted.Render("pipenv install"))
			fmt.Println(pipe() + "  " + muted.Render("pipenv shell"))
			fmt.Println(pipe() + "  " + muted.Render("fastapi dev app"))
		} else {
			fmt.Println(pipe() + "  " + muted.Render("pip install -r requirements.txt"))
			fmt.Println(pipe() + "  " + muted.Render("source .venv/bin/activate"))
			fmt.Println(pipe() + "  " + muted.Render("fastapi dev app"))
		}
		fmt.Println(pipe())
	}
}
