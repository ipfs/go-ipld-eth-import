package devp2p

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/p2p"
)

// PeerStore keeps track of the ethereum peers after a succesful
// handshake (i.e. a match in protocols and version).
// It is also, able to tell the best nodes based on their difficulties (td).
type PeerStore struct {
	peers map[string]Peer

	// TODO
	// A sorted index of ids for difficulties
}

// Peer is an arrangement we use to organize the life cycle of a
// devp2p peer after it is dialed, and the encryption and protocol
// handshakes are performed.
// After a succesuful ethereum handshake, this Peer can be added
// in the peerstore, to be managed for further requesting.
type Peer struct {
	// the id of the devp2p node, shorted to 8 chars.
	id string

	// the communication pipeline
	rw p2p.MsgReadWriter

	// total difficulty informed by the peer in the eth handshake
	td *big.Int

	// current block informed by the peer in the eth handshake
	currentBlock common.Hash
}

// TODO
// new peerstore

// TODO
// add peer
// needs to be asynchronous

// TODO
// remove peer
// needs to be asynchronous

// TODO
// returns the X best peers
// HINT: Those are the ones with the bigger difficulty
// Use the to-be-built index of ids for difficulties
