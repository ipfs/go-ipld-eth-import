package lib

import (
	"bytes"
	"context"
	"io/ioutil"
	"math"

	cid "github.com/ipfs/go-cid"

	"github.com/ipfs/go-ipfs/commands/files"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coredag"
	"github.com/ipfs/go-ipfs/pin"
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

	// Also, load the eth-ipld plugins HERE
	coredag.DefaultInputEncParsers.AddParser("raw", "eth-state-trie", ipldeth.EthStateTrieRawInputParser)

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

func (m *IPFS) DagPut(raw []byte, format string) string {
	// Dag Put command options
	ienc := "raw"
	mhType := uint64(math.MaxUint64)

	// Stands for --pin ?
	defer m.n.Blockstore.PinLock().Unlock()

	// Convert the raw bytes into a NopCloser, which
	// in turn will create a file object
	r := ioutil.NopCloser(bytes.NewReader(raw))
	file := files.NewReaderFile("", "", r, nil)

	// Parse your raw data into a DAG Node
	nds, err := coredag.ParseInputs(ienc, format, file, mhType, -1)
	if err != nil {
		panic(err)
	}
	if len(nds) == 0 {
		panic("no nodes returned from parse inputs")
	}

	// Adding the IPLD block
	b := m.n.DAG.Batch()
	_, err = b.Add(nds[0])
	if err != nil {
		panic(err)
	}
	err = b.Commit()
	if err != nil {
		panic(err)
	}

	// Pin it!
	m.n.Pinning.PinWithMode(nds[0].Cid(), pin.Direct)
	err = m.n.Pinning.Flush()
	if err != nil {
		panic(err)
	}

	return nds[0].String()
}
