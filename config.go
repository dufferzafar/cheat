package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type JSONData struct {
	Highlight bool     `json:"highlight"`
	Linewrap  int      `json:"linewrap"`
	Editor    string   `json:"editor"`
	Cheatdirs []string `json:"cheatdirs"`
}

var defaults = `{
    "highlight": true,
    "linewrap": 79,
    "cheatdirs": [
        "$HomeDir/.cheatsheets"
    ],
    "editor": "vim"
}`

func (q *JSONData) ReadConfig() error {
	usr, _ := user.Current()
	rcfile := filepath.Join(usr.HomeDir, ".cheatrc")

	settings := []byte(defaults)
	if _, err := os.Stat(rcfile); os.IsNotExist(err) {
		defaults = strings.Replace(defaults, "$HomeDir", usr.HomeDir, 1)
		ioutil.WriteFile(rcfile, []byte(defaults), 0777)
	} else {
		settings, _ = ioutil.ReadFile(rcfile)
	}

	//Umarshalling JSON into struct
	var data = &q
	err := json.Unmarshal(settings, data)
	if err != nil {
		return err
	}
	for i, dir := range q.Cheatdirs {
		if strings.HasPrefix(dir, "~/") {
			q.Cheatdirs[i] = filepath.Join(usr.HomeDir, dir[2:])
		}
	}
	return nil
}
