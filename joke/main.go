package main

import (
	"github.com/rwxrob/cmdtab"
	_ "github.com/anders-14/cmdtab-joke"
)

func main() {
	cmdtab.Execute("joke")
}
