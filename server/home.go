package server

import (
	"encoding/json"
	"log"
	"net/http"
	"text/template"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})

	if r.Method != http.MethodGet {
		response["error"] = "Invalid request method."
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		tmpl := `<html>
                    <head><title>Page Not Found</title></head>
                    <body>
                        <h1>404 - Page Not Found</h1>
                    </body>
                 </html>`
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(tmpl))
		return
	}

	tmpl, err := template.ParseFiles("./static/index.html")
	if err != nil {
		log.Printf("Template parsing error: %v", err)
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Error rendering posts", http.StatusInternalServerError)
		return
	}
}


