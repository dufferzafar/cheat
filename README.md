# cheat

Reimplementation of [Chris Lane's cheatsheet](https://github.com/chrisallenlane/cheat) script in [Go](http://golang.org/). 

I'm mostly doing this as a means of learning Go as it seemed like a nice first project to start with. 

I'll update the readme with more fodder later.

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

# Shout outs

Thanks to [Chris Lane](http://github.com/chrisallenlane/) for creating cheat, and to [Jahendrie](https://github.com/jahendrie/) and [Lucas Werkmeister](https://github.com/lucaswerkmeister) for writing their own personal versions of cheat (I've stolen most of the features from them.)
