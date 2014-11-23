# cheat

*Status: Up for grabs. I would really love to see this created but am busy with other stuff :( If you are interested in learning Go and want to do this, ping me! I'll help you get started.*

Reimplementation of [Chris Lane's cheatsheet](https://github.com/chrisallenlane/cheat) script in [Go](http://golang.org/). 

I'm mostly doing this as a means of learning Go as it seemed like a nice first project to start with. 

# Todo

* Colors on the AppHelpTemplate.

* Edit cheatsheets `cheat --edit git`, if the file does not exist, it is created, also create the dir.

* Copy commands to clipboard, something like `cheat git 12 --copy`
* Or execute a command by, `cheat git 12`

* Allow multiple cheat directories
* User's favorite editor, with support for command line parameters.
* Open the `.cheatrc` for editing via `cheat --config`
* Store all settings in a json file `.cheatrc` in user's home directory.

* List all available sheets. `cheat --list`

* Wrap the output to a fit width? like 79 characters?

* Update cheat sheets from chris' repo, `cheat --update` for updating it the safe way, and `cheat --update --force` for overwriting all the cheats with the downloaded version.

* Should grep support be added? or can that be achived by `grep`ping things?

# Far Future

* Add support for bro pages
* A complete client for http://www.commandlinefu.com/site/api

# Prior Art

* [/chrisallenlane/cheat](http://github.com/chrisallenlane/cheat) in Python
* [/jahendrie/cheat](https://github.com/jahendrie/cheat) in Bash
* [/lucaswerkmeister/cheats](https://github.com/lucaswerkmeister/cheats) in Bash
* [/defunkt/cheat](https://github.com/defunkt/) in Ruby
* [/torsten/cheat](https://github.com/torsten/cheat) in Ruby (single file)
* [/arthurnn/cheatly](https://github.com/arthurnn/cheatly) in Ruby

<!--

Markdown Cheatsheets - https://github.com/rstacruz/cheatsheets
Kapeli's Sheets - https://github.com/Kapeli/cheatsheets
Git Cheat - https://github.com/0xAX/git-cheat
More Sheets - https://github.com/Dmitrii-I/cheat

-->
