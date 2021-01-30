package repl

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/peterh/liner"

	"github.com/x-research-team/vm"
	"github.com/x-research-team/vm/eval"
	"github.com/x-research-team/vm/lexer"
	"github.com/x-research-team/vm/parser"
)

var kernelKeywords = []string{
	"fn", "let", "true", "false", "if", "else", "elsif", "elseif",
	"elif", "return", "import", "and", "or", "struct", "do", "while",
	"break", "continue", "for", "in", "where", "grep", "map", "case",
	"is", "try", "catch", "finally", "throw", "qw", "unless", "thread",
	"enum", "defer", "nil", "class", "new", "this", "parent", "property",
	"get", "set", "static", "public", "private", "protected", "interface", "default",
	"from", "select", "group", "into", "orderby", "join", "on", "equals", "by", "ascending", "descending",
	"async", "await", "service",
}

//Note: we should put the longest operators first.
var kernelOperators = []string{
	"+=", "-=", "*=", "/=", "%=", "^=",
	"++", "--",
	"&&", "||",
	"<<", ">>",
	"->", "=>",
	"==", "!=", "<=", ">=", "=~", "!~",
	"+", "-", "*", "/", "%", "^",
	"(", ")", "{", "}", "[", "]",
	"=", "<", ">",
	"!", "&", "|", ".",
	",", "?", ":", ";",
}

//var colors = map[liner.Category]string{
//	liner.NumberType:   liner.COLOR_YELLOW,
//	liner.KeywordType:  liner.COLOR_MAGENTA,
//	liner.StringType:   liner.COLOR_CYAN,
//	liner.CommentType:  liner.COLOR_GREEN,
//	liner.OperatorType: liner.COLOR_RED,
//}

const PROMPT = "kernel>> "

func Start(out io.Writer, color bool) {
	history := filepath.Join(os.TempDir(), ".kernel_history")
	l := liner.NewLiner()
	defer l.Close()

	l.SetCtrlCAborts(true)
	l.SetMultiLineMode(true)

	if color {
		//l.SetSyntaxHighlight(color) //use syntax highlight or not
		//l.RegisterKeywords(kernelKeywords)
		//l.RegisterOperators(kernelOperators)
		//l.RegisterColors(colors)
	}

	if f, err := os.Open(history); err == nil {
		l.ReadHistory(f)
		f.Close()
	}

	if color {
		eval.REPLColor = true
	}
	scope := eval.NewScope(vm.Scope)
	wd, err := os.Getwd()
	if err != nil {
		io.WriteString(out, err.Error())
		return
	}

	var tmplines []string
	for {
		if line, err := l.Prompt(PROMPT); err == nil {
			if line == "exit" || line == "quit" {
				if f, err := os.Create(history); err == nil {
					l.WriteHistory(f)
					f.Close()
				}
				break
			}

			tmpline := strings.TrimSpace(line)
			if len(tmpline) == 0 { //empty line
				continue
			}
			//check if the `line` variable is ended with '\'
			if tmpline[len(tmpline)-1:] == "\\" { //the expression/statement has remaining part
				tmplines = append(tmplines, strings.TrimRight(tmpline, "\\"))
				continue
			} else {
				tmplines = append(tmplines, line)
			}

			resultLine := strings.Join(tmplines, "")
			l.AppendHistory(resultLine)
			tmplines = nil // clear the array

			lex := lexer.New("", resultLine)
			p := parser.New(lex, wd)
			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				printParserErrors(out, p.Errors())
				continue
			}

			eval.Eval(program, scope)
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
