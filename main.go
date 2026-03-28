package main

import (
	"context"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	if err := RunMigrations(context.Background(), pool); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	if err := SeedContacts(context.Background(), pool); err != nil {
		log.Fatalf("Failed to seed contacts: %v", err)
	}

	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	funcMap := template.FuncMap{
		"formatTime": formatTime,
	}

	templates, err := template.New("").Funcs(funcMap).ParseGlob(filepath.Join("templates", "*.html"))
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}
	e.Renderer = &TemplateRenderer{templates: templates}

	e.Static("/static", "static")

	h := &Handler{db: pool}

	// Page routes
	e.GET("/", h.ListContacts)
	e.GET("/contacts/new", h.NewContactForm)
	e.POST("/contacts", h.CreateContact)
	e.GET("/contacts/:id", h.ViewContact)
	e.GET("/contacts/:id/edit", h.EditContactForm)
	e.POST("/contacts/:id", h.UpdateContact)
	e.POST("/contacts/:id/delete", h.DeleteContact)

	// API routes
	e.GET("/api/contacts", h.APIListContacts)
	e.POST("/api/contacts", h.APICreateContact)
	e.GET("/api/contacts/:id", h.APIGetContact)
	e.PUT("/api/contacts/:id", h.APIUpdateContact)
	e.DELETE("/api/contacts/:id", h.APIDeleteContact)

	// Health check
	e.GET("/healthz", HealthCheck(pool))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func formatTime(t interface{}) string {
	if ts, ok := t.(string); ok {
		return ts
	}
	return ""
}
