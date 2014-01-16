package main

import (
   "flag"
   "fmt"
   "encoding/json"
   "io/ioutil"
   "log"
   "os"
   "os/user"
   "path"
   "strings"
   "sync"
)

type config struct {
   User string
}

var defaultConfig = &config{
   "",   // User
}

func loadConfig() (*config, error) {
   cfgfile, _ := configFile()
   bytes, err := ioutil.ReadFile(cfgfile)
   if err != nil {
      return nil, err
   }
   var cfg config
   err = json.Unmarshal(bytes, &cfg)
   return &cfg, err
}

func saveConfig(c *config) error {
   cfgfile, _ := configFile()
   bytes, err := json.MarshalIndent(c, "", "   ")
   if err != nil {
      return err
   }
   err = ioutil.WriteFile(cfgfile, bytes, 0600)
   return err
}

// A Command is an implementation of a dm command.
// This is borrowed largely from golang's cmd/go/main.go.
type Command struct {
   Run func(cfg *config, cmd *Command, args []string)
   UsageLine string
   Short string
   Long string
   Flag flag.FlagSet
}

func (c *Command) Name() string {
   name := c.UsageLine
   i := strings.Index(name, " ")
   if i >= 0 {
      name = name[:i]
   }
   return name
}

func (c *Command) Usage() {
   fmt.Fprintf(os.Stderr, "usage: %s\n\n", c.UsageLine)
   fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(c.Long))
   os.Exit(2)
}

func (c *Command) Runnable() bool {
   return c.Run != nil
}

// Available commands
var commands = []*Command{
   cmdUser,
   cmdSync,
   cmdFind,
}

var exitStatus = 0
var exitMu sync.Mutex

func setExitStatus(n int) {
   exitMu.Lock()
   if exitStatus < n {
      exitStatus = n
   }
   exitMu.Unlock()
}

func main() {
   var mainUser string  // -u user name

   flag.StringVar(&mainUser, "u", "", "")
   flag.Parse()
   args := flag.Args()

   if len(args) < 1 {
      usage()
   }

   // make working directory if it doesn't exist...or die trying
   makeWorkDir()

   cfg, err := loadConfig()
   if err != nil {
      if os.IsNotExist(err) {
         // create a default config file
         err = saveConfig(defaultConfig)
         if err != nil {
            fatalf("%s [dm.go - saveConfig]", err)
         }
         cfg = defaultConfig
      } else {
         fatalf("%s [dm.go - loadConfig]", err)
      }
   }

   // command line flags override config file
   if mainUser != "" {
      cfg.User = mainUser
   }

   if cfg.User == "" {
      fatalf(
`No user set.  Either use the 'dm user <user name>' command or
the '-u <user name>' command line argument.`)
   }

   for _, cmd := range commands {
      if cmd.Name() == args[0] && cmd.Run != nil {
         cmd.Flag.Usage = func() { cmd.Usage() }
         cmd.Flag.Parse(args[1:])
         args = cmd.Flag.Args()
         cmd.Run(cfg, cmd, args)
         exit()
         return
      }
   }

   fmt.Fprintf(os.Stderr, "dm: unknown subcommand %q\n", args[0])
   setExitStatus(2)
   exit()
}

var atexitFuncs []func()

func atexit(f func()) {
   atexitFuncs = append(atexitFuncs, f)
}

func exit() {
   for _, f := range atexitFuncs {
      f()
   }
   os.Exit(exitStatus)
}

func userDir(user string) (string, error) {
   wrk, err := workDir()
   if err != nil {
      return "", err
   }
   return path.Join(wrk, user), nil
}

func makeUserDir(user string) {
   usr, err := userDir(user)
   if err != nil {
      fatalf("%s [dm.go]", err)
   }
   exists, err := isDir(usr)
   if err != nil {
      fatalf("%s [dm.go]", err)
   }
   if exists == false {
      err = os.MkdirAll(usr, 0700)
      if err != nil {
         fatalf("%s [dm.go]", err)
      }
   }
}

func homeDir() (string, error) {
   usr, err := user.Current()
   if err != nil {
      return "", err
   }
   return  usr.HomeDir, nil
}

func workDir() (string, error) {
   home, err := homeDir()
   if err != nil {
      return "", err
   }
   return path.Join(home, ".dailymile_cli"), nil
}

func makeWorkDir() {
   wrkdir, err := workDir()
   if err != nil {
      fatalf("%s [dm.go]", err)
   }

   exists, err := isDir(wrkdir)
   if err != nil {
      fatalf("%s [dm.go]", err)
   }

   if exists == false {
      err = os.MkdirAll(wrkdir, 0700)
      if err != nil {
         fatalf("%s [dm.go]", err)
      }
   }
}

func configFile() (string, error) {
   wrkdir, err := workDir()
   if err != nil {
      return "", err
   }
   cfgfile := path.Join(wrkdir, "config")
   return cfgfile, nil
}

func isDir(path string) (bool, error) {
   fi, err := os.Stat(path)
   if err == nil {
      return fi.IsDir(), nil
   } else if os.IsNotExist(err) {
      return false, nil
   }
   return false, err
}

func isFile(path string) (bool, error) {
   fi, err := os.Stat(path)
   if err == nil {
      return !fi.IsDir(), nil
   } else if os.IsNotExist(err) {
      return false, nil
   }
   return false, err
}

func fatalf(format string, args ...interface{}) {
   errorf(format, args...)
   exit()
}

func errorf(format string, args ...interface{}) {
   log.Printf(format, args...)
   setExitStatus(1)
}

var logf = log.Printf

func usage() {
   fmt.Fprintf(os.Stderr, "dm is a command line tool for Dailymile.com\n\n  usage: dm command [args]\n")
   os.Exit(2)
}
