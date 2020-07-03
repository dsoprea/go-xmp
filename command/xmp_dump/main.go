package main

import (
	"fmt"
	"os"

	"encoding/json"

	"github.com/dsoprea/go-logging"
	"github.com/jessevdk/go-flags"

	"github.com/dsoprea/go-xmp"
	_ "github.com/dsoprea/go-xmp/namespace"
)

var (
	mainLogger = log.NewLogger("main.main")
)

type parameters struct {
	Filepath            string `short:"f" long:"filepath" required:"true" description:"File-path of image"`
	PrintAsJson         bool   `short:"j" long:"json" description:"Print out as JSON"`
	IsVerbose           bool   `short:"v" long:"verbose" description:"Print logging"`
	DoNotSimplifyExport bool   `short:"n" long:"no-simplify" description:"If exporting, return the raw, unsimplified structure"`
}

var (
	arguments = new(parameters)
)

func main() {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err := errRaw.(error)
			log.PrintError(err)

			os.Exit(-2)
		}
	}()

	_, err := flags.Parse(arguments)
	if err != nil {
		os.Exit(-1)
	}

	if arguments.IsVerbose == true {
		cla := log.NewConsoleLogAdapter()
		log.AddAdapter("console", cla)

		scp := log.NewStaticConfigurationProvider()
		scp.SetLevelName(log.LevelNameDebug)

		log.LoadConfiguration(scp)
	}

	f, err := os.Open(arguments.Filepath)
	log.PanicIf(err)

	defer f.Close()

	xp := xmp.NewParser(f)

	xpi, err := xp.Parse()
	log.PanicIf(err)

	if arguments.PrintAsJson == true {
		doSimplify := arguments.DoNotSimplifyExport == false

		flat, err := xpi.Export(doSimplify)
		log.PanicIf(err)

		encoded, err := json.MarshalIndent(flat, "", "  ")
		log.PanicIf(err)

		fmt.Println(string(encoded))

		return
	}

	xpi.Dump()
}
