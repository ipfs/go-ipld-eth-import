package network

import (
	"encoding/hex"
	"fmt"
	"net"
	"time"

	"github.com/ethereum/go-ethereum/logger"
	"github.com/ethereum/go-ethereum/logger/glog"
)

// process does most of the heavy lifting here
// We will get a peer and attempt to:
// * TCP dial it
// * Do an encrypted handshake to get a RLPx pipe
// * Perform a protocol handshake to get its name and caps
// * Try to obtain its status data
// * Try to obtain its highest block
func (nm *NetworkManager) process(peer *database.Peer) {
	// 1x. TCP Dial
	// Will give us a net.Conn
	glog.V(logger.Debug).Infof("Attempting to dial %v %v:%v", peer.EthID[:8], peer.IP, peer.TCP)
	dialer := &net.Dialer{Timeout: 5 * time.Second}

	conn, err := dialer.Dial("tcp", fmt.Sprintf("%v:%v", peer.IP, peer.TCP))
	if err != nil {
		glog.V(logger.Debug).Infof("Connection attempt to %v failed", peer.EthID[:8])

		peer.StatusCode = "10"
		peer.StatusMsg = "Dial failed"
		peer.StatusInfo = err.(*net.OpError).Err.Error()
		nm.updatePeer(peer)
		return
	}
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer conn.Close()

	// 2x. Encryption Handshake
	// This is very similar to the idea of an SSL handshake.
	// We want to obtain a secrets struct filled for future use.
	glog.V(logger.Debug).Infof("Attempting to get an encryption handshake with %v", peer.EthID[:8])

	secrets, err := nm.doEncryptionHandshake(conn, peer)
	if err != nil {
		if _err, ok := err.(*net.OpError); ok {
			err = _err
		}
		glog.V(logger.Debug).Infof("Enc Handshake attempt to %v failed", peer.EthID[:8])

		peer.StatusCode = "20"
		peer.StatusMsg = "Enc Handshake failed"
		peer.StatusInfo = err.Error()
		nm.updatePeer(peer)
		return
	}

	// This object facilitates communication on each following stage
	r := &rlpx{
		fd: conn,
		rw: newRLPXFrameRW(conn, secrets),
	}

	// 3x. Protocol Handshake
	// Our node sends an encrypted message to the peer
	// telling its id, name, caps (protocols supported).
	// A similar package is expected from the peer.
	// Name and Caps of the peer will be stored in the DB.
	glog.V(logger.Debug).Infof("Attempting to get a protocol handshake with %v", peer.EthID[:8])

	ph, err := nm.doProtocolHandshake(r)
	if err != nil {
		glog.V(logger.Debug).Infof("Proto Handshake attempt to %v failed", peer.EthID[:8])

		peer.StatusCode = "30"
		peer.StatusMsg = "Proto Handshake failed"
		peer.StatusInfo = err.Error()
		nm.updatePeer(peer)
		return
	}
	// It worked, let's get the peer's info
	peer.Name = ph.Name
	for i, cap := range ph.Caps {
		peer.Caps += cap.String()
		if i != len(ph.Caps)-1 {
			peer.Caps += " "
		}
	}

	// 4x. Get Status Message
	// We will send a status (using the eth63 protocol)
	// and wait for the peer to send us theirs.
	glog.V(logger.Debug).Infof("Attempting to get status of peer %v", peer.EthID[:8])

	status, err := nm.doGetStatusMessage(r)
	if err != nil {
		glog.V(logger.Debug).Infof("Get Status attempt to %v failed", peer.EthID[:8])

		peer.StatusCode = "40"
		peer.StatusMsg = "Get Status failed"
		peer.StatusInfo = err.Error()
		nm.updatePeer(peer)
		return
	}
	// Worked, let's put that info to good use
	peer.Protocol = fmt.Sprint(status.ProtocolVersion)
	peer.Network = fmt.Sprint(status.NetworkId)
	peer.TD = status.TD.String()
	peer.Genesis = hex.EncodeToString(status.CurrentBlock.Bytes())
	peer.Block = hex.EncodeToString(status.GenesisBlock.Bytes())

	// TODO
	// Get height of the peer (5x)

	// DEBUG
	// Temporal tag and code
	peer.StatusCode = "49"
	peer.StatusMsg = "Get Status succeeded"
	peer.StatusInfo = ""
	// DEBUG

	// Let's update and remove the peer from the list
	nm.updatePeer(peer)
	return
}
