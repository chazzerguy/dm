package main

import (
   "github.com/dgnorton/dmapi"
   "fmt"
   "log"
   "path"
)

var cmdSync = &Command {
   UsageLine:  "sync [-u username]",
   Short:      "sync local entry DB to dailymile.com",
   Long: `
Downloads all new entries, for the specified user, from dailymile.com.
If no user is specified on the command line, the default user is synced.
Use 'dm user [user name]' to change the default user.`,
}

func init() {
   cmdSync.Run = runIncSync
   addSyncFlags(cmdSync)
}

// Add command line flags specific to the sync command here
   //var syncU string    // -u user name

func addSyncFlags(cmd *Command) {
   //cmd.Flag.StringVar(&syncU, "u", "", "")
}

func runFullSync(cfg *config, cmd *Command, args []string) {
   user := cfg.User
   makeUserDir(user)

   var entries dmapi.Entries
   for pgNbr := 1;; pgNbr = pgNbr + 1 {
      pageEntries, err := dmapi.EntriesByPage(user, pgNbr)
      if err != nil {
         log.Fatalf("%s", err)
      } else if len(pageEntries.Entries) == 0 {
         break
      }
      for _, entry := range pageEntries.Entries {
         entries.Entries = append(entries.Entries, entry)
      }
   }

   usrDir, _ := userDir(user)
   userFile := path.Join(usrDir, "entries.json")
   err := dmapi.SaveEntries(userFile, &entries)
   if err != nil {
      log.Fatalf("%s", err)
   }
}

func runIncSync(cfg *config, cmd *Command, args []string) {
   user := cfg.User
   makeUserDir(user)
   usrDir, _ := userDir(user)
   userFile := path.Join(usrDir, "entries.json")

   haveUserFile, err := isFile(userFile)
   if err != nil {
      log.Fatalf("%s", err)
   }

   if haveUserFile == false {
      fmt.Println("Performing initial (full) sync.  This may take a few minutes.")
      runFullSync(cfg, cmd, args)
      return
   }

   entries, err := dmapi.LoadEntries(userFile)
   if err != nil {
      log.Fatalf("%s", err)
   }

   if len(entries.Entries) == 0 {
      fmt.Println("Performing initial (full) sync.  This may take a few minutes.")
      runFullSync(cfg, cmd, args)
      return
   }

   t, err := entries.Entries[0].Time()
   if err != nil {
      log.Fatalf("%s", err)
   }

   newEntries, err := dmapi.EntriesSince(user, t.Unix())
   if err != nil {
      log.Fatalf("%s", err)
   } else if len (newEntries.Entries) == 0 {
      fmt.Println("Already up-to-date.")
      return
   }

   entries.Entries = append(newEntries.Entries, entries.Entries...)

   usrDir, _ = userDir(user)
   userFile = path.Join(usrDir, "entries.json")
   err = dmapi.SaveEntries(userFile, entries)
   if err != nil {
      log.Fatalf("%s", err)
   }

   if len(newEntries.Entries) == 1 {
      fmt.Println("Synced 1 new entry.")
   } else {
      fmt.Printf("Synced %d new entries.\n", len(newEntries.Entries))
   }
}
