package lib

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"

	"github.com/ethereum/go-ethereum/rlp"
)

// protoHandshake is the RLP structure of the protocol handshake.
type protoHandshake struct {
	Version    uint64
	Name       string
	Caps       []Cap
	ListenPort uint64
	ID         []byte
}

// Cap is the structure of a peer capability.
type Cap struct {
	Name    string
	Version uint
}

func (nm *NetworkManager) doProtocolHandshake(r *rlpx) (*protoHandshake, error) {
	ph := &protoHandshake{
		Version: 4, // baseProtocolVersion. Why?
		Name:    "フクロウ central",
		ID:      PubkeyID(&nm.nodeKey.PublicKey),
		Caps: []Cap{
			Cap{
				Name:    "eth",
				Version: uint(63),
			},
			Cap{
				Name:    "eth",
				Version: uint(62),
			},
		},
	}

	err := Send(r, 0x00, ph)
	if err != nil {
		return nil, fmt.Errorf("error protoHandshake send: %v", err)
	}

	msg, err := r.ReadMsg()
	if err != nil {
		return nil, fmt.Errorf("error protoHandshake read: %v", err)
	}
	if msg.Size > baseProtocolMaxMsgSize {
		return nil, fmt.Errorf("message too big")
	}
	if msg.Code == 0x01 { // DiscMsg
		// Disconnect before protocol handshake is valid according to the
		// spec and we send it ourself if the posthanshake checks fail.
		// We can't return the reason directly, though, because it is echoed
		// back otherwise. Wrap it in a string instead.
		var reason [1]DiscReason
		rlp.Decode(msg.Payload, &reason)
		return nil, reason[0]
	}
	if msg.Code != 0x00 { // handshakeMsg
		return nil, fmt.Errorf("expected handshake, got %x", msg.Code)
	}

	var hs protoHandshake
	if err := msg.Decode(&hs); err != nil {
		return nil, err
	}

	return &hs, nil
}

func (cap Cap) String() string {
	return fmt.Sprintf("%s/%d", cap.Name, cap.Version)
}

func (cap Cap) RlpData() interface{} {
	return []interface{}{cap.Name, cap.Version}
}

// PubkeyID returns a marshaled representation of the given public key.
func PubkeyID(pub *ecdsa.PublicKey) []byte {
	id := make([]byte, 64)
	pbytes := elliptic.Marshal(pub.Curve, pub.X, pub.Y)
	if len(pbytes)-1 != 64 {
		panic(fmt.Errorf("need %d bit pubkey, got %d bits", (64+1)*8, len(pbytes)))
	}
	copy(id[:], pbytes[1:])
	return id
}
