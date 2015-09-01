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
	"net/url"
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
			Usage: "Fetch cheats from Github",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "dir, d",
					Value: config.Cheatdirs[0],
					Usage: "cheats directory",
				},
				cli.StringFlag{
					Name:  "repo, r",
					Value: "https://github.com/chrisallenlane/cheat/cheat/cheatsheets",
					Usage: "repository to fetch cheats from",
				},
				cli.StringFlag{
					Name:  "local, l",
					Usage: "local path to store repository",
				},
			},
			Action: func(c *cli.Context) {
				if c.String("local") == "" && os.Getenv("GOPATH") == "" {
					fmt.Fprintf(os.Stderr, "Local path to store repo is required.\n")
					return
				}

				fetchCheats(c)
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
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, cmdname) {
			i++
		}

		if cheatno == i {
			re := regexp.MustCompile(`([^#]*)`)
			res := re.FindAllStringSubmatch(line, -1)
			line = strings.TrimSpace(res[0][0])
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
		line := strings.TrimSpace(scanner.Text())

		// Pretty print the output
		// Todo: Will have to be tested on other platforms and terminals.

		if strings.HasPrefix(line, "#") {
			fmt.Fprintln(stdout, "\x1b[33;5m"+line+"\x1b[0m")
		} else if len(line) > 0 {
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

func fetchCheats(c *cli.Context) {
	// parse repo url
	repo, err := url.Parse(c.String("repo"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	repoPath := strings.Split(repo.Path, "/")
	if len(repoPath) <= 3 {
		fmt.Fprintln(os.Stderr, "Invalid Repo URL")
		return
	}

	cheatsPath := repoPath[3:]
	repo.Path = fmt.Sprintf("/%s/%s", repoPath[1], repoPath[2])

	// directory where the cloned repository is stored
	var cloneDir string
	if c.String("local") != "" {
		cloneDir = c.String("local")
	} else if os.Getenv("GOPATH") != "" {
		cloneDir = filepath.Join(os.Getenv("GOPATH"), "src", repo.Host, repoPath[1], repoPath[2])
	}

	// update the repo
	updated, err := updateLocalRepo(repo.String(), cloneDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	// copy cheats
	if updated {
		srcPath := cloneDir
		for _, p := range cheatsPath {
			srcPath = filepath.Join(srcPath, p)
		}

		count, err := copyCheatFiles(srcPath, c.String("dir"))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		fmt.Fprintf(os.Stderr, "%d cheats updated.\n", count)
	} else {
		fmt.Fprintf(os.Stderr, "No cheats updated.\n")
	}
}

func updateLocalRepo(url, dir string) (bool, error) {
	var cmd *exec.Cmd

	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		cmd = exec.Command("git", "clone", url, dir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return true, cmd.Run()
	} else {
		cmd = exec.Command("git", "pull", url)
		cmd.Dir = dir

		out, err := cmd.CombinedOutput()
		if err != nil {
			return false, err
		}

		res := string(out)
		fmt.Fprint(os.Stderr, res)

		updated := true
		if strings.Contains(res, "Already up-to-date.") {
			updated = false
		}

		return updated, nil
	}
}

func copyCheatFiles(cloneDir, cheatsDir string) (int, error) {
	fmt.Fprintf(os.Stderr, "Copying from %s to %s\n", cloneDir, cheatsDir)
	count := 0

	files, err := ioutil.ReadDir(cloneDir)
	if err != nil {
		return count, err
	}

	err = os.MkdirAll(cheatsDir, 0755)
	if err != nil {
		return count, err
	}

	for _, f := range files {
		count += 1

		err := copyFile(filepath.Join(cloneDir, f.Name()), filepath.Join(cheatsDir, f.Name()))
		if err != nil {
			return count, err
		}
	}

	return count, nil
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
	defer d.Close()

	if _, err := io.Copy(d, s); err != nil {
		return err
	}

	return nil
}
