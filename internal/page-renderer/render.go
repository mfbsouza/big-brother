package render

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type page struct {
	template *template.Template
	status   os.FileInfo
}

// page cache
var pc = make(map[string]*page)

func RenderTemplate(w http.ResponseWriter, f string) {
	// check if the required page is on the cache
	_, exists := pc[f]
	if !exists {
		log.Println("[render] loading template:", f)
		err := loadPage(f)
		if err != nil {
			log.Println("[render] failed loading template:", err)
			return
		}
	} else {
		// check if there was an update to the disk file
		cs, err := os.Stat("./templates/" + f)
		if err != nil {
			log.Println("[render] failed loading file status:", err)
			return
		}
		if cs.ModTime() != pc[f].status.ModTime() {
			// reload template
			log.Println("[render] reloading template:", f)
			err := loadPage(f)
			if err != nil {
				log.Println("[render] failed loading template:", err)
				return
			}
		} else {
			log.Printf("[render] using template %s from cache", f)
		}
	}

	p := pc[f]
	err := p.template.Execute(w, nil)
	if err != nil {
		log.Println("[render] error: failed to execute the template:", err)
		return
	}
}

func loadPage(f string) error {
	templates := []string{
		fmt.Sprintf("./templates/%s", f),
		"./templates/navbar.html",
		"./templates/base.html",
	}

	// load templates from disk
	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		return err
	}

	// load file status
	stat, err := os.Stat(templates[0])
	if err != nil {
		return err
	}

	// create the page struct
	p := page{
		template: tmpl,
		status:   stat,
	}

	// add page to the cache
	pc[f] = &p

	return nil
}
