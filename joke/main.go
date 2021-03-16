package main

import (
	_ "github.com/anders-14/cmdtab-joke"
	"github.com/rwxrob/cmdtab"
)

func main() {
	cmdtab.Execute("joke")
}
