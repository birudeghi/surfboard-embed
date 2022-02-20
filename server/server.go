package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/birudeghi/surfboard-embed/session"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Dimension struct {
	Width  string
	Height string
}

type ReqSchema struct {
	EventType  string
	WebsiteUrl string
	SessionId  string
	CopyPaste  CopyPaste
	TimeTaken  float64
	ResizeFrom Dimension
	ResizeTo   Dimension
}

type Board struct {
	Fqbn         string `bson:"fqbn" json:"fqbn"`
	Name         string `bson:"name" json:"name"`
	Version      string `bson:"version" json:"version"`
	PropertiesId string `bson:"properties_id" json:"propertiesId"`
}

type Lib struct {
	Name    string `bson:"name" json:"name"`
	Author  string `bson:"author" json:"author"`
	Url     string `bson:"url" json:"url"`
	Version string `bson:"version" json:"version"`
}

type File struct {
	Content string `bson:"content" json:"content"`
	Name    string `bson:"name" json:"name"`
}

type MetricResSchema struct {
	Created bool
}

type ErrorSchema struct {
	Errors []Error `json:"errors"`
}

type SurfSchema struct {
	AppId            string                   `bson:"app_id,omitempty" json:"appId,omitempty"`
	Email            string                   `bson:"email" json:"email"`
	CompatibleBoards []map[string]interface{} `bson:"compatible_boards" json:"compatibleBoards"`
	Libs             []Lib                    `bson:"libs" json:"libs"`
	Files            []File                   `bson:"files" json:"files"`
}

type Error struct {
	ErrorType string `json:"errorType"`
	Message   string `json:"message"`
}

type CopyPaste struct {
	FormId string
	Pasted bool
}

var globalSessions *session.Manager
var database = "surfboard-user"
var collection = "app-data"

//go:embed build
var content embed.FS

var err = os.Setenv("PASSWORD", "NnGvyg2BNKTUdybR")

var uri = fmt.Sprintf("mongodb+srv://birudeghi:%s@surfboard-basic.1v9tc.mongodb.net/surfboard-basic", os.Getenv("PASSWORD"))

func init() {
	globalSessions, _ = session.NewManager("memory", "sessionId", 3600)

	go globalSessions.GC()
}

func generateError(errorType string, message string) ErrorSchema {
	res := ErrorSchema{}
	res.Errors = append(res.Errors, Error{errorType, message})
	return res
}

func clientHandler() http.Handler {
	fsys := fs.FS(content)

	contentStatic, _ := fs.Sub(fsys, "build")

	return http.FileServer(http.FS(contentStatic))
}

// Generate a JSON response to be encoded into the HTTP response
func jsonResponse(w http.ResponseWriter, res interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options:", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(res)
}

func main() {
	r := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"X-PINGOTHER", "Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	r.Use(cors.Handler)

	r.Mount("/", clientHandler())
	r.Post("/metric/submit", submitHandler)
	r.Post("/metric/screen-resize", scrResizeHandler)
	r.Post("/metric/copy-paste", cPasteHandler)
	r.Post("/onboard/create", createHandler)
	r.Get("/app/info/{appId}", appInfoHandler)

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
