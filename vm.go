package vm

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/x-research-team/vm/internal/eval"
	"github.com/x-research-team/vm/internal/lexer"
	"github.com/x-research-team/vm/internal/message"
	"github.com/x-research-team/vm/internal/parser"
)

var Scope *eval.Scope

func RunProgram(debug bool, filename string) {
	Scope = eval.NewScope(nil, os.Stdout)
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	f, err := ioutil.ReadFile(wd + "/" + filename)
	if err != nil {
		fmt.Println("vm: ", err.Error())
		os.Exit(1)
	}

	input := string(f)
	l := lexer.New(filename, input)

	p := parser.New(l, wd)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		for _, err := range p.Errors() {
			fmt.Println(err)
		}
		os.Exit(1)
	}
	RegisterGoGlobals()

	if debug {
		eval.REPLColor = true
		Lines := strings.Split(input, "\n")
		//pre-append an empty line, so the Lines start with 1, not zero.
		Lines = append([]string{""}, Lines...)

		eval.Dbg = eval.NewDebugger(Lines)
		eval.Dbg.SetFunctions(p.Functions)
		eval.Dbg.ShowBanner()

		dbgInfosArr := parser.SplitSlice(parser.DebugInfos) //[][]ast.Node
		eval.Dbg.SetDbgInfos(dbgInfosArr)
		// for idx, dbgInfos := range dbgInfosArr {
		// 	for _, dbgInfo := range dbgInfos {
		// 		fmt.Printf("idx:%d, Line:<%d-%d>, node.Type=%T, node=<%s>\n", idx, dbgInfo.Pos().Line, dbgInfo.End().Line, dbgInfo, dbgInfo.String())
		// 	}
		// }

		eval.MsgHandler = message.NewMessageHandler()
		eval.MsgHandler.AddListener(eval.Dbg)

	}

	result := eval.Eval(program, Scope)
	if result.Type() == eval.ERROR_OBJ {
		fmt.Println(result.Inspect())
	}

	//	e := eval.Eval(program, scope)
	//	if e.Inspect() != "nil" {
	//		fmt.Println(e.Inspect())
	//	}
}

// Run VM program
func Run(debug bool, filename string, s *eval.Scope) eval.Object {
	Scope = eval.NewScope(s, nil)
	var o eval.Object
	root, err := filepath.Abs("..")
	if err != nil {
		log.Printf("[VM] %v\n", err)
		return o
	}
	f, err := ioutil.ReadFile(path.Join(root, filename))
	if err != nil {
		log.Printf("[VM] %v\n", err)
		return o
	}

	input := string(f)
	l := lexer.New(filename, input)

	p := parser.New(l, root)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		for _, err := range p.Errors() {
			log.Printf("[VM] %v\n", err)
		}
		return o
	}

	if debug {
		eval.REPLColor = true
		Lines := strings.Split(input, "\n")
		//pre-append an empty line, so the Lines start with 1, not zero.
		Lines = append([]string{""}, Lines...)

		eval.Dbg = eval.NewDebugger(Lines)
		eval.Dbg.SetFunctions(p.Functions)
		eval.Dbg.ShowBanner()

		dbgInfosArr := parser.SplitSlice(parser.DebugInfos) //[][]ast.Node
		eval.Dbg.SetDbgInfos(dbgInfosArr)

		eval.MsgHandler = message.NewMessageHandler()
		eval.MsgHandler.AddListener(eval.Dbg)
	}

	return eval.Eval(program, Scope)
}

