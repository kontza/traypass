/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"embed"
	"waitwhat/cmd"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	cmd.Execute(assets)
}
