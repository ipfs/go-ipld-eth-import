package main

type IPFS struct {
	cmd string
}

func ipfsInit(cmd string) *IPFS {
	return &IPFS{
		cmd: cmd,
	}
}
