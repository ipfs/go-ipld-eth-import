package lib

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type statusData struct {
	ProtocolVersion uint32
	NetworkId       uint32
	TD              *big.Int
	CurrentBlock    common.Hash
	GenesisBlock    common.Hash
}

func (nm *NetworkManager) doGetStatusMessage(r *rlpx) (*statusData, error) {
	status := &statusData{}

	// Prepare the status message we are going to send
	_td := new(big.Int)
	td, _ := _td.SetString("17179869184", 10)
	ourStatus := &statusData{ProtocolVersion: uint32(63),
		NetworkId:    uint32(1),
		TD:           td,
		CurrentBlock: common.HexToHash("d4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3"),
		GenesisBlock: common.HexToHash("d4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3"),
	}

	err := Send(r.rw, 0x10, ourStatus) // StatusData (0x00) + offset
	if err != nil {
		return status, err
	}

	msg, err := r.ReadMsg()
	if err != nil {
		return status, err
	}

	if msg.Code != 0x10 { // StatusMsg + offset
		return status, fmt.Errorf("StatusMsg error: first msg has code %x (!= %x)", msg.Code, 0x10)
	}
	protocolMaxMsgSize := uint32(10 * 1024 * 1024)
	if msg.Size > protocolMaxMsgSize {
		return status, fmt.Errorf("Message too large error: ", "%v > %v", msg.Size, protocolMaxMsgSize)
	}
	// Decode the return message into a statusData struct
	if err := msg.Decode(status); err != nil {
		return status, fmt.Errorf("Decode error: msg %v: %v", msg, err)
	}

	return status, nil
}
