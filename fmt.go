package vm

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/x-research-team/vm/formatter"
)

func _main_fmt() {
	args := os.Args[1:]

	var f []byte
	var err error
	if len(args) == 0 {
		f, err = ioutil.ReadAll(os.Stdin)
	} else {
		f, err = ioutil.ReadFile(args[0])
		if err != nil {
			log.Printf("[VM] Formatter: cannot read file; %v\n", err.Error())
			return
		}
	}

	formatter := formatter.New(string(f))
	formatter.Format()
}
