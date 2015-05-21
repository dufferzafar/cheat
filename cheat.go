// Create and view command-line cheatsheets.
// A Go reimplementation of Chris Lane's python script.
package main

import (
	"bufio"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/codegangsta/cli"
	"github.com/mattn/go-colorable"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const version string = "0.5"

var (
	stdout = colorable.NewColorableStdout()
)

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
				var cheatfile = filepath.Join(config.Cheatdirs[0], cmdname)

				if _, err := os.Stat(cheatfile); os.IsNotExist(err) {
					fmt.Fprintf(os.Stderr, "No cheatsheet found for '%s'\n", cmdname)
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
				var cheatfile = filepath.Join(config.Cheatdirs[0], c.Args().First())
				editCheat(cheatfile, config.Editor)
			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "List all available cheats",
			Action: func(c *cli.Context) {
				files, _ := ioutil.ReadDir(config.Cheatdirs[0])
				for _, f := range files {
					fmt.Println(f.Name())
				}
			},
		},
		{
			Name:  "config",
			Usage: "Edit the config file",
			Action: func(c *cli.Context) {
				usr, _ := user.Current()
				rcfile := filepath.Join(usr.HomeDir, ".cheatrc")
				editCheat(rcfile, config.Editor)
			},
		},
		{
			Name:  "fetch",
			Usage: "Download cheats",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "dir, d",
					Value: config.Cheatdirs[0],
					Usage: fmt.Sprintf("cheats directory (default: %s)", config.Cheatdirs[0]),
				},
				cli.BoolFlag{
					Name: "verbose, v",
				},
			},
			Action: func(c *cli.Context) {
				url := "https://github.com/chrisallenlane/cheat"
				exitCode := download(url, c.String("dir"), c.Bool("verbose"))
				os.Exit(exitCode)
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
			re := regexp.MustCompile(`([^#]*)`)
			res := re.FindAllStringSubmatch(line, -1)
			line = strings.Trim(res[0][0], " ")
			clipboard.WriteAll(line)
			fmt.Fprintln(stdout, "\x1b[32;5m"+"Copied to Clipboard: "+"\x1b[0m"+line)
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
			fmt.Fprintln(stdout, "\x1b[33;5m"+line+"\x1b[0m")
		} else if strings.HasPrefix(line, cmdname) {
			fmt.Fprintln(stdout, "\x1b[36;5m("+strconv.Itoa(i)+")\x1b[0m "+line)
			i++
		} else {
			fmt.Fprintln(stdout, line)
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

func download(url string, cheatsDir string, verbose bool) int {
	cloneDir, err := ioutil.TempDir(os.TempDir(), "cheat")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	defer func() {
		if verbose {
			fmt.Fprintf(os.Stderr, "Removing temporary directory: %s\n", cloneDir)
		}
		os.RemoveAll(cloneDir)
	}()

	if runGitClone(url, cloneDir, verbose) != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	if copyCheatFiles(path.Join(cloneDir, "cheat", "cheatsheets"), cheatsDir) != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	return 0
}

func runGitClone(url, dir string, verbose bool) error {
	cmd := exec.Command("git", "clone", url, dir)
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd.Run()
}

func copyCheatFiles(cloneDir, cheatsDir string) error {
	files, err := ioutil.ReadDir(cloneDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		err := copyFile(path.Join(cloneDir, f.Name()), path.Join(cheatsDir, f.Name()))
		if err != nil {
			return err
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}

	return d.Close()
}
