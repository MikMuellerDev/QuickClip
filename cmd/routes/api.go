package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	utils "github.com/MikMuellerDev/QuickClip/utils"
)

// Returns the mode which is currently active
func getVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	version, production := utils.GetVersion()
	json.NewEncoder(w).Encode(VersionStruct{Version: version, Production: production})
}

// Returns the clip list
func getClips(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(utils.GetClips())
}

func modClip(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t utils.Clip
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}

	success, clip := utils.ModClip(t.Id, t.Name, t.Content)

	w.Header().Set("Content-Type", "application/json")
	if success {
		json.NewEncoder(w).Encode(clip)
	} else {
		json.NewEncoder(w).Encode(ResponseStruct{false, 404, "Unknown clip", fmt.Sprintf("The ID: %s is not associated with any Object.", t.Id)})
	}
}
