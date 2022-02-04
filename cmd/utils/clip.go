package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
var prevSave []byte
var saveCounter = 10

func ReadClipFile() {
	path := "../config/clipboard.json"
	content, err := ioutil.ReadFile(path)
	if err != nil {
		err = ioutil.WriteFile("../config/clipboard.json", []byte("{}"), 0644)
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

func DoesClipExist(id string) bool {
	for _, v := range clips.Clips {
		if v.Id == id {
			return true
		}
	}
	log.Trace(fmt.Sprintf("Requested board that does not exist %s", id))
	return false
}

func GetClipById(id string, user string) (bool, Clip) {
	for _, v := range clips.Clips {
		if v.Id == id {
			if v.Restricted {
				if HasPermission(user, v.Id) || user == "admin" {
					return true, v
				} else {
					log.Warn(fmt.Sprintf("User: %s requested restricted board: %s", user, id))
					return false, Clip{"", "", "", "", false, false, -1, false}
				}
			} else {
				return true, v
			}
		}
	}
	log.Error(fmt.Sprintf("Requested board that does not exist: %s, DoesBoardExist() might have failed.", id))
	return false, Clip{"", "", "", "", false, false, -1, false}
}

func ModClipInList(clip Clip) {
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
	fmt.Printf("Removing clip with id: %s\n", id)
	var clipsCpy []Clip
	// _, clip := GetClipById(id, "admin")
	for _, v := range clips.Clips {
		if id != v.Id {
			clipsCpy = append(clipsCpy, v)
		}
	}
	clips.Clips = clipsCpy
	writeClips(clips)
}

func AddClip(clip Clip) {
	clips.Clips = append(clips.Clips, clip)
	writeClips(clips)
}

func ModClip(clip Clip) (bool, Clip) {
	if DoesClipExist(clip.Id) {
		ModClipInList(clip)
		writeClips(clips)
		return true, Clip{Name: clip.Name, Id: clip.Id, Description: clip.Description, Content: clip.Content, Restricted: clip.Restricted, Refresh: clip.Refresh, RefreshInterval: clip.RefreshInterval, ReadOnly: clip.ReadOnly}
	} else {
		log.Warn(fmt.Sprintf("The Clip ID: %s does not exist.", clip.Id))
	}
	return false, Clip{}
}

func writeClips(clips Clips) {
	var jsonBlob = []byte(`{}`)
	err := json.Unmarshal(jsonBlob, &clips)
	if err != nil {
		log.Fatal("Error during unmarshal", err.Error())
	}

	clipJson, _ := json.MarshalIndent(clips, "", "    ")
	prevSave = clipJson
	err = ioutil.WriteFile("../config/clipboard.json", clipJson, 0644)
	if err != nil {
		log.Fatal("Error writing clipboard: %s", err.Error())
	}
	log.Debug("Written clip contents to clipboard.json.")
}

func RequestSave() bool {
	// First get a byte arr from the current state
	var jsonBlob = []byte(`{}`)
	err := json.Unmarshal(jsonBlob, &clips)
	if err != nil {
		log.Fatal("Error during unmarshal", err.Error())
	}

	clipJson, _ := json.Marshal(clips)

	if string(prevSave) == string(clipJson) {
		log.Trace("No save was triggered: identical states.")
		return false
	} else {
		log.Trace("Saving changes due to changes.")
		writeClips(clips)
		prevSave = clipJson
		return true
	}
}

func EditClip(id string, content string, user string) bool {
	_, clip := GetClipById(id, "admin")
	if clip.Restricted {
		if !HasPermission(user, id) && user != "admin" {
			return false
		}
	}
	var lenBefore int = 0
	var clipsCpy []Clip

	for _, v := range clips.Clips {
		if v.Id == id {
			lenBefore = len(v.Content)
			clipsCpy = append(clipsCpy, Clip{Name: v.Name, Id: v.Id, Content: content, Description: v.Description, Restricted: v.Restricted, Refresh: v.Refresh, RefreshInterval: v.RefreshInterval, ReadOnly: v.ReadOnly})
		} else {
			clipsCpy = append(clipsCpy, v)
		}
	}

	clips.Clips = clipsCpy
	if len(content)-lenBefore > 1 {
		saveCounter = 10
		writeClips(clips)
	} else if saveCounter <= 0 {
		writeClips(clips)
		saveCounter = 10
	}
	saveCounter--
	log.Trace(fmt.Sprintf("Save Counter at: %d", saveCounter))
	return true
}
