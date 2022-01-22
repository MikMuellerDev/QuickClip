package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MikMuellerDev/QuickClip/middleware"
	"github.com/MikMuellerDev/QuickClip/sessions"
	utils "github.com/MikMuellerDev/QuickClip/utils"
	"github.com/gorilla/mux"
)

func getUser(r *http.Request) (bool, string) {
	session, _ := sessions.Store.Get(r, "session")

	query := r.URL.Query()
	password := query.Get("password")

	uname2 := query.Get("username")
	uname1, _ := session.Values["username"]
	uname1Parsed, okParse := uname1.(string)

	if okParse && utils.DoesUserExist(uname1Parsed) {
		return true, uname1Parsed
	} else if middleware.TestCredentials(uname2, password, true) {
		return true, uname2
	} else {
		return false, ""
	}
}

// Returns the mode which is currently active
func getVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	version, production := utils.GetVersion()
	json.NewEncoder(w).Encode(VersionStruct{Version: version, Production: production})
}

// Returns the clip list
func getClips(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, username := getUser(r)
	json.NewEncoder(w).Encode(utils.GetClips(username))
}

func saveToFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	action := utils.RequestSave()

	if action {
		json.NewEncoder(w).Encode(ResponseStruct{true, 0, "Changes Saved.", "Changes you made will now be saved to disk."})
	} else {
		json.NewEncoder(w).Encode(ResponseStruct{true, 0, "Not saved.", "Your changes are already saved."})
	}

}

// Returns Clip Based on Id
func getClipById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	_, user := getUser(r)
	success, clip := utils.GetClipById(id, user)
	w.Header().Set("Content-Type", "application/json")
	if success {
		json.NewEncoder(w).Encode(clip)
	} else {
		json.NewEncoder(w).Encode(ResponseStruct{false, 404, "Unknown clip", fmt.Sprintf("The ID: %s is not associated with any Object.", id)})
		return
	}
}

func modClip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var t utils.Clip
	err := decoder.Decode(&t)
	if err != nil {
		log.Error(fmt.Sprintf("Invalid Request: %s", err.Error()))
		json.NewEncoder(w).Encode(ResponseStruct{false, 400, "Invalid Request", "Your request could not be parsed to a clip"})
		return
	}

	_, user := getUser(r)
	if user == "admin" {
		success, clip := utils.ModClip(t)
		if success {
			json.NewEncoder(w).Encode(clip)
		} else {
			json.NewEncoder(w).Encode(ResponseStruct{false, 404, "Unknown clip", fmt.Sprintf("The ID: %s is not associated with any Object.", t.Id)})
		}
	} else {
		json.NewEncoder(w).Encode(ResponseStruct{false, 401, "Permission denied", "You mus be admin to edit clips."})
		return
	}
}

func addClip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var clip utils.Clip
	err := decoder.Decode(&clip)

	if err != nil {
		log.Warn(fmt.Sprintf("Invalid Request: %s", err.Error()))
		json.NewEncoder(w).Encode(ResponseStruct{false, 400, "Invalid Request", "Your request could not be parsed to a clip"})
		return
	}
	_, user := getUser(r)

	if user != "admin" {
		json.NewEncoder(w).Encode(ResponseStruct{false, 401, "Permission denied", "You mus be admin to add clips."})
	}

	if utils.DoesClipExist(clip.Id) {
		json.NewEncoder(w).Encode(ResponseStruct{false, 400, "Invalid Request", fmt.Sprintf("The Id: %s is already matched to a board.", clip.Id)})
		return
	}

	utils.AddClip(clip)
	json.NewEncoder(w).Encode(ResponseStruct{Success: true, ErrorCode: 0, Title: "Success", Message: "Clip was added."})
}

func getApiUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	loggedIn, user := getUser(r)
	json.NewEncoder(w).Encode(UserResponse{User: user, LoggedIn: loggedIn})
}

func removeClip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	_, user := getUser(r)
	if user != "admin" {
		json.NewEncoder(w).Encode(ResponseStruct{false, 401, "Permission denied", "You mus be admin to add clips."})
	}

	if !utils.DoesClipExist(id) {
		json.NewEncoder(w).Encode(ResponseStruct{false, 404, "Unknown clip", fmt.Sprintf("The ID: %s is not associated with any Object.", id)})
		return

	}
	utils.RemoveClip(id)
	json.NewEncoder(w).Encode(ResponseStruct{Success: true, ErrorCode: 0, Title: "Success", Message: fmt.Sprintf("Clip with ID: %s was removed.", id)})
}

func probeWriteAccess(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]
	_, user := getUser(r)

	_, clip := utils.GetClipById(id, user)
	if clip.ReadOnly {
		if !utils.HasWritePermission(user, id) {
			json.NewEncoder(w).Encode(ResponseStruct{false, 401, "Read Only", fmt.Sprintf("[PROBE] The ID: %s is set to read-only.", id)})
			return
		}
	}
	json.NewEncoder(w).Encode(ResponseStruct{true, 200, "Write allowed", fmt.Sprintf("[PROBE] Your user is allowed to write to: %s", id)})
}

func editClip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var req TextInput
	err := decoder.Decode(&req)

	vars := mux.Vars(r)
	id := vars["id"]

	if err != nil {
		log.Warn(fmt.Sprintf("Invalid Request: %s", err.Error()))
		json.NewEncoder(w).Encode(ResponseStruct{false, 400, "Invalid Request", "Your request could not be parsed to a TextInput"})
		return
	}

	if !utils.DoesClipExist(id) {
		json.NewEncoder(w).Encode(ResponseStruct{false, 404, "Unknown clip", fmt.Sprintf("The ID: %s is not associated with any Object.", id)})
		return
	}

	_, user := getUser(r)

	_, clip := utils.GetClipById(id, user)
	if clip.ReadOnly {
		if !utils.HasWritePermission(user, id) {
			json.NewEncoder(w).Encode(ResponseStruct{false, 401, "Read Only", fmt.Sprintf("The ID: %s is set to read-only.", id)})
			return
		}
	}

	success := utils.EditClip(id, req.Content, user)

	if success {
		json.NewEncoder(w).Encode(ResponseStruct{Success: true, ErrorCode: 0, Title: "", Message: ""})
	} else {
		json.NewEncoder(w).Encode(ResponseStruct{Success: false, ErrorCode: 401, Title: "Permission Denied", Message: "An error occurred."})
	}
}

type TextInput struct {
	Content string
}