// Register go package methods/types
// Note here, we use 'gfmt', 'glog', 'gos' 'gtime', because in vm
// we already have built in module 'fmt', 'log' 'os', 'time'.
// And Here we demonstrate the use of import go language's methods.
func RegisterGoGlobals() {
	eval.RegisterFunctions("gfmt", map[string]interface{}{
		"Errorf":   fmt.Errorf,
		"Println":  fmt.Println,
		"Print":    fmt.Print,
		"Printf":   fmt.Printf,
		"Fprint":   fmt.Fprint,
		"Fprintln": fmt.Fprintln,
		"Fscan":    fmt.Fscan,
		"Fscanf":   fmt.Fscanf,
		"Fscanln":  fmt.Fscanln,
		"Scan":     fmt.Scan,
		"Scanf":    fmt.Scanf,
		"Scanln":   fmt.Scanln,
		"Sscan":    fmt.Sscan,
		"Sscanf":   fmt.Sscanf,
		"Sscanln":  fmt.Sscanln,
		"Sprint":   fmt.Sprint,
		"Sprintf":  fmt.Sprintf,
		"Sprintln": fmt.Sprintln,
	})

	eval.RegisterFunctions("glog", map[string]interface{}{
		"Fatal":     log.Fatal,
		"Fatalf":    log.Fatalf,
		"Fatalln":   log.Fatalln,
		"Flags":     log.Flags,
		"Panic":     log.Panic,
		"Panicf":    log.Panicf,
		"Panicln":   log.Panicln,
		"Print":     log.Print,
		"Printf":    log.Printf,
		"Println":   log.Println,
		"SetFlags":  log.SetFlags,
		"SetOutput": log.SetOutput,
		"SetPrefix": log.SetPrefix,
	})

	eval.RegisterFunctions("gos", map[string]interface{}{
		"Chdir":    os.Chdir,
		"Chmod":    os.Chmod,
		"Chown":    os.Chown,
		"Exit":     os.Exit,
		"Getpid":   os.Getpid,
		"Hostname": os.Hostname,
		"Environ":  os.Environ,
		"Getenv":   os.Getenv,
		"Setenv":   os.Setenv,
		"Create":   os.Create,
		"Open":     os.Open,
	})

	argsStart := 1
	if len(os.Args) > 2 {
		argsStart = 2
	}
	eval.RegisterVars("gos", map[string]interface{}{
		"Args": os.Args[argsStart:],
	})

	eval.RegisterVars("runtime", map[string]interface{}{
		"GOOS":   runtime.GOOS,
		"GOARCH": runtime.GOARCH,
	})

	eval.RegisterVars("gtime", map[string]interface{}{
		"Duration": time.Duration(0),
		"Ticker":   time.Ticker{},
		"Time":     time.Time{},
	})
	eval.RegisterFunctions("gtime", map[string]interface{}{
		"After":           time.After,
		"Sleep":           time.Sleep,
		"Tick":            time.Tick,
		"Since":           time.Since,
		"FixedZone":       time.FixedZone,
		"LoadLocation":    time.LoadLocation,
		"NewTicker":       time.NewTicker,
		"Date":            time.Date,
		"Now":             time.Now,
		"Parse":           time.Parse,
		"ParseDuration":   time.ParseDuration,
		"ParseInLocation": time.ParseInLocation,
		"Unix":            time.Unix,
		"AfterFunc":       time.AfterFunc,
		"NewTimer":        time.NewTimer,
		"Nanosecond":      time.Nanosecond,
		"Microsecond":     time.Microsecond,
		"Millisecond":     time.Millisecond,
		"Second":          time.Second,
		"Minute":          time.Minute,
		"Hour":            time.Hour,
	})

	eval.RegisterFunctions("math/rand", map[string]interface{}{
		"New":         rand.New,
		"NewSource":   rand.NewSource,
		"Float64":     rand.Float64,
		"ExpFloat64":  rand.ExpFloat64,
		"Float32":     rand.Float32,
		"Int":         rand.Int,
		"Int31":       rand.Int31,
		"Int31n":      rand.Int31n,
		"Int63":       rand.Int63,
		"Int63n":      rand.Int63n,
		"Intn":        rand.Intn,
		"NormFloat64": rand.NormFloat64,
		"Perm":        rand.Perm,
		"Seed":        rand.Seed,
		"Uint32":      rand.Uint32,
	})

	eval.RegisterFunctions("io/ioutil", map[string]interface{}{
		"WriteFile": ioutil.WriteFile,
		"ReadFile":  ioutil.ReadFile,
		"TempDir":   ioutil.TempDir,
		"TempFile":  ioutil.TempFile,
		"ReadAll":   ioutil.ReadAll,
		"ReadDir":   ioutil.ReadDir,
		"NopCloser": ioutil.NopCloser,
	})

	eval.RegisterFunctions("bufio", map[string]interface{}{
		"NewWriter":     bufio.NewWriter,
		"NewReader":     bufio.NewReader,
		"NewReadWriter": bufio.NewReadWriter,
		"NewScanner":    bufio.NewScanner,
	})
	eval.RegisterFunctions("gregex", map[string]interface{}{
		"Match":            regexp.Match,
		"MatchReader":      regexp.MatchReader,
		"MatchString":      regexp.MatchString,
		"QuoteMeta":        regexp.QuoteMeta,
		"Compile":          regexp.Compile,
		"CompilePOSIX":     regexp.CompilePOSIX,
		"MustCompile":      regexp.MustCompile,
		"MustCompilePOSIX": regexp.MustCompilePOSIX,
	})
}
