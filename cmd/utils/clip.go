package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

type Clip struct {
	Name            string
	Id              string
	Content         string
	Description     string
	Restricted      bool
	Refresh         bool
	RefreshInterval int
	ReadOnly        bool
}

type Clips struct {
	Clips []Clip
}

var clips Clips
var clipsMutex sync.RWMutex
var prevSave []byte
var saveCounter = 10

func ReadClipFile() {
	clipsMutex.Lock()
	defer clipsMutex.Unlock()
	path := "../config/clipboard.json"
	content, err := ioutil.ReadFile(path)
	if err != nil {
		err = ioutil.WriteFile("../config/clipboard.json", []byte("{}"), 0600)
		if err != nil {
			log.Fatal("Error when opening file: ", err)
			return
		}
		return
	}
	err = json.Unmarshal(content, &clips)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	log.Debug(fmt.Sprintf("Loaded QuickClip Clipboard from %s", path))
}

func GetClips(username string) Clips {
	clipsMutex.RLock()
	defer clipsMutex.RUnlock()
	var clipsCpy Clips
	for _, v := range clips.Clips {
		var clip = Clip{Name: v.Name, Id: v.Id, Content: "", Restricted: v.Restricted, Description: v.Description, Refresh: v.Refresh, RefreshInterval: v.RefreshInterval, ReadOnly: v.ReadOnly}
		if v.Restricted {
			if HasPermission(username, v.Id) || username == "admin" {
				clipsCpy.Clips = append(clipsCpy.Clips, clip)
			}
		} else {
			clipsCpy.Clips = append(clipsCpy.Clips, clip)
		}
	}
	return clipsCpy
}

// Caller must hold clipsMutex
func clipExists(id string) bool {
	for _, v := range clips.Clips {
		if v.Id == id {
			return true
		}
	}
	return false
}

func DoesClipExist(id string) bool {
	clipsMutex.RLock()
	defer clipsMutex.RUnlock()
	if clipExists(id) {
		return true
	}
	log.Trace(fmt.Sprintf("Requested board that does not exist %q", id))
	return false
}

func GetClipById(id string, user string) (bool, Clip) {
	clipsMutex.RLock()
	defer clipsMutex.RUnlock()
	for _, v := range clips.Clips {
		if v.Id == id {
			if v.Restricted {
				if HasPermission(user, v.Id) || user == "admin" {
					return true, v
				} else {
					log.Warn(fmt.Sprintf("User: %q requested restricted board: %q", user, id))
					return false, Clip{"", "", "", "", false, false, -1, false}
				}
			} else {
				return true, v
			}
		}
	}
	log.Debug(fmt.Sprintf("Requested board that does not exist: %q", id))
	return false, Clip{"", "", "", "", false, false, -1, false}
}

// Caller must hold clipsMutex
func modClipInList(clip Clip) {
	var clipsCpy []Clip
	for _, v := range clips.Clips {
		if v.Id == clip.Id {
			clipsCpy = append(clipsCpy, Clip{Name: clip.Name, Id: v.Id, Content: v.Content, Description: clip.Description, Restricted: clip.Restricted, Refresh: clip.Refresh, RefreshInterval: clip.RefreshInterval, ReadOnly: clip.ReadOnly})
		} else {
			clipsCpy = append(clipsCpy, v)
		}
	}
	clips.Clips = clipsCpy
}

func RemoveClip(id string) {
	clipsMutex.Lock()
	defer clipsMutex.Unlock()
	log.Info(fmt.Sprintf("Removing clip with id: %q", id))
	var clipsCpy []Clip
	for _, v := range clips.Clips {
		if id != v.Id {
			clipsCpy = append(clipsCpy, v)
		}
	}
	clips.Clips = clipsCpy
	writeClips()
}

// Returns false if a clip with the same id already exists
func AddClip(clip Clip) bool {
	clipsMutex.Lock()
	defer clipsMutex.Unlock()
	if clipExists(clip.Id) {
		return false
	}
	clips.Clips = append(clips.Clips, clip)
	writeClips()
	return true
}

func ModClip(clip Clip) (bool, Clip) {
	clipsMutex.Lock()
	defer clipsMutex.Unlock()
	if clipExists(clip.Id) {
		modClipInList(clip)
		writeClips()
		return true, Clip{Name: clip.Name, Id: clip.Id, Description: clip.Description, Content: clip.Content, Restricted: clip.Restricted, Refresh: clip.Refresh, RefreshInterval: clip.RefreshInterval, ReadOnly: clip.ReadOnly}
	} else {
		log.Warn(fmt.Sprintf("The Clip ID: %q does not exist.", clip.Id))
	}
	return false, Clip{}
}

// Caller must hold clipsMutex. Write errors are logged instead of stopping the server,
// the in-memory state stays intact and is written again on the next save.
func writeClips() bool {
	clipJson, err := json.MarshalIndent(clips, "", "    ")
	if err != nil {
		log.Error("Error during marshal: ", err.Error())
		return false
	}
	err = ioutil.WriteFile("../config/clipboard.json", clipJson, 0600)
	if err != nil {
		log.Error(fmt.Sprintf("Error writing clipboard: %s", err.Error()))
		return false
	}
	prevSave = clipJson
	log.Debug("Written clip contents to clipboard.json.")
	return true
}

func RequestSave() bool {
	clipsMutex.Lock()
	defer clipsMutex.Unlock()
	clipJson, _ := json.MarshalIndent(clips, "", "    ")

	if string(prevSave) == string(clipJson) {
		log.Trace("No save was triggered: identical states.")
		return false
	} else {
		log.Trace("Saving changes due to changes.")
		return writeClips()
	}
}

func EditClip(id string, content string, user string) bool {
	clipsMutex.Lock()
	defer clipsMutex.Unlock()
	var lenBefore int = 0
	var clipsCpy []Clip

	for _, v := range clips.Clips {
		if v.Id == id {
			if v.Restricted && !HasPermission(user, id) {
				return false
			}
			lenBefore = len(v.Content)
			clipsCpy = append(clipsCpy, Clip{Name: v.Name, Id: v.Id, Content: content, Description: v.Description, Restricted: v.Restricted, Refresh: v.Refresh, RefreshInterval: v.RefreshInterval, ReadOnly: v.ReadOnly})
		} else {
			clipsCpy = append(clipsCpy, v)
		}
	}

	clips.Clips = clipsCpy
	if len(content)-lenBefore > 1 {
		saveCounter = 10
		writeClips()
	} else if saveCounter <= 0 {
		writeClips()
		saveCounter = 10
	}
	saveCounter--
	log.Trace(fmt.Sprintf("Save Counter at: %d", saveCounter))
	return true
}
