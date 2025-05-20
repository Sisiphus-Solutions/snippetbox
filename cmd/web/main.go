package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"snippetbox.sisiphus.dev/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

type config struct {
	addr      string
	staticDir string
}

type application struct {
	logger   *slog.Logger
	snippets *models.SnippetModel
}

func main() {
	// Configuration
	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	// Setup logging
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Setup Database Connection
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// We also defer a call to db.Close(), so that the connection pool is closed
	// before the main() function exits.
	defer db.Close()

	// Initialize the struct
	app := &application{
		logger:   logger,
		snippets: &models.SnippetModel{DB: db},
	}

	// Start the server and log any errors
	logger.Info("starting server on %s", "addr", cfg.addr)
	err = http.ListenAndServe(cfg.addr, app.routes(cfg))
	logger.Error(err.Error())
	os.Exit(1)
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
