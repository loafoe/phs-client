package main

import (
	_ "encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/joho/godotenv"
	"github.com/loafoe/cfutil"
	"github.com/mitchellh/cli"
	"os"
)

type AuthRequest struct {
	LoginId  string `json:"loginId"`
	Password string `json:"password"`
}

func init() {
	goEnv := os.Getenv("GOENV")
	if goEnv != "" {
		err := godotenv.Load(goEnv + ".env")
		if err != nil {
			log.Error(err)
		}
	} else {
		godotenv.Load("development.env")
	}
	//cfInit()
}

func main() {
	os.Exit(realMain())
}

func realMain() int {

	args := os.Args[1:]
	// Get the command line args. We shortcut "--version" and "-v" to
	// just show the version.
	for _, arg := range args {
		if arg == "--" {
			break
		}
		if arg == "-v" || arg == "--version" {
			newArgs := make([]string, len(args)+1)
			newArgs[0] = "version"
			copy(newArgs[1:], args)
			args = newArgs
			break
		}
	}

	cli := &cli.CLI{
		Args:     args,
		Commands: Commands,
		HelpFunc: cli.BasicHelpFunc("phsclient"),
	}

	exitCode, err := cli.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		return 1
	}

	return exitCode
}

func cfInit() {
	db, connectString, err := cfutil.NewConnection("postgres", "")
	if err != nil {
		log.Fatal(err)
	}
	log.Print(db.DriverName())
	if cfutil.IsFirstInstance() {
		log.Print("Migrating..")
		cfutil.Migrate(connectString, "./migrations")
	}
}
