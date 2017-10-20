package lib

import (
	"crypto/ecdsa"
	"sync"
	"time"

	gorp "gopkg.in/gorp.v1"
)

type NetworkManager struct {
	dbmap          *gorp.DbMap
	nodeKey        *ecdsa.PrivateKey
	peersInProcess PeersInProcess
	updatePeerChan chan *database.Peer
}

type PeersInProcess struct {
	mutex sync.RWMutex
	Map   map[string]*database.Peer
	Cnt   int
}

func newNetworkManager() *NetworkManager {
	response := &NetworkManager{}
	response.peersInProcess = PeersInProcess{}
	response.peersInProcess.Map = make(map[string]*database.Peer)
	response.peersInProcess.Cnt = 0

	return response
}

func (nm *NetworkManager) getProcessingPeers() int {
	var response int
	nm.peersInProcess.mutex.Lock()
	response = nm.peersInProcess.Cnt
	nm.peersInProcess.mutex.Unlock()

	return response
}

func (nm *NetworkManager) addProcessingPeer(peer *database.Peer) {
	nm.peersInProcess.mutex.Lock()

	nm.peersInProcess.Map[peer.EthID] = peer
	nm.peersInProcess.Cnt += 1

	// Give the actual order
	go nm.process(peer)

	nm.peersInProcess.mutex.Unlock()
}

func (nm *NetworkManager) removeProcessingPeer(id string) {
	nm.peersInProcess.mutex.Lock()

	delete(nm.peersInProcess.Map, id)
	nm.peersInProcess.Cnt -= 1

	nm.peersInProcess.mutex.Unlock()
}

// UpdatePeerResult is a convenience method
func (nm *NetworkManager) updatePeer(peer *database.Peer) {
	peer.UpdatedAt = time.Now().UnixNano()
	peer.UpdateSource = "central"

	nm.updatePeerChan <- peer
	nm.removeProcessingPeer(peer.EthID)
}
