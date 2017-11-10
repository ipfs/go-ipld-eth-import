package devp2p

import (
	"fmt"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
)

// eth protocol message codes
const (
	// Protocol messages belonging to eth/62
	StatusMsg          = 0x00
	NewBlockHashesMsg  = 0x01
	TxMsg              = 0x02
	GetBlockHeadersMsg = 0x03
	BlockHeadersMsg    = 0x04
	GetBlockBodiesMsg  = 0x05
	BlockBodiesMsg     = 0x06
	NewBlockMsg        = 0x07

	// Protocol messages belonging to eth/63
	GetNodeDataMsg = 0x0d
	NodeDataMsg    = 0x0e
	GetReceiptsMsg = 0x0f
	ReceiptsMsg    = 0x10
)

//
func getEth63CompatibleSubProtocol() p2p.Protocol {
	return p2p.Protocol{
		Name:    "eth",
		Version: 63,
		Length:  17,
		Run:     protocolHandler,
	}
}

//
func protocolHandler(p *p2p.Peer, rw p2p.MsgReadWriter) error {
	// This peer is formatted as an eth peer
	ethPeer := &Peer{
		id: p.String(),
		rw: rw,
	}

	log.Trace("protocolHandler", "stage", "start eth protocol handshake", "peer", p)
	if err := ethPeer.sendStatusMsg(); err != nil {
		log.Error("protocolHandler", "stage", "failed eth protocol handshake", "peer", p, "error", err)
		return err
	}

	log.Trace("protocolHandler", "stage", "adding peer to store", "peer", p)

	// TODO
	// Add the peer to our registries
	// Needed to identify the "best" peers

	// TODO
	// Defer the removal of this peer of our registries.
	// This marks the lifecycle of the peer in protocol terms.

	// TODO
	// Handle all incoming messages in a loop.
	// At the first error, we exit the loop, cutting the connection,
	// And removing the peer from our registries

	// PLACEHOLDER
	return nil
	// PLACEHOLDER
}

//
func handleIncomingMsg(rw p2p.MsgReadWriter) error {
	//
	msg, err := rw.ReadMsg()
	if err != nil {
		return err
	}
	defer msg.Discard()

	//
	//
	switch {
	case msg.Code == StatusMsg:
		log.Trace("handleIncomingMsg", "msg", "StatusMsg")
		return fmt.Errorf("not Implemented")

	case msg.Code == NewBlockHashesMsg:
		log.Trace("handleIncomingMsg", "msg", "NewBlockHashesMsg")
		return fmt.Errorf("not Implemented")

	case msg.Code == TxMsg:
		log.Trace("handleIncomingMsg", "msg", "TxMsg")
		return fmt.Errorf("not Implemented")

	case msg.Code == GetBlockHeadersMsg:
		log.Trace("handleIncomingMsg", "msg", "GetBlockHeadersMsg")
		return fmt.Errorf("not Implemented")

	case msg.Code == BlockHeadersMsg:
		log.Trace("handleIncomingMsg", "msg", "BlockHeadersMsg")
		return fmt.Errorf("not Implemented")

	case msg.Code == GetBlockBodiesMsg:
		log.Trace("handleIncomingMsg", "msg", "GetBlockBodiesMsg")
		return fmt.Errorf("not Implemented")

	case msg.Code == BlockBodiesMsg:
		log.Trace("handleIncomingMsg", "msg", "BlockBodiesMsg")
		return fmt.Errorf("not Implemented")

	case msg.Code == NewBlockMsg:
		log.Trace("handleIncomingMsg", "msg", "NewBlockMsg")
		return fmt.Errorf("not Implemented")

	case msg.Code == GetNodeDataMsg:
		log.Trace("handleIncomingMsg", "msg", "GetNodeDataMsg")
		return fmt.Errorf("not Implemented")

	case msg.Code == NodeDataMsg:
		log.Trace("handleIncomingMsg", "msg", "NodeDataMsg")
		return fmt.Errorf("not Implemented")

	case msg.Code == GetReceiptsMsg:
		log.Trace("handleIncomingMsg", "msg", "GetReceiptsMsg")
		return fmt.Errorf("not Implemented")

	case msg.Code == ReceiptsMsg:
		log.Trace("handleIncomingMsg", "msg", "ReceiptsMsg")
		return fmt.Errorf("not Implemented")

	default:
		return fmt.Errorf("message code not supported")
	}
}
