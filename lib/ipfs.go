package lib

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"math"

	cid "github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs/commands/files"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coredag"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	ipldeth "github.com/ipfs/go-ipld-eth/plugin"
	node "github.com/ipfs/go-ipld-format"
	mh "github.com/multiformats/go-multihash"
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

	coredag.DefaultInputEncParsers.AddParser("raw", "eth-state-trie", ipldeth.EthStateTrieRawInputParser)
	coredag.DefaultInputEncParsers.AddParser("raw", "importer-ipld-raw-data", ipldRawNodeInputParser)

	return &IPFS{n: ipfsNode, ctx: ctx}
}

// ipldRawNodeInputParser custom input parser to be able to introduce a 0x55 = keccak256 IPLD BLock
func ipldRawNodeInputParser(r io.Reader, mhtype uint64, mhLen int) ([]node.Node, error) {
	rawdata, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	c, err := cid.Prefix{
		Codec:    0x55,
		Version:  1,
		MhType:   mh.KECCAK_256,
		MhLength: -1,
	}.Sum(rawdata)
	if err != nil {
		panic(err)
	}

	rawNode := &IpldRawNode{
		cid:     c,
		rawdata: rawdata,
	}

	return []node.Node{rawNode}, nil
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

	return nds[0].String()
}
