package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgnorton/dmapi"
	"github.com/dgnorton/nltime"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"
)

var cmdFind = &Command{
	UsageLine: "find [-s start date] [-e end date] [-p regex] [-f output template file | json]",
	Short:     "search local data store for specified entries",
	Long: `
Searchs the local data for an entries matching the specified
parameters.  If only a start date is given, all entries from that
time forward will be shown.  If no start or end dates are given,
all entries will be considered. The pattern is a regex pattern
with the exception that giving only "*" as the pattern will
match anything.`,
}

func init() {
	cmdFind.Run = runFind
	addFindFlags(cmdFind)
}

// Add command line flags specific to the find command here
var findStart string   // -s start date
var findEnd string     // -e end date
var findPattern string // -p regex pattern or *
var findFormat string  // -f

func addFindFlags(cmd *Command) {
	cmd.Flag.StringVar(&findStart, "s", "", "")
	cmd.Flag.StringVar(&findEnd, "e", "", "")
	cmd.Flag.StringVar(&findPattern, "p", "", "")
	cmd.Flag.StringVar(&findFormat, "f", "entries.tsv", "")
}

func runFind(cfg *config, cmd *Command, args []string) {
	user := cfg.User
	usrDir, err := userDir(user)
	if err != nil {
		log.Fatalf("%s", err)
	}

	userFile := path.Join(usrDir, "entries.json")

	haveUserFile, err := isFile(userFile)
	if err != nil {
		log.Fatalf("%s", err)
	}

	if haveUserFile == false {
		log.Fatalf("No entries to search.  Need to run 'dm sync'?")
	}

	entries, err := dmapi.LoadEntries(userFile)
	if err != nil {
		log.Fatalf("%s", err)
	}

	if len(entries.Entries) == 0 {
		log.Fatalf("No entries to search.  Need to run 'dm sync'?")
	}

	if len(args) > 0 {
		args, params, err := parseNaturalLangArgs(args)
		if err != nil {
			log.Fatalf("%s", err)
		}

		cmd.Flag.Parse(args)
		findStart = params.StartDate
		findEnd = params.EndDate
	}

	matches, err := entries.Find(findStart, findEnd, findPattern)
	if err != nil {
		log.Fatalf("%s", err)
	}

	if findFormat != "json" {
		var wr io.Writer = os.Stdout

		if filepath.Ext(findFormat) == "tsv" {
			wr = tabwriter.NewWriter(wr, 0, 8, 1, '\t', 0)
		}

		if filepath.Ext(findFormat) == "html" {
			err = fprintHTML(wr, matches, findFormat)
		} else {
			err = fprintText(wr, matches, findFormat)
		}

		if err != nil {
			log.Fatalf("%s", err)
		}
	} else {
		bytes, err := json.Marshal(matches)
		if err != nil {
			log.Fatalf("%s", err)
		}
		fmt.Fprintf(os.Stdout, "%s", string(bytes))
	}
}

type findParams struct {
	StartDate, EndDate, Pattern string
}

func parseNaturalLangArgs(args []string) ([]string, findParams, error) {
	var nlargs, otherArgs []string
	for i, arg := range args {
		if arg[0] == '-' {
			otherArgs = args[i:]
			break
		}
		nlargs = append(nlargs, arg)
	}
	d, err := nltime.ParseRange(strings.Join(nlargs, " "), time.Monday)
	if err != nil {
		return nil, findParams{}, err
	}
	return otherArgs, findParams{d[0].Format("06/1/2"), d[1].Format("06/1/2"), ""}, nil
}
