package main

import (
	"fmt"

	docopt "github.com/docopt/docopt-go"
	tufnotary "github.com/notaryproject/tuf/tuf-notary/tuf-notary"
)

func init() {
	register("init", cmdInit, `
usage: tuf-notary init <registry> [--repo=<repository>]

Initialize a new repository and push it to the TUF repository on the
registry

Options:
  --repo	Set the tuf repository name. By default this will be 'tuf-repo'
  `)
}

func cmdInit(args []string, opts docopt.Opts) error {
	repository := "tuf-repo"
	if r := opts["--repo"]; r != nil {
		repository = r.(string)
	}

	registry := args[0]

	err := tufnotary.Init(repository)

	if err != nil {
		return err
	}

	//upload root with no references
	root_desc, err := tufnotary.UploadTUFMetadata(registry, repository, "root", "")
	if err != nil {
		return err
	}
	fmt.Println("uploaded root " + root_desc.Digest.String())

	//upload targets with a reference to root metadata
	targets_desc, err := tufnotary.UploadTUFMetadata(registry, repository, "targets", "root")
	if err != nil {
		return err
	}
	fmt.Println("uploaded targets " + targets_desc.Digest.String())

	return err
}
