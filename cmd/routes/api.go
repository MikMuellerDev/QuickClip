package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/MikMuellerDev/QuickClip/middleware"
	"github.com/MikMuellerDev/QuickClip/sessions"
	utils "github.com/MikMuellerDev/QuickClip/utils"
	"github.com/gorilla/mux"
)

func getUser(r *http.Request) (bool, string) {
	user, ok := middleware.CurrentUser(r)
	return ok, user
}

var clipIdPattern = regexp.MustCompile(`^[A-Za-z0-9._-]{1,64}$`)

// Returns an error message if the username or password is not acceptable
func validateCredentials(username string, password string, checkPassword bool) string {
	if username == "" || len(username) > 64 || strings.ContainsAny(username, ":/") || strings.TrimSpace(username) != username {
		return "The username must be 1-64 characters long, without leading or trailing spaces, ':' or '/'."
	}
	if checkPassword && (len(password) < 8 || len(password) > 128) {
		return "The password must be 8-128 characters long."
	}
	return ""
}

func getVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	version, production := utils.GetVersion()
	json.NewEncoder(w).Encode(VersionStruct{Version: version, Production: production})
}

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

func getClipById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	_, user := getUser(r)
	success, clip := utils.GetClipById(id, user)
	w.Header().Set("Content-Type", "application/json")
	if success {
		json.NewEncoder(w).Encode(clip)
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ResponseStruct{false, 404, "Unknown clip", fmt.Sprintf("The ID: %s is not associated with any Object.", id)})
		return
	}
}

func modClip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var tempClip utils.Clip

	requestError := decoder.Decode(&tempClip)
	if requestError != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(fmt.Sprintf("Invalid Request: %s", requestError.Error()))
		json.NewEncoder(w).Encode(ResponseStruct{false, 400, "Invalid Request", "Your request could not be parsed to a clip"})
		return
	}

	_, user := getUser(r)
	if user == "admin" {
		success, clip := utils.ModClip(tempClip)
		if success {
			json.NewEncoder(w).Encode(clip)
		} else {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ResponseStruct{false, 404, "Unknown clip", fmt.Sprintf("The ID: %s is not associated with any Object.", tempClip.Id)})
		}
	} else {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(ResponseStruct{false, 403, "Permission denied", "You mus be admin to edit clips."})
		return
	}
}

func addClip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var clip utils.Clip
	requestError := decoder.Decode(&clip)

	if requestError != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Warn(fmt.Sprintf("Invalid Request: %s", requestError.Error()))
		json.NewEncoder(w).Encode(ResponseStruct{false, 400, "Invalid Request", "Your request could not be parsed to a clip"})
		return
	}
	_, user := getUser(r)

	if user != "admin" {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(ResponseStruct{false, 403, "Permission denied", "You mus be admin to add clips."})
		return
	}

	if !clipIdPattern.MatchString(clip.Id) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseStruct{false, 400, "Invalid Request", "The Id must be 1-64 characters long and may only contain letters, digits, '.', '_' and '-'."})
		return
	}

	if !utils.AddClip(clip) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseStruct{false, 400, "Invalid Request", fmt.Sprintf("The Id: %s is already matched to a board.", clip.Id)})
		return
	}
	json.NewEncoder(w).Encode(ResponseStruct{Success: true, ErrorCode: 0, Title: "Success", Message: "Clip was added."})
}

// USERS
func getUserList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(utils.GetUsers())
}

func getApiUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	loggedIn, user := getUser(r)
	json.NewEncoder(w).Encode(UserResponse{User: user, LoggedIn: loggedIn})
}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var user utils.User
	requestError := decoder.Decode(&user)

	if requestError != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Warn(fmt.Sprintf("Invalid Request: %s", requestError.Error()))
		json.NewEncoder(w).Encode(ResponseStruct{false, 400, "Invalid Request", "Your request could not be parsed to a User"})
		return
	}

	if message := validateCredentials(user.Name, user.Password, true); message != "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseStruct{false, 400, "Invalid Request", message})
		return
	}

	success := utils.AddUser(user)
	if success {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ResponseStruct{Success: true, ErrorCode: 0, Title: "Success", Message: fmt.Sprintf("User: %s was added.", user.Name)})
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseStruct{false, 400, "Invalid Request", "This user already exists."})
	}
}

func modifyUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["username"]

	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var newUser utils.User
	requestError := decoder.Decode(&newUser)

	if requestError != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Warn(fmt.Sprintf("Invalid Request: %s", requestError.Error()))
		json.NewEncoder(w).Encode(ResponseStruct{false, 400, "Invalid Request", "Your request could not be parsed to a User"})
		return
	}

	changePassword := newUser.Password != "" && newUser.Password != "?"
	if changePassword {
		if message := validateCredentials(id, newUser.Password, true); message != "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ResponseStruct{false, 400, "Invalid Request", message})
			return
		}
	}

	success := utils.AlterUser(id, newUser)
	if success {
		if changePassword {
			sessions.RevokeUser(id)
		}
		json.NewEncoder(w).Encode(ResponseStruct{Success: true, ErrorCode: 0, Title: "Success", Message: fmt.Sprintf("User: %s was altered.", id)})
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ResponseStruct{false, 404, "Unknown User", fmt.Sprintf("The Username: %s does not exist.", id)})
		return
	}
}

func alterPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var newPassword utils.Password
	requestError := decoder.Decode(&newPassword)

	if requestError != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Warn(fmt.Sprintf("Invalid Request: %s", requestError.Error()))
		json.NewEncoder(w).Encode(ResponseStruct{false, 400, "Invalid Request", "Your request could not be parsed to a Password struct"})
		return
	}
	_, user := getUser(r)
	if message := validateCredentials(user, newPassword.Password, true); message != "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseStruct{false, 400, "Invalid Request", message})
		return
	}
	success := utils.AlterPassword(user, newPassword.Password)
	if success {
		// Log out every other session, keep the current one logged in
		_, hadSession := sessions.User(r)
		sessions.RevokeUser(user)
		if hadSession {
			sessions.Login(w, r, user)
		}
		json.NewEncoder(w).Encode(ResponseStruct{Success: true, ErrorCode: 0, Title: "Success", Message: fmt.Sprintf("Password of User: %s was altered.", user)})
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ResponseStruct{false, 404, "Unknown User", fmt.Sprintf("The Username: %s does not exist.", user)})
		return
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["username"]
	w.Header().Set("Content-Type", "application/json")
	if id == "admin" {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(ResponseStruct{false, 403, "Denied", "Cannot delete admin user."})
		return
	}
	success := utils.DeleteUser(id)
	if success {
		sessions.RevokeUser(id)
		json.NewEncoder(w).Encode(ResponseStruct{Success: true, ErrorCode: 0, Title: "Success", Message: fmt.Sprintf("User: %s was deleted.", id)})
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ResponseStruct{false, 404, "Unknown User", fmt.Sprintf("The Username: %s does not exist.", id)})
		return
	}
}

func removeClip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	requestedId := vars["id"]
	_, user := getUser(r)
	if user != "admin" {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(ResponseStruct{false, 403, "Permission denied", "You mus be admin to remove clips."})
		return
	}

	if !utils.DoesClipExist(requestedId) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ResponseStruct{false, 404, "Unknown clip", fmt.Sprintf("The ID: %s is not associated with any Object.", requestedId)})
		return

	}
	utils.RemoveClip(requestedId)
	json.NewEncoder(w).Encode(ResponseStruct{Success: true, ErrorCode: 0, Title: "Success", Message: fmt.Sprintf("Clip with ID: %s was removed.", requestedId)})
}

func probeWriteAccess(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]
	_, user := getUser(r)

	found, clip := utils.GetClipById(id, user)
	if !found {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ResponseStruct{false, 404, "Unknown clip", fmt.Sprintf("[PROBE] The ID: %s is not associated with any Object.", id)})
		return
	}
	if clip.ReadOnly {
		if !utils.HasWritePermission(user, id) {
			json.NewEncoder(w).Encode(ResponseStruct{false, 403, "Read Only", fmt.Sprintf("[PROBE] The ID: %s is set to read-only.", id)})
			return
		}
	}
	json.NewEncoder(w).Encode(ResponseStruct{true, 200, "Write allowed", fmt.Sprintf("[PROBE] Your user is allowed to write to: %s", id)})
}

func editClip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var inputRequest TextInput
	requestError := decoder.Decode(&inputRequest)

	vars := mux.Vars(r)
	id := vars["id"]

	if requestError != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Warn(fmt.Sprintf("Invalid Request: %s", requestError.Error()))
		json.NewEncoder(w).Encode(ResponseStruct{false, 400, "Invalid Request", "Your request could not be parsed to a TextInput"})
		return
	}
	_, user := getUser(r)
	found, clipFromId := utils.GetClipById(id, user)
	if !found {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ResponseStruct{false, 404, "Unknown clip", fmt.Sprintf("The ID: %s is not associated with any Object.", id)})
		return
	}

	if clipFromId.ReadOnly {
		if !utils.HasWritePermission(user, id) {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(ResponseStruct{false, 403, "Read Only", fmt.Sprintf("The ID: %s is set to read-only.", id)})
			return
		}
	}

	editSuccess := utils.EditClip(id, inputRequest.Content, user)

	if editSuccess {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ResponseStruct{Success: true, ErrorCode: 0, Title: "", Message: ""})
	} else {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(ResponseStruct{Success: false, ErrorCode: 403, Title: "Permission Denied", Message: "An error occurred."})
	}
}

type TextInput struct {
	Content string
}
