package main

import (
   "fmt"
)

var cmdUser = &Command {
   UsageLine:  "user [user name]",
   Short:      "get or set default user",
   Long: `
User gets or sets the default user name.  If the optional [user name] parameter
is given, the default user is changed to this name.  If the parameter is not
provided, the default user name is printed.`,
}

func init() {
   cmdUser.Run = runUser
   addUserFlags(cmdUser)
}

// Add command line flags specific to the user command here
//   var userD bool    // -d dummy example

func addUserFlags(cmd *Command) {
   //cmd.Flag.BoolVar(&userD, "d", false, "")
   // etc.
}

func runUser(cfg *config, cmd *Command, args []string) {
   if len(args) == 0 {
      fmt.Println(cfg.User)
   } else {
      oldUser := cfg.User
      cfg.User = args[0]
      err := saveConfig(cfg)
      if err != nil {
         cfg.User = oldUser
         errorf("%s [user.go]", err)
      } else {
         fmt.Printf("Default user changed to: %s\n", cfg.User)
      }
   }
}
