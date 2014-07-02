// Create and view command-line cheatsheets.
// A Go reimplementation of Chris Lane's python script.
package main

import (
    "fmt"
    "github.com/codegangsta/cli"
    "io/ioutil"
    "os"
    "path"
)

const version string = "0.1"

// These settings will be stored in a json once we get to that point.
const cheatdir string = "cheatsheets"

func main() {
    app := cli.NewApp()

    // Use our custom templates
    cli.AppHelpTemplate = AppHelpTemplate

    app.Name = "cheat"
    app.Usage = "Create and view command-line cheatsheets."
    app.Version = version

    app.Action = mainCmd

    app.Run(os.Args)
}

func mainCmd(c *cli.Context) {
    if !c.Args().Present() {
        // Help will be displayed here
        println("No Args were passed!")
    } else {
        var cmdname = os.Args[1]
        var cheatfile = path.Join(cheatdir, cmdname)

        if _, err := os.Stat(cheatfile); os.IsNotExist(err) {
            fmt.Printf("No cheatsheat found for '%s'\n", cmdname)
            fmt.Printf("To create a new sheet, run: cheat -e %s\n", cmdname)
            return
        } else {
            data, _ := ioutil.ReadFile(cheatfile)
            fmt.Print(string(data))
        }
    }
}
