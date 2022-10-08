package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	git "github.com/go-git/go-git/v5"
)

type programOptions struct {
	repositoryPath   string
	executionCommand string
	pullInterval     int
	maxRetries       int
}

var usage = func() {
	fmt.Fprintf(os.Stderr, "GitPullExecutor is a lightweight tool that continuously pulls an already cloned git repository and executes a given command whenever a change is pulled. \n Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func parseProgramOptions() (programOptions, error) {
	var parsed programOptions

	flag.StringVar(&parsed.repositoryPath, "repository-path", ".", "Local path of the repository to pull")
	flag.StringVar(&parsed.executionCommand, "execute", "", "Command to execute on successful pull")
	flag.IntVar(&parsed.pullInterval, "pull-interval", 60, "Pull interval")
	flag.IntVar(&parsed.maxRetries, "max-retries", 3, "Maximum successive retries after failing to pull")
	flag.Parse()

	requiredFlags := []string{parsed.executionCommand}

	// Checking if required flags are set, currently only works for string flags
	for _, requiredFlag := range requiredFlags {
		if len(requiredFlag) == 0 {
			usage()
			return programOptions{}, fmt.Errorf("missing required parameter. Please refer to %s --help to see which parameters are not defaulted", os.Args[0])
		}
	}

	return parsed, nil
}

func logFatalIfError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	parsed, err := parseProgramOptions()
	logFatalIfError(err)

	repository, err := git.PlainOpen(parsed.repositoryPath)
	logFatalIfError(fmt.Errorf("failed to open repository: '%s'", err))

	worktree, err := repository.Worktree()
	logFatalIfError(fmt.Errorf("failed to get work tree: '%s'", err))

	cmd := exec.Command(parsed.executionCommand)

	var successiveErrorCount int
	for {
		err = worktree.Pull(&git.PullOptions{RemoteName: "origin"})

		// No changes
		if err == git.NoErrAlreadyUpToDate {
			goto TIMEOUT
		}

		// Actual error during pull
		if err != nil {
			successiveErrorCount++
			if successiveErrorCount > parsed.maxRetries {
				log.Fatalln("Reached pull retry count limit, exiting...")
			}

			fmt.Printf("Unable to pull repository: %s. Retrying...\n", err)
		}

		// Successful pull, execute command
		err = cmd.Run()
		logFatalIfError(fmt.Errorf("failed to run command: '%s'", err))

	TIMEOUT:
		time.Sleep(time.Duration(parsed.pullInterval))
	}

}
