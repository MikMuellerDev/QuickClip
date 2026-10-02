package routes

import (
	"fmt"
	"net/http"

	"github.com/MikMuellerDev/QuickClip/middleware"
	"github.com/MikMuellerDev/QuickClip/sessions"
	"github.com/MikMuellerDev/QuickClip/templates"
	"github.com/MikMuellerDev/QuickClip/utils"
	"github.com/gorilla/mux"
)

func indexGetHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/dash", http.StatusFound)
}

func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	if _, ok := sessions.User(r); ok {
		http.Redirect(w, r, "/dash", http.StatusFound)
		return
	}
	templates.ExecuteTemplate(w, "login.html", http.StatusOK)
}

func logoutGetHandler(w http.ResponseWriter, r *http.Request) {
	sessions.Logout(w, r)
	http.Redirect(w, r, "/dash", http.StatusFound)
}

func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	success, limited := middleware.TestCredentials(r, username, password)
	if limited {
		templates.ExecuteTemplate(w, "login.html", http.StatusTooManyRequests)
		return
	}
	if success {
		if err := sessions.Login(w, r, username); err != nil {
			log.Error("Could not create session: ", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if username == "admin" {
			http.Redirect(w, r, "/admin", http.StatusFound)
			return
		}
		http.Redirect(w, r, "/dash", http.StatusFound)
	} else {
		templates.ExecuteTemplate(w, "login.html", http.StatusForbidden)
	}
}

func dashGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "dash.html", http.StatusOK)
}

func adminGetHandler(w http.ResponseWriter, r *http.Request) {
	_, user := getUser(r)
	if user == "admin" {
		templates.ExecuteTemplate(w, "admin.html", http.StatusOK)
	} else {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
}

func adminUserSettingsGetHandler(w http.ResponseWriter, r *http.Request) {
	_, user := getUser(r)
	if user == "admin" {
		templates.ExecuteTemplate(w, "user.html", http.StatusOK)
	} else {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
}

func editGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestedId := vars["id"]
	_, username := getUser(r)
	if utils.DoesClipExist(requestedId) {
		if _, v := utils.GetClipById(requestedId, "admin"); v.Restricted {
			if utils.HasPermission(username, requestedId) {
				templates.ExecuteTemplate(w, "edit.html", http.StatusOK)
				return
			} else {
				templates.ExecuteTemplate(w, "404.html", http.StatusNotFound)
				return
			}
		} else {
			templates.ExecuteTemplate(w, "edit.html", http.StatusOK)
		}
		return
	} else {
		templates.ExecuteTemplate(w, "404.html", http.StatusNotFound)
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "404.html", http.StatusNotFound)
}

func robotsTxtHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "User-agent: * Disallow:\n")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Status: Ok")
}
