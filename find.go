package main

import (
   "github.com/dgnorton/dmapi"
   "encoding/json"
   "fmt"
   "os"
   "path"
)

var cmdFind = &Command {
   UsageLine:  "find",
   Short:      "search local data store for specified entries",
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
   var findStart string    // -s start date
   var findEnd string      // -e end date
   var findPattern string  // -p regex pattern or *
   var findFormat string   // -f
   var findHTML string     // -html

func addFindFlags(cmd *Command) {
   cmd.Flag.StringVar(&findStart, "s", "", "")
   cmd.Flag.StringVar(&findEnd, "e", "", "")
   cmd.Flag.StringVar(&findPattern, "p", "", "")
   cmd.Flag.StringVar(&findFormat, "format", "", "")
   cmd.Flag.StringVar(&findHTML, "html", "", "")
}

func runFind(cfg *config, cmd *Command, args []string) {
   user := cfg.User
   usrDir, err := userDir(user)
   if err != nil {
      fatalf("%s [find.go - runFind - userDir]", err)
   }

   userFile := path.Join(usrDir, "entries.json")

   haveUserFile, err := isFile(userFile)
   if err != nil {
      fatalf("%s [find.go - runFind - isFile]", err)
   }

   if haveUserFile == false {
      fatalf("No entries to search.  Need to run 'dm sync'?")
   }

   entries, err := dmapi.LoadEntries(userFile)
   if err != nil {
      fatalf("%s [find.go - runFind - dmapi.LoadEntries]", err)
   }

   if len(entries.Entries) == 0 {
      fatalf("No entries to search.  Need to run 'dm sync'?")
   }

   matches, err := entries.Find(findStart, findEnd, findPattern)
   if err != nil {
      fatalf("%s [find.go - runFind - dmapi.Find]", err)
   }

   if findFormat != "" {
      err = fprintText(os.Stdout, matches, findFormat)
      if err != nil {
         fatalf("%s [find.go - runFind - fprintText]", err)
      }
   } else if findHTML != "" {
     err = fprintHTML(os.Stdout, matches, findHTML)
      if err != nil {
         fatalf("%s [find.go - runFind - fprintHTML]", err)
      }
   } else {
       bytes, err := json.Marshal(matches)
       if err != nil {
         fatalf("%s [find.go - runFind - json.Marshal]", err)
       }
       fmt.Fprintf(os.Stdout, "%s", string(bytes))
   }
}

