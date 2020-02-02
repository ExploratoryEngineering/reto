package gitutil

import (
	"errors"
	"time"

	"github.com/ExploratoryEngineering/reto/pkg/toolbox"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// HasChanges returns true if there's uncomitted or unstaged changes on the
// current branch.
func HasChanges(rootDir string) bool {
	src, err := git.PlainOpen(rootDir)
	if err != nil {
		toolbox.PrintError("Could not open Git repo at %s: %v", rootDir, err)
		return true
	}
	tree, err := src.Worktree()
	if err != nil {
		toolbox.PrintError("Could not read the working tree for %s: %v", rootDir, err)
		return true
	}
	status, err := tree.Status()
	if err != nil {
		toolbox.PrintError("Could not read status for the working tree at %s: %v", rootDir, err)
		return true
	}
	// The returned values is a map of changes. If map length is 0 there is no
	// staged, unstaged or uncommited files.
	for _, v := range status {
		if v.Staging == git.Untracked {
			// Untracked files are OK
			continue
		}
		// Any other: Not OK
		return true
	}
	return false
}

// GetCurrentHash returns the current hash for HEAD by digging through
// the .git directory. The hash is stored somewhere in .git/refs/heads and
// the file .git/HEAD points to the current branch
func GetCurrentHash(rootDir string) (string, error) {
	src, err := git.PlainOpen(rootDir)
	if err != nil {
		toolbox.PrintError("Could not open Git repo at %s: %v", rootDir, err)
		return "", err
	}
	ref, err := src.Head()
	if err != nil {
		toolbox.PrintError("Could not read the HEAD of %s: %v", rootDir, err)
		return "", err
	}
	if ref.Hash().IsZero() {
		toolbox.PrintError("Could not find the hash for the latest commit at %s", rootDir)
		return "", errors.New("no hash")
	}
	return ref.Hash().String(), nil
}

// TagVersion creates a version tag in Git
func TagVersion(rootDir, name, email, tagName, message string) error {
	src, err := git.PlainOpen(rootDir)
	if err != nil {
		toolbox.PrintError("Could not open Git repo at %s: %v", rootDir, err)
		return err
	}

	ref, err := src.Head()
	if err != nil {
		toolbox.PrintError("Could not read the HEAD of %s: %v", rootDir, err)
		return err
	}
	_, err = src.CreateTag(tagName, ref.Hash(), &git.CreateTagOptions{
		Tagger: &object.Signature{
			Name:  name,
			Email: email,
		},
		Message: message,
	})
	if err != nil {
		toolbox.PrintError("Could not create a tag in %s: %v", rootDir, err)
		return err
	}
	return nil
}

// CreateCommit creates a new commit.
func CreateCommit(rootDir, name, email, message string, files ...string) (string, error) {
	src, err := git.PlainOpen(rootDir)
	if err != nil {
		toolbox.PrintError("Could not open Git repo at %s: %v", rootDir, err)
		return "", err
	}
	tree, err := src.Worktree()
	if err != nil {
		toolbox.PrintError("Could not read the working tree for %s: %v", rootDir, err)
		return "", err
	}
	for _, v := range files {
		_, err := tree.Add(v)
		if err != nil {
			toolbox.PrintError("Could not add %s to the working tree: %v", v, err)
			return "", err
		}
	}
	hash, err := tree.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  name,
			Email: email,
			When:  time.Now(),
		},
		Committer: nil,
	})
	if err != nil {
		toolbox.PrintError("Could not commit to the working tree: %v", err)
		return "", err
	}
	return hash.String(), nil
}
