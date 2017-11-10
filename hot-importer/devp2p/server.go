package devp2p

import (
	"crypto/ecdsa"
	"net"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
)

//
func newServer(mgrConfig Config) *p2p.Server {
	dialer := p2p.TCPDialer{&net.Dialer{Timeout: 60 * time.Second}}

	name := getClientName()

	privateKey := getPrivateKey(mgrConfig.PrivateKeyFilePath)

	protocols := []p2p.Protocol{getEth63CompatibleSubProtocol()}

	serverConfig := p2p.Config{
		BootstrapNodes:  mgrConfig.bootnodes,
		Dialer:          dialer,
		ListenAddr:      ":30303",
		MaxPeers:        1000000,
		MaxPendingPeers: 1000000,
		Name:            name,
		NoDiscovery:     false,
		NodeDatabase:    mgrConfig.NodeDatabasePath,
		PrivateKey:      privateKey,
		Protocols:       protocols,
	}
	server := &p2p.Server{Config: serverConfig}
	log.Debug("new devp2p server configured", "instance:", serverConfig.Name)

	return server
}

//
//
// random one if not file
func getPrivateKey(filePath string) *ecdsa.PrivateKey {
	// TODO
	// PLACEHOLDER
	privKeyString := "af8bf8bb4c634b8716880aa44e82da72b902144940a56e1fa787505aa513ba46"
	privateKey, _ := crypto.HexToECDSA(privKeyString)
	// PLACEHOLDER

	// TODO
	// If error, zero tolerance (you don't wanna delete that key file)
	// If empty, just randomize one

	return privateKey
}

//
// Use this name, v 0.0.1 and this commit
func getClientName() string {
	// TODO
	// PLACEHOLDER
	return "Spanish Flea"
	// PLACEHOLDER
}
