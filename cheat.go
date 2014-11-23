// Create and view command-line cheatsheets.
// A Go reimplementation of Chris Lane's python script.
package main

import (
    "bufio"
    "fmt"
    "github.com/codegangsta/cli"
    "io/ioutil"
    "os"
    "os/exec"
    "path"
    "strconv"
    "strings"
)

const version string = "0.3"

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

    app.Flags = []cli.Flag{
        cli.BoolFlag{"edit, e", "edit a cheatsheet"},
    }

    app.Run(os.Args)
}

func mainCmd(c *cli.Context) {
    // Read settings
    config := &JSONData{}
    config.ReadConfig()

    if !c.Args().Present() {
        cli.ShowAppHelp(c)
    } else {
        var cmdname = c.Args().First()
        var cheatfile = path.Join(cheatdir, cmdname)

        if c.Bool("edit") {
            editCheat(cheatfile, config.Editor)
        } else {
            if _, err := os.Stat(cheatfile); os.IsNotExist(err) {
                fmt.Fprintf(os.Stderr, "No cheatsheat found for '%s'\n", cmdname)
                fmt.Fprintf(os.Stderr, "To create a new sheet, run: cheat -e %s\n", cmdname)
                os.Exit(1)
            } else {
                showCheats(cheatfile, cmdname)
            }
        }
    }
}

func showCheats(cheatfile string, cmdname string) {
    file, _ := os.Open(cheatfile)
    scanner := bufio.NewScanner(file)

    var i = 1
    for scanner.Scan() {
        line := strings.Trim(scanner.Text(), " ")

        // Pretty print the output
        // Todo: Will have to be tested on other platforms and terminals.

        if strings.HasPrefix(line, "#") {
            fmt.Println("\x1b[33;1m" + line + "\x1b[0m")
        } else if strings.HasPrefix(line, cmdname) {
            fmt.Println("\x1b[36;1m(" + strconv.Itoa(i) + ")\x1b[0m " + line)
            i++
        } else {
            fmt.Println(line)
        }
    }
    file.Close()
}

func editCheat(cheatfile string, configEditor string) {
    editor, err := exec.LookPath(configEditor)

    if err != nil {
        fmt.Fprintf(os.Stderr, "Editor not found: %s", editor)
    }

    cmd := exec.Command(editor, cheatfile)
    cmd.Stdout = os.Stdout
    cmd.Stdin = os.Stdin
    cmd.Stderr = os.Stderr
    cmd.Run()
}
