package main

// These templates help modify the output presented by the cli package.
var AppHelpTemplate = `{{.Name}} - {{.Usage}}

Version: {{.Version}}

Usage:
    {{.Name}} <command> [cheatsheet] [<args>]


Commands:
    {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
    {{end}}

Global Options:
    {{range .Flags}}{{.}}
    {{end}}

Examples:
    {{.Name}} show git              Shows git cheatsheet
    {{.Name}} show git -c 12        Copy the 12th git cheat

    {{.Name}} edit at               Edit cheatsheet named at
`
