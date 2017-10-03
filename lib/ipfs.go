package lib

import (
	"context"

	cid "github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
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
