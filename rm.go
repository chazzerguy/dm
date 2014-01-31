package main

import (
   "github.com/dgnorton/dmapi"
   "fmt"
   "log"
   "path"
   "strconv"
)

var cmdRm = &Command {
   UsageLine:  "rm id | last",
   Short:      "removes the specified entry from the local data store",
   Long: `
Removes the entry with the given ID from the local copy of your DailyMile data.
Run 'dm rm last' to remove the most recent entry or use the 'dm find'
command with the '-f entries_id.csv' output format to find the ID of the entry
you want to remove.  This will NOT update dailymile.com.`,
}

func init() {
   cmdRm.Run = runRm
   addRmFlags(cmdRm)
}

// Add command line flags specific to the rm command here
   //var rmDummy string      // -d dummy example

func addRmFlags(cmd *Command) {
   //cmd.Flag.StringVar(&rmDummy, "d", "", "")
}

func runRm(cfg *config, cmd *Command, args []string) {
   if len(args) < 1 {
      usage()
   }

   var err error
   id := -1

   if args[0] != "last" {
      id, err = strconv.Atoi(args[0])
      if err != nil {
         log.Fatalf("%s", err)
      }
   }

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

   err = entries.Remove(id)
   if err != nil {
      log.Fatalf("%s", err)
   }

   err = dmapi.SaveEntries(userFile, entries)
   if err != nil {
      log.Fatalf("%s", err)
   }

   fmt.Println("Entry removed.")
}

