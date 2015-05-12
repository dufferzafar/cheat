// Create and view command-line cheatsheets.
// A Go reimplementation of Chris Lane's python script.
package main

import (
	"bufio"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

const version string = "0.5"

// These settings will be stored in a json once we get to that point.
const cheatdir string = "/mnt/Work/Github/cheat/cheatsheets"

func main() {
	app := cli.NewApp()

	// Use our custom templates
	cli.AppHelpTemplate = AppHelpTemplate

	app.Name = "cheat"
	app.Usage = "Create and view command-line cheatsheets."
	app.Version = version

	config := &JSONData{}
	config.ReadConfig()

	app.Commands = []cli.Command{
		{
			Name:    "show",
			Aliases: []string{"s"},
			Usage:   "Show cheats related to a command",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "copy, c",
					Usage: "cheat number to copy",
				},
			},
			Action: func(c *cli.Context) {
				var cmdname = c.Args().First()
				var cheatfile = path.Join(cheatdir, cmdname)

				if _, err := os.Stat(cheatfile); os.IsNotExist(err) {
					fmt.Fprintf(os.Stderr, "No cheatsheat found for '%s'\n", cmdname)
					fmt.Fprintf(os.Stderr, "To create a new sheet, run: cheat edit %s\n", cmdname)
					os.Exit(1)
				} else {
					if c.Int("copy") != 0 {
						copyCheat(cheatfile, cmdname, c.Int("copy"))
					} else {
						showCheats(cheatfile, cmdname)
					}
				}
			},
		},
		{
			Name:    "edit",
			Aliases: []string{"e"},
			Usage:   "Add/Edit a cheat",
			Action: func(c *cli.Context) {
				var cheatfile = path.Join(cheatdir, c.Args().First())
				editCheat(cheatfile, config.Editor)
			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "List all available cheats",
			Action: func(c *cli.Context) {
				files, _ := ioutil.ReadDir(cheatdir)
				for _, f := range files {
					fmt.Println(f.Name())
				}
			},
		},
	}

	app.Run(os.Args)
}

func copyCheat(cheatfile string, cmdname string, cheatno int) {
	file, _ := os.Open(cheatfile)
	scanner := bufio.NewScanner(file)

	var i = 0
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " ")

		if strings.HasPrefix(line, cmdname) {
			i++
		}

		if cheatno == i {
			clipboard.WriteAll(line)
			fmt.Println("\x1b[32;5m" + "Copied to Clipboard: " + "\x1b[0m" + line)
			break
		}
	}
	file.Close()
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
			fmt.Println("\x1b[33;5m" + line + "\x1b[0m")
		} else if strings.HasPrefix(line, cmdname) {
			fmt.Println("\x1b[36;5m(" + strconv.Itoa(i) + ")\x1b[0m " + line)
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
