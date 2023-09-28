package render

import (
	"log"
	"net/http"
	"html/template"
)

func RenderTemplate(w http.ResponseWriter, n string) {
	p := "./templates/" + n
	log.Println("[render] loading template:", p)
	t, _ := template.ParseFiles(p, "./templates/base.html")
	err := t.Execute(w, nil)
	if err != nil {
		log.Println("[render] error: failed to parse the template:", err)
		return
	}
}
