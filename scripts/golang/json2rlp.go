package main

/*
  This is a WIP

  Needs a number of iterations
*/

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Defines the output of the JSON RPC API for either
// "eth_BlockByHash" or "eth_BlockByHeader".
type objJSONBlock struct {
	Result objJSONBlockResult `json:"result"`
}

// Nested struct that takes the contents of the JSON field "result".
type objJSONBlockResult struct {
	types.Header           // Use its fields and unmarshaler
	*objJSONBlockResultExt // Add these fields to the parsing
}

// Facilitates the composition of the field "result", adding to the
// Header fields, both uncles and transactions.
type objJSONBlockResultExt struct {
	UncleHashes  []common.Hash        `json:"uncles"`
	Transactions []*types.Transaction `json:"transactions"`
}

// Overrides the function types.Header.UnmarshalJSON, allowing us
// to parse the fields of Header, plus uncles and transactions.
func (o *objJSONBlockResult) UnmarshalJSON(input []byte) error {
	err := o.Header.UnmarshalJSON(input)
	if err != nil {
		return err
	}

	o.objJSONBlockResultExt = &objJSONBlockResultExt{}
	err = json.Unmarshal(input, o.objJSONBlockResultExt)
	if err != nil {
		return err
	}

	return nil
}

func parseJSONBlock(filepath string, obj *objJSONBlock) error {
	fi, err := os.Open(filepath)
	if err != nil {
		return err
	}

	dec := json.NewDecoder(fi)

	err = dec.Decode(obj)
	if err != nil {
		return err
	}

	return nil
}

// Este en el fondo es la proto-version del JSON to RLP que quiero tener
func TestTest(t *testing.T) {

	var blockBody objJSONBlock
	err := parseJSONBlock("test_data/eth-block-body-json-997522", &blockBody)
	if err != nil {
		t.Fatal(err)
	}

	var uncleHeader0 objJSONBlock
	err = parseJSONBlock("test_data/eth-uncle-json-997522-0", &uncleHeader0)
	if err != nil {
		t.Fatal(err)
	}

	var uncleHeader1 objJSONBlock
	err = parseJSONBlock("test_data/eth-uncle-json-997522-1", &uncleHeader1)
	if err != nil {
		t.Fatal(err)
	}

	// We have our elements, let's build the block body to RLP
	body := types.NewBlock(
		&blockBody.Result.Header,
		blockBody.Result.Transactions,
		[]*types.Header{&uncleHeader0.Result.Header, &uncleHeader1.Result.Header},
		[]*types.Receipt{},
	)

	buf := new(bytes.Buffer)
	err = body.EncodeRLP(buf)
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile("test_data/eth-block-body-rlp-997522", buf.Bytes(), 0644)
	if err != nil {
		t.Fatal(err)
	}
}
