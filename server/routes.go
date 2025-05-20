package server

import (
	"html/template"
	"log"
	"net/http"
	"sync"
)


func addRoutes(mux *http.ServeMux, logger *log.Logger) error {
	th := NewTmplHandler(logger);
	mux.HandleFunc("GET /", th.HomePage)
	mux.HandleFunc("GET /spending", th.SpendingPage)
	return nil
}

type TmplHandler struct {
	logger *log.Logger
	// lazy load each page's template
	templatesMu sync.Mutex
	templates   map[string]*template.Template
}


func NewTmplHandler(logger *log.Logger) *TmplHandler {
	return &TmplHandler{
		logger: logger, 
		templates: make(map[string]*template.Template), 
	}
}

// lazily load a template
func (h *TmplHandler) getTemplate(name string) (*template.Template, error) {
	h.templatesMu.Lock()
	defer h.templatesMu.Unlock()

	if tmpl, ok := h.templates[name]; ok {
		return tmpl, nil
	}

	h.logger.Printf("Parsing template '%s'", name)

	tmpl, err := template.ParseFiles("templates/layout.html")
	if err != nil {
		return nil, err
	}

	// Parse the specific page template
	tmpl, err = tmpl.ParseFiles("templates/pages/" + name + ".html")
	if err != nil {
		return nil, err
	}

	h.templates[name] = tmpl
	return tmpl, nil
}

func (h *TmplHandler) HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.renderError(w, "Not Found", http.StatusNotFound)
		return
	}
	
	tmpl, err := h.getTemplate("home")
	if err != nil {
		h.logger.Printf("Error loading home template: %v", err)
		h.renderError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"Title": "Home",
	}
	if err := tmpl.Execute(w, data); err != nil {
		h.logger.Printf("Error rendering home template: %v", err)
		h.renderError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *TmplHandler) SpendingPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := h.getTemplate("spending")
	if err != nil {
		h.logger.Printf("Error loading spending template: %v", err)
		h.renderError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"Title": "Spending",
	}
	if err := tmpl.Execute(w, data); err != nil {
		h.logger.Printf("Error rendering spending template: %v", err)
		h.renderError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *TmplHandler) NotFound(w http.ResponseWriter, r *http.Request) {
	h.renderError(w, "Not Found", http.StatusNotFound)
}


func (h *TmplHandler) renderError(w http.ResponseWriter, errMsg string, statusCode int) {
	w.WriteHeader(statusCode)

	tmpl, err := h.getTemplate("error")
	if err != nil {
		h.logger.Printf("Error loading error template: %v", err)
		http.Error(w, errMsg, statusCode)
		return
	}

	data := map[string]any{
		"Title":        "Error",
		"ErrorMessage": errMsg,
		"StatusCode":   statusCode,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		h.logger.Printf("Error rendering error template: %v", err)
		http.Error(w, errMsg, statusCode)
		return
	}
}