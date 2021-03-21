// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"fmt"
	git "github.com/amzn/golang-gitoperations/gitoperations"
	"os"
	"os/exec"
)

const (
	green   = "\033[1;32m"
	red     = "\033[1;31m"
	noColor = "\033[0m"
	usage   = `amigitenough

amigitenough will query git for the values of some recommended settings, and provide feedback.
More details: https://github.com/amzn/amigitenough`
)

func colorText(color, text string) string {
	return color + text + noColor
}

func printlnGreen(text string) { fmt.Println(colorText(green, text)) }
func printlnRed(text string) { fmt.Println(colorText(red, text)) }

var verboseFlag bool
const (
	TooManyArgs = 1
	CantFindGit = 2
	CantExecuteGit = 4
	RerereNotEnabled = 8
	PullRebaseNotEnabled = 16
	MissingIdentity = 32
	)

func myUsage() {
	fmt.Fprintf(flag.CommandLine.Output(), "%s\n", usage)
	flag.PrintDefaults()
}

func main() {
	flag.Usage = myUsage
	flag.BoolVar(&verboseFlag, "verbose", false, "Print the git commands used during execution.")
	flag.Parse()
	if verboseFlag {
		git.SetTrace(true)
	}
	if flag.NArg() > 0 {
		printlnRed("Error: too many arguments.")
		flag.Usage()
		os.Exit(TooManyArgs)
	}

	if _, err := exec.LookPath("git"); err != nil {
		fmt.Printf("Unable to find git executable: %q\n",err.Error())
		os.Exit(CantFindGit)
	}

	if err := git.GitCanExecute(exec.Command); err != nil {
		fmt.Printf("git fails to execute with error: %q\n", err.Error())
		os.Exit(CantExecuteGit)
	}

	exitCode := 0

	if setting, _ := git.GetGlobalConfigSetting(exec.Command, "pull.rebase"); setting != "true" {
		exitCode |= PullRebaseNotEnabled
		fmt.Println("Rebase when performing \"git pull\" is not enabled. Correct this by running: git config --global pull.rebase true")
	}
	if setting, _ := git.GetGlobalConfigSetting(exec.Command, "rerere.enabled"); setting != "true" {
		exitCode |= RerereNotEnabled
		fmt.Println("Rerere when performing rebase is not enabled. Correct this by running: git config --global rerere.enabled true")
	}

	// We track the more boring defects of name and email address with exitCode 32.
	// We discovered a case of a user who uses includeIfs to set the idenitity fields, and decided that
	// amigitenough will not demand a global setting for those settings.  Hence we use GetConfigSetting instead of
	// GetGlobalConfigSetting.
	if _, err := git.GetConfigSetting(exec.Command, "user.name"); err != nil {
		exitCode |= MissingIdentity
		fmt.Println("User's name is not defined. Correct this by running: git config --global user.name \"Your Name\"")
	}
	if _, err := git.GetConfigSetting(exec.Command, "user.email"); err != nil {
		exitCode |= MissingIdentity
		fmt.Println("User's Email is not defined. Correct this by running: git config --global user.email \"youremail@yourdomain.com\"")
	}

	if exitCode == 0 {
		printlnGreen("VALIDATION SUCCEEDED")
	} else {
		printlnRed("VALIDATION FAILED")
		fmt.Println("Learn more at: https://github.com/amzn/amigitenough")
	}
	os.Exit(exitCode)
}
