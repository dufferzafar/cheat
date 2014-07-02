package main

// These templates help modify the output presented by the cli package.
var AppHelpTemplate = `{{.Name}} - {{.Usage}}

Version: {{.Version}}

Usage:
    {{.Name}} [global options] [cheatsheet] <command> [<args>]


Commands:
    {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
    {{end}}

Global Options:
    {{range .Flags}}{{.}}
    {{end}}

Examples:
    {{.Name}} git           Shows git cheatsheet

    {{.Name}} -e at         Create a new cheatsheet named at
`
