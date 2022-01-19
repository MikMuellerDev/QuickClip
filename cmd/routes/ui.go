package routes

import (
	"net/http"

	"github.com/MikMuellerDev/QuickClip/middleware"
	"github.com/MikMuellerDev/QuickClip/sessions"
	"github.com/MikMuellerDev/QuickClip/templates"
	"github.com/MikMuellerDev/QuickClip/utils"
	"github.com/gorilla/mux"
)

// UI Pages
func indexGetHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/dash", http.StatusFound)
}

func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Store.Get(r, "session")
	value, ok := session.Values["valid"]
	value2, ok2 := session.Values["username"]
	_, okParse2 := value2.(string)
	valid, okParse := value.(bool)

	log.Printf("username: %s, valid: %t\n", value2, valid)

	if ok && ok2 && okParse2 && okParse && valid {
		http.Redirect(w, r, "/dash", http.StatusFound)
		return
	}

	templates.ExecuteTemplate(w, "login.html", http.StatusOK)
}

func logoutGetHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Store.Get(r, "session")
	session.Values["valid"] = false
	session.Values["username"] = ""
	session.Save(r, w)
	http.Redirect(w, r, "/dash", http.StatusFound)
}

func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")

	if middleware.TestCredentials(username, password, false) {
		session, _ := sessions.Store.Get(r, "session")
		session.Values["valid"] = true
		session.Values["username"] = username
		session.Save(r, w)
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

func editGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, username := getUser(r)
	if utils.DoesClipExist(id) {
		if _, v := utils.GetClipById(id, "admin"); v.Restricted {
			if utils.HasPermission(username, id) {
				templates.ExecuteTemplate(w, "edit.html", http.StatusOK)
				return
			} else {
				templates.ExecuteTemplate(w, "404.html", http.StatusOK)
				return
			}
		} else {
			templates.ExecuteTemplate(w, "edit.html", http.StatusOK)
		}
		return
	} else {
		templates.ExecuteTemplate(w, "404.html", http.StatusOK)
	}
}
