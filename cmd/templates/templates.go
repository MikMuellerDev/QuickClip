package templates

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func InitLogger(logger *logrus.Logger) {
	log = logger
}

var templates *template.Template

func LoadTemplates(pattern string) {
	templates = template.Must(template.ParseGlob(pattern))
	log.Debug(fmt.Sprintf("Templates loaded: %s", pattern))
}

func ExecuteTemplate(responseWriter http.ResponseWriter, templateName string, statusCode int) {
	responseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
	responseWriter.WriteHeader(statusCode)
	if err := templates.ExecuteTemplate(responseWriter, templateName, nil); err != nil {
		log.Error(fmt.Sprintf("Could not execute template %s: %s", templateName, err.Error()))
	}
}
