package main

import (
	"fmt"
	"github.com/alecthomas/kingpin"
	"log"
)

type Command []string

var (
	cmdApp     = kingpin.New("bin-vcs", "A binary version control system.")
	cmdDebug   = cmdApp.Flag("debug", "Print verbose/debug messages.").Bool()
	cmdVerbose = cmdApp.Flag("v", "Be verbose about things!").Bool()

	cmdAdd       = cmdApp.Command("add", "Stage changed files.")
	cmdAddForce  = cmdAdd.Flag("f", "Whether to force ignored files to be added.").Bool()
	cmdAddDry    = cmdAdd.Flag("n", "Dry run (don't actually add the files)").Bool()
	cmdAddUpdate = cmdAdd.Flag("u", "Only update already staged files which have changed, don't add new ones.").Bool()
	cmdAddGlob   = cmdAdd.Arg("glob", "Glob of files to stage.").Required().String()

	cmdLog       = cmdApp.Command("log", "Shows a log of the last commits.")
	cmdLogNumber = cmdLog.Flag("n", "Number of commits to show.").Default("10").Int()

	cmdCheckout     = cmdApp.Command("checkout", "Checks out the given branch.")
	cmdCheckoutNew  = cmdCheckout.Flag("b", "Whether to check out a new branch.").Bool()
	cmdCheckoutName = cmdCheckout.Arg("branch", "Name of branch or hash to checkout.").Required().String()

	cmdStatus        = cmdApp.Command("status", "Prints the status of the repository.")
	cmdStatusBranch  = cmdStatus.Flag("b", "Show all branches.").Default("true").Bool()
	cmdStatusChanged = cmdStatus.Arg("c", "Show changed files").Default("true").Bool()

	cmdVersion = cmdApp.Command("version", "Prints the version number of bin-vcs")
)

func LogVerbose(message string, formats ...interface{}) {
	if *cmdVerbose {
		log.Printf(message+"\n", formats...)
	}
}

func Resolve(parts []string) {
	switch kingpin.MustParse(cmdApp.Parse(parts)) {
	// case cmdStatus.FullCommand():
	// 	return Status()
	// case cmdCheckout.FullCommand():
	// 	return Checkout()
	// case cmdLog.FullCommand():
	// 	return Log()
	case cmdAdd.FullCommand():
		Stage()
	case cmdVersion.FullCommand():
		fmt.Println("bin-vcs version %s", VERSION)
	default:
		log.Panic("Unrecognized command. Run `bin-vcs --help` for usage.")
	}
}
