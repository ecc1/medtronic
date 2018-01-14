package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/ecc1/medtronic"
)

type (
	// Printer represents a function that prints an arbitrary value.
	Printer func(interface{})
)

var (
	formatFlag = flag.String("f", "openaps", "print result in specified `format`")

	format = map[string]Printer{
		"internal": showInternal,
		"json":     showJSON,
		"openaps":  showOpenAPS,
	}

	openAPSMode bool
)

func usage() {
	log.Printf("usage: %s [options] command [ arg ...]", os.Args[0])
	log.Printf("   or: %s [options] command [ args.json ]", os.Args[0])
	flag.PrintDefaults()
	fmts := ""
	for k := range format {
		fmts += " " + k
	}
	log.Printf("output formats:%s", fmts)
	keys := make([]string, len(command))
	i := 0
	for k := range command {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	cmds := ""
	for _, k := range keys {
		cmds += " " + k
	}
	log.Fatalf("commands:%s", cmds)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	printFn := format[*formatFlag]
	if printFn == nil {
		log.Printf("%s: unknown format", *formatFlag)
		usage()
	}
	openAPSMode = *formatFlag == "openaps"
	if flag.NArg() == 0 {
		usage()
	}
	name := flag.Arg(0)
	cmd, found := command[name]
	if !found {
		log.Printf("%s: unknown command", name)
		usage()
	}
	args := getArgs(name, cmd)
	pump := medtronic.Open()
	defer pump.Close()
	pump.Wakeup()
	result := cmd.Cmd(pump, args)
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
	if result == nil {
		return
	}
	printFn(result)
}

type (
	// Arguments represents the formal and actual parameters for a command.
	Arguments map[string]interface{}
)

// Float returns the float value associated with the given key.
func (args Arguments) Float(key string) (float64, error) {
	if openAPSMode {
		f, ok := args[key].(float64)
		if !ok {
			return f, fmt.Errorf("%q parameter must be a number", key)
		}
		return f, nil
	}
	return strconv.ParseFloat(args[key].(string), 64)
}

// Int returns the int value associated with the given key.
func (args Arguments) Int(key string) (int, error) {
	if openAPSMode {
		f, ok := args[key].(float64)
		if !ok {
			return int(f), fmt.Errorf("%q parameter must be a number", key)
		}
		return int(f), nil
	}
	return strconv.Atoi(args[key].(string))
}

// String returns the string value associated with the given key.
func (args Arguments) String(key string) (string, error) {
	if openAPSMode {
		s, ok := args[key].(string)
		if !ok {
			return s, fmt.Errorf("%q parameter must be a string", key)
		}
		return s, nil
	}
	return args[key].(string), nil
}

func getArgs(name string, cmd Command) Arguments {
	params := cmd.Params
	argv := flag.Args()[1:]
	if len(params) == 0 {
		if len(argv) != 0 {
			log.Fatalf("%s does not take any arguments", name)
		}
		return nil
	}
	if openAPSMode {
		return openAPSArgs(name, params, argv)
	}
	return cliArgs(name, params, argv)
}

// Parse an openaps JSON file for arguments.
func openAPSArgs(name string, params []string, argv []string) Arguments {
	if len(argv) != 1 || !strings.HasSuffix(argv[0], ".json") {
		log.Fatalf("%s: openaps format requires single JSON argument file", name)
	}
	file := argv[0]
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("%s: %v", name, err)
	}
	args := make(Arguments)
	err = json.NewDecoder(f).Decode(&args)
	if err != nil {
		log.Fatalf("%s: %v", name, err)
	}
	_ = f.Close()
	for _, k := range params {
		_, present := args[k]
		if !present {
			log.Fatalf("%s: argument file %s is missing %q parameter", name, file, k)
		}
	}
	return args
}

// Collect command-line arguments.
func cliArgs(name string, params []string, argv []string) Arguments {
	args := make(Arguments)
	for i, k := range params {
		if i >= len(argv) {
			// Bind remaining parameters to "".
			args[k] = ""
			continue
		}
		if strings.HasSuffix(k, "...") {
			// Bind all remaining args to this parameter.
			args[k] = argv[i:]
			if i != len(params)-1 {
				panic(k + " is not the final parameter")
			}
			continue
		}
		args[k] = argv[i]
	}
	return args
}
