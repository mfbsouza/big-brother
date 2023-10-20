package render

import (
	"fmt"
	"log"
	"net/http"
	"html/template"
)

// templates cache
var tc = make(map[string]*template.Template)

func RenderTemplate(w http.ResponseWriter, f string) {
	// check if the required template is on the cache
	_, exists := tc[f]
	if !exists {
		log.Println("[render] loading template:", f)
		err := loadTemplate(f)
		if err != nil {
			log.Println("[render] failed loading template:", err)
			return
		}
	} else {
		log.Printf("[render] using template %s from cache", f)
	}

	t := tc[f]
	err := t.Execute(w, nil)
	if err != nil {
		log.Println("[render] error: failed to execute the template:", err)
		return
	}
}

func loadTemplate(f string) error {
	templates := []string{
		fmt.Sprintf("./templates/%s", f),
		"./templates/base.html",
	}

	// load templates from disk
	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		return err
	}

	// add template to the cache
	tc[f] = tmpl

	return nil
}
