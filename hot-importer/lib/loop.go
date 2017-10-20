package lib

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/logger/glog"
	gorp "gopkg.in/gorp.v1"
)

// TODO
// Fill this constant from either a file or a console flag
const (
	MAX_PROCESSING_NODES = 500
)

func GetPeerInformation(dbmap *gorp.DbMap) {
	manager := newNetworkManager()
	manager.dbmap = dbmap

	// To be able to talk with other nodes, we need to have our own key
	nodeKey, err := crypto.GenerateKey()
	if err != nil {
		glog.Fatalf("could not generate key: %v", err)
	}
	manager.nodeKey = nodeKey
	manager.updatePeerChan = make(chan *database.Peer)

	// Updating Loop
	go func() {
		for {
			peer := <-manager.updatePeerChan
			if _, err := manager.dbmap.Update(peer); err != nil {
				glog.Infof("Error updating n-tuple in the DB %v", err)
			}
		}
	}()

	// Feeding Loop
	go func() {
		for {
			var peers []*database.Peer
			_, err := dbmap.Select(&peers,
				fmt.Sprintf(
					"SELECT * FROM peers ORDER BY updated_at ASC LIMIT %v",
					MAX_PROCESSING_NODES-manager.getProcessingPeers()))
			if err != nil {
				glog.Infof("Error retrieving peers: %v", err)
			}

			// Add all found peers to server
			for _, peer := range peers {
				manager.addProcessingPeer(peer)
			}

			time.Sleep(5 * time.Second)
		}
	}()
}
