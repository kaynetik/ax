package ax

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// TODO: Add all static flags as const.
const cmdGit = "git"

// ErrCmdWrapFn - func used to wrap errors caused by executeCommand, if any.
var ErrCmdWrapFn = func(cmd, cmdArgs string, err error) error {
	cmdSlc := strings.Fields(cmdArgs)

	return fmt.Errorf("%s: failed executing the command with args[%v]: %w", cmd, cmdSlc, err)
}

func executeCommand(cmd, cmdArgs string) error {
	cmdSliceInit := strings.Fields(cmdArgs)

	_, err := exec.Command(cmd, cmdSliceInit...).Output()
	if err != nil {
		return fmt.Errorf("%s: failed executing the command with args[%v]: %w", cmd, cmdSliceInit, err)
	}

	return nil
}

// PushToGIT - used to commit&push created archive(s) to the remote GIT Repository.
func PushToGIT(gitRepo string) error {
	// push to git repo
	// init git
	err := executeCommand(cmdGit, "init")
	if err != nil {
		return ErrCmdWrapFn(cmdGit, "init", err)
	}

	log.Println("initialized repo")

	// add remote
	err = executeCommand(cmdGit, fmt.Sprintf("remote add origin %s", gitRepo))
	if err != nil {
		return fmt.Errorf(": %w", err)
	}

	fmt.Println("added git remote")

	// commit
	err = executeCommand(cmdGit, "add .")
	if err != nil {
		return fmt.Errorf(": %w", err)
	}

	fmt.Println("added all archives to the repo")

	err = executeCommand(cmdGit, `commit -m "test_commit_msg"`)
	if err != nil {
		return fmt.Errorf(": %w", err)
	}

	fmt.Println("made a commit")

	// push
	err = executeCommand(cmdGit, "push -u origin master --force")
	if err != nil {
		return fmt.Errorf(": %w", err)
	}

	fmt.Println("pushed commit to remote")

	return nil
}
