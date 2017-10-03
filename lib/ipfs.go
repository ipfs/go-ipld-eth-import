package lib

import (
	"context"
	"fmt"
	"math"
	"os"

	cid "github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs/commands/files"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coredag"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	ipldeth "github.com/ipfs/go-ipld-eth/plugin"
)

type IPFS struct {
	n   *core.IpfsNode
	ctx context.Context
}

func IpfsInit(repoPath string) *IPFS {
	r, err := fsrepo.Open(repoPath)
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

	// Also, load the eth-ipld plugin HERE
	coredag.DefaultInputEncParsers.AddParser("raw", "eth-block", ipldeth.EthBlockRawInputParser)

	return &IPFS{n: ipfsNode, ctx: ctx}
}

func (m *IPFS) HasBlock(cidString string) bool {
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

	if b != nil {
		return true
	}
	return false
}

func (m *IPFS) DagPut() string {
	// Dag Put command options
	// TODO
	// We may want to parametrize those ones
	ienc := "raw"
	format := "eth-state-trie"
	mhType := uint64(math.MaxUint64)
	// Stands for --pin ?
	defer m.n.Blockstore.PinLock().Unlock()

	// TODO
	// Oh! I need to solve this one somehow!
	// Basically, convert from []byte to a []io.ReadCloser()
	file := files.NewReaderFile("", "", os.Stdin, nil)

	//
	nodes, err := coredag.ParseInputs(ienc, format, file, mhType, -1)
	if err != nil {
		panic(err)
	}

	// TMP
	// TODO
	// What do we return here?
	fmt.Printf("%v\n", nds)
	return "NOT IMPLEMENTED"
	// TMP
}
