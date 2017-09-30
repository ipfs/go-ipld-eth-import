package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

type IPFS struct {
	cmd string
}

func ipfsInit(cmd string) *IPFS {
	return &IPFS{
		cmd: cmd,
	}
}

func (ipfs *IPFS) DagGet(cidString string) (string, error) {
	cmd := exec.Command(ipfs.cmd, "dag", "get", "--local", cidString)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("%s\n", stderr.String())
		return "", err
	} else {
		return stdout.String(), nil
	}
}
