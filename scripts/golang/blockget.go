package main

import (
	"context"
	"fmt"
	"time"

	cid "gx/ipfs/QmNp85zy9RLrQ5oQD4hPyS39ezrrXpcaa7R4Y9kxdWQLLQ/go-cid"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
)

// Yeah, yeah... MyStruct...
type MyStruct struct {
	n   *core.IpfsNode
	ctx context.Context
}

// Trying "dag get" by loading an ipfs node at init time
func main() {
	m := InitIPFS()

	m.BlockGet("z45oqTRuZDa8Kvo9MWtEmGZJfnsg39ngDa4CyHY77oRD5Yd8PiX")
	m.BlockGet("z45oqTS26iKhYHbe5shTQhyW6DZE1Ffy24ENRVKQUvvZZwAyFaV")
}

func InitIPFS() *MyStruct {
	r, err := fsrepo.Open("~/.ipfs")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	cfg := &core.BuildCfg{
		Online: false,
		Repo:   r,
	}

	ipfsNode, err := core.NewNode(ctx, cfg)
	if err != nil {
		panic(err)
	}

	return &MyStruct{n: ipfsNode, ctx: ctx}
}

func (m *MyStruct) BlockGet(cidString string) {
	start := time.Now().UnixNano()

	c, err := cid.Decode(cidString)
	if err != nil {
		panic(err)
	}

	b, err := m.n.Blocks.GetBlock(m.ctx, c)
	if err != nil {
		if err.Error() != "blockservice: key not found" {
			panic(err)
		}
	}

	fmt.Printf("Time Diff: %.0d\n", time.Now().UnixNano()-start)

	// Some output
	if b != nil {
		fmt.Printf("%d bytes\n", len(b.RawData()))
	} else {
		fmt.Printf("not found\n")
	}
}
