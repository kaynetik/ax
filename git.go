package ax

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	cmdGit                  = "git"
	cmdGitInit              = "init"
	cmdGitRemoteAddOrigin   = "remote add origin"
	cmdGitAddDot            = "add ."
	cmdGitCommitDashM       = "commit -m"
	cmdGitForcePushToMaster = "push -u origin master --force"

	commitMessageStatic = ":art:Test_Commit_Message"
	zeroInt             = int(0)
)

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
func PushToGIT(gitRepo string, args ...gitChain) error {
	gc := gitChain{}

	if args == nil {
		gc = buildDefaultForcePushGitChain()
	} else {
		gc = args[zeroInt]
	}

	err := gc.gitInit()
	if err != nil {
		return ErrCmdWrapFn(cmdGit, cmdGitInit, err)
	}

	err = gc.gitAddRemote(gitRepo)
	if err != nil {
		return fmt.Errorf("failed adding git remote: %w", err)
	}

	err = gc.gitStageDot()
	if err != nil {
		return fmt.Errorf("failed staging changes: %w", err)
	}

	err = gc.gitCommitM(commitMessageStatic)
	if err != nil {
		return fmt.Errorf("failed committing changes: %w", err)
	}

	err = gc.gitForcePushMaster()
	if err != nil {
		return fmt.Errorf("failed pushing changes: %w", err)
	}

	return nil
}

type (
	gitCmdExec    func() error
	gitCmdStrExec func(string) error
)

type gitChain struct {
	gitInit            gitCmdExec
	gitAddRemote       gitCmdStrExec
	gitStageDot        gitCmdExec
	gitCommitM         gitCmdStrExec
	gitForcePushMaster gitCmdExec
}

// buildDefaultForcePushGitChain - used to generate defaults for gitChain which will be used to force push
// to the initialized Git Repo.
func buildDefaultForcePushGitChain() gitChain {
	return gitChain{
		gitInit:            gitInit,
		gitAddRemote:       gitAddRemote,
		gitStageDot:        gitStageDot,
		gitCommitM:         gitCommitM,
		gitForcePushMaster: gitForcePushMaster,
	}
}

func gitInit() error {
	err := executeCommand(cmdGit, cmdGitInit)
	if err != nil {
		return err
	}

	printStdoutLn("Initialized Repo")

	return nil
}

func gitAddRemote(gitRepo string) error {
	err := executeCommand(cmdGit, fmt.Sprintf("%s %s", cmdGitRemoteAddOrigin, gitRepo))
	if err != nil {
		return err
	}

	printStdoutLn("Added Git Remote")

	return nil
}

func gitStageDot() error {
	err := executeCommand(cmdGit, cmdGitAddDot)
	if err != nil {
		return err
	}

	printStdoutLn("Added all Archives to the Repo")

	return nil
}

// TODO: Make commit message dynamic - based on metadata
// [GH Issue #2](https://github.com/kaynetik/ax/issues/2)
func gitCommitM(commitMsg string) error {
	err := executeCommand(cmdGit, fmt.Sprintf(`%s %s`, cmdGitCommitDashM, commitMsg))
	if err != nil {
		return err
	}

	printStdoutLn("Made a Commit")

	return nil
}

// TODO: Improve this flow
// [GH Issue #1(https://github.com/kaynetik/ax/issues/1)
func gitForcePushMaster() error {
	err := executeCommand(cmdGit, cmdGitForcePushToMaster)
	if err != nil {
		return fmt.Errorf(": %w", err)
	}

	printStdoutLn("Pushed a Commit to origin/master")

	return nil
}

func printStdoutLn(args ...interface{}) {
	_, _ = fmt.Fprintln(os.Stdout, args...)
}
