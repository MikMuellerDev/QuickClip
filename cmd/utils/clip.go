package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Mode struct {
	Name        string
	Description string
	ImagePath   string
	Url         string
	Id          string
	Volume      int
}

type Clip struct {
	Name    string
	Id      string
	Content string
}

type Clips struct {
	Clips []Clip
}

var clips Clips

func ReadClipFile() {
	path := "../config/clipboard.json"
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	err = json.Unmarshal(content, &clips)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	log.Debug(fmt.Sprintf("Loaded QuickClip Clipboard from %s", path))
}

func GetClips() Clips {
	return clips
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

func GetClipById(id string) *Clip {
	for _, v := range clips.Clips {
		if v.Id == id {
			return &v
		}
	}
	log.Error(fmt.Sprintf("Requested board that does not exist: %s, DoesBoardExist() might have failed.", id))
	return &Clip{"", "", ""}
}

func ModClipInList(clip Clip) {
	var clipsCpy []Clip
	for _, v := range clips.Clips {
		if v.Id == clip.Id {
			clipsCpy = append(clipsCpy, Clip{Name: clip.Name, Id: clip.Id, Content: clip.Content})
		} else {
			clipsCpy = append(clipsCpy, v)
		}
	}
	clips.Clips = clipsCpy
}

func ModClip(id string, name string, content string) (bool, Clip) {
	if DoesClipExist(id) {
		var clip *Clip = GetClipById(id)
		clip.Content = content
		clip.Id = id
		clip.Name = name

		ModClipInList(*clip)

		var jsonBlob = []byte(`{}`)
		err := json.Unmarshal(jsonBlob, &clip)
		if err != nil {
			fmt.Println("opening config file", err.Error())
		}

		clipJson, _ := json.Marshal(clip)
		err = ioutil.WriteFile("output.json", clipJson, 0644)
		fmt.Printf("%s", clipJson)
		return true, Clip{Name: clip.Name, Id: clip.Id, Content: clip.Content}
	} else {
		log.Warn(fmt.Sprintf("The Clip ID: %s does not exist.", id))
	}
	return false, Clip{}
}
