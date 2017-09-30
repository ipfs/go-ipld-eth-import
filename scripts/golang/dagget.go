package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"time"
)

func main() {
	DagGet("z45oqTRuZDa8Kvo9MWtEmGZJfnsg39ngDa4CyHY77oRD5Yd8PiX")
	DagGet("z45oqTS26iKhYHbe5shTQhyW6DZE1Ffy24ENRVKQUvvZZwAyFaV")
}

func DagGet(cidString string) {

	cmd := exec.Command("ipfs", "dag", "get", "--local", cidString)

	var out bytes.Buffer
	cmd.Stdout = &out

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Here!
	start := time.Now().UnixNano()
	err := cmd.Run()
	fmt.Printf("Time Diff: %.0d\n", time.Now().UnixNano()-start)

	if err != nil {
		fmt.Printf("%s\n", stderr.String())
	} else {
		fmt.Printf("%s\n", out.String())
	}
}
