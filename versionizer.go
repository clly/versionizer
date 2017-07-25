package versionizer

import (
	"fmt"

	"bytes"

	gogit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
)

type GitInfo struct {
	Hash string
	Log  []byte
}

func GetGit(repopath string, commits int) (*GitInfo, error) {
	r, err := gogit.PlainOpen(repopath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open git repo at %s %s", repopath, err)
	}

	ref, err := r.Head()
	if err != nil {
		return nil, fmt.Errorf("Failed to get head of git repo at %s %s", repopath, err)
	}

	commitIter, err := r.Log(&gogit.LogOptions{From: ref.Hash()})

	// Setup iterator
	i := 0
	buf := &bytes.Buffer{}
	err = commitIter.ForEach(func(c *object.Commit) error {
		if i >= commits {
			return storer.ErrStop
		}
		buf.WriteString(fmt.Sprintf("\n%s", c.String()))

		return nil
	})
	commitIter.Close()
	if err != nil {
		return nil, fmt.Errorf("Failed to iterate over commits %s", err)
	}

	g := &GitInfo{
		Hash: ref.String(),
		Log:  buf.Bytes(),
	}
	return g, nil
}
