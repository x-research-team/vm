package vm

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/x-research-team/vm/highlight"
)

func _main_highlight() {
	args := os.Args[1:]

	var outFileName string = "output.html"
	var f []byte
	var err error

	if len(args) == 0 {
		f, err = ioutil.ReadAll(os.Stdin)
	} else {
		f, err = ioutil.ReadFile(args[0])
		if err != nil {
			log.Printf("[VM] Highlighter: cannot read file; %v\n", err.Error())
			return
		}
		outFileName = fmt.Sprintf("%s.html", args[0])
	}

	highlighter := highlight.New(string(f))
	highlighter.RegisterGenerator(&highlight.ConsoleHighlighter{})

	file, err := os.Create(outFileName)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	highlighter.RegisterGenerator(&highlight.HtmlHighlighter{Out: file})
	highlighter.Highlight()
}
