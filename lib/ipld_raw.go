package lib

import (
	"fmt"

	cid "github.com/ipfs/go-cid"
	node "github.com/ipfs/go-ipld-format"
)

type IpldRawNode struct {
	cid     *cid.Cid
	rawdata []byte
}

// Static check
var _ node.Node = (*IpldRawNode)(nil)

/*
  Block INTERFACE
*/

// RawData returns the binary of the RLP encode of the block header.
func (i *IpldRawNode) RawData() []byte {
	return i.rawdata
}

// Cid returns the cid of the block header.
func (i *IpldRawNode) Cid() *cid.Cid {
	return i.cid
}

// String is a helper for output
func (i *IpldRawNode) String() string {
	return fmt.Sprintf("<IpldRawNode %s>", i.cid)
}

// Loggable returns a map the type of IPLD Link.
func (i *IpldRawNode) Loggable() map[string]interface{} {
	return map[string]interface{}{
		"type": "ipld-raw-node",
	}
}

/*
  Node INTERFACE
*/

func (i *IpldRawNode) Resolve(p []string) (interface{}, []string, error) {
	return nil, nil, nil
}

func (i *IpldRawNode) Tree(p string, depth int) []string {
	return nil
}

func (i *IpldRawNode) ResolveLink(p []string) (*node.Link, []string, error) {
	return nil, nil, nil
}

func (i *IpldRawNode) Copy() node.Node {
	return nil
}

func (i *IpldRawNode) Links() []*node.Link {
	return nil
}

func (i *IpldRawNode) Stat() (*node.NodeStat, error) {
	return nil, nil
}

func (i *IpldRawNode) Size() (uint64, error) {
	return 0, nil
}
