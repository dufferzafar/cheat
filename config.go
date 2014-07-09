package main

import (
    "encoding/json"
    "io/ioutil"
    "os"
    "os/user"
    "path"
)

type JSONData struct {
    Highlight bool     `json:"highlight"`
    Linewrap  int      `json:"linewrap"`
    Editor    string   `json:"editor"`
    Cheatdirs []string `json:"cheatdirs"`
}

var defaults = []byte(`{
    "highlight": true,
    "linewrap": 79,
    "cheatdirs": [
        "~/.cheatsheets"
    ],
    "editor": "vim"
}`)

func (q *JSONData) ReadConfig() error {
    usr, _ := user.Current()
    rcfile := path.Join(usr.HomeDir, ".cheatrc")

    settings := defaults
    if _, err := os.Stat(rcfile); os.IsNotExist(err) {
        ioutil.WriteFile(rcfile, []byte(defaults), 0777)
    } else {
        settings, _ = ioutil.ReadFile(rcfile)
    }

    //Umarshalling JSON into struct
    var data = &q
    return json.Unmarshal(settings, data)
}
