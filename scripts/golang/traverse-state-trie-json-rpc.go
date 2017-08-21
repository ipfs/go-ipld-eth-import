package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	goque "github.com/beeker1121/goque"
	rlp "github.com/ethereum/go-ethereum/rlp"
	cid "github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
)

const MEthStateTrie = 0x96

// trieStack wraps the goque stack, enabling the adding of specific
// methods for dealing with the state trie.
type trieStack struct {
	*goque.Stack
}

func main() {
	var (
		stateRoot string
		err       error
	)

	// Take command parameters
	// ethapiHost := flag.String("ethapi", "http://localhost:15001", "ethereum json api endpoint")
	flag.StringVar(&stateRoot, "state-root", "", "state trie root")
	flag.Parse()

	// Init our leveldb backed stack.
	ts := &trieStack{}
	dataDirectoryName := "data_dir/" + stateRoot
	ts.Stack, err = goque.OpenStack(dataDirectoryName)
	if err != nil {
		panic(err)
	}
	defer ts.Close()

	// DEBUG
	// Dump checking the contents of the stack
	for {
		item, err := ts.Pop()
		if err == goque.ErrEmpty {
			break
		}
		fmt.Printf("%v\n", item.ToString())
	}
	// DEBUG

	// Init the traversal with the root
	_, err = ts.PushString(stateRoot)
	if err != nil {
		panic(err)
	}

	// Do the loop
	ts.runStateTrieTraversal()
}

// runStateTrieTraversal is the main loop
func (ts *trieStack) runStateTrieTraversal() {
	importedNodesCount := 0

	for {
		item, err := ts.Pop()
		if err == goque.ErrEmpty {
			fmt.Println("Stack Empty. We are done here :D")
			break
		}
		if err != nil {
			panic(err)
		}

		log.Printf("From the stack: %v\n", item.ToString())
		c := keccak256ToCid(MEthStateTrie, item.ToString()[2:])
		log.Printf("\t%v\n", c)

		nodeData, err := fetchTrieNodeData(c)
		if err != nil {
			ts.errorAndPushNodeBack(item, "Error during fething of the data", err)
			continue
		}

		err = ts.processRawData(nodeData)
		if err != nil {
			ts.errorAndPushNodeBack(item, "Error during rawdata processing", err)
			continue
		}

		err = importTrieNode(c, nodeData)
		if err != nil {
			// We don't want to push back to the stack.
			// For now, we will just log it
			log.Printf("\t\t\tThere was an error importing your node: %v\n", err)
		}

		importedNodesCount++
		log.Printf("\t\t\tNode imported. Count = %d\n", importedNodesCount)
	}
}

// errorAndPushNodeBack is a helper of runStateTrieTraversal
func (ts *trieStack) errorAndPushNodeBack(item *goque.Item, message string, err error) {
	log.Printf("\t%s: %v", message, err)

	log.Printf("\tTaking node back to the stack.")
	_, err = ts.PushString(item.ToString())
	if err != nil {
		// If the stack is not working, just panic.
		panic(err)
	}

	// Give it some time before trying it again
	time.Sleep(1 * time.Second)
}

// processRawData will examine the given data to identify what kind of node
// we are dealing with. If children are found, they will be added to the stack.
// In all cases the node is stored if there are no errors processing it.
func (ts *trieStack) processRawData(encoded []byte) error {
	var decoded []interface{}
	err := rlp.DecodeBytes(encoded, &decoded)
	if err != nil {
		return err
	}

	// What kind of node are we dealing with?
	switch len(decoded) {
	case 2:
		val := decoded[1].([]byte)

		if len(val) == 32 {
			// This node is referencing to another one below,
			// so we add that one to the stack
			log.Printf("\t\tAdding 0x%x (from a 'case 2') to the stack\n", val)
			_, err = ts.PushString(fmt.Sprintf("0x%x", val))
			if err != nil {
				panic(err)
			}
		} else {
			// A proper leaf!
			// Just log it, as we will store the raw contents anyways
			log.Printf("\t\tThis is a leaf")
		}
	case 17:
		// A node with 16 children
		for i, _ := range decoded {
			// Adding them in reverse order
			vb := decoded[17-1-i].([]byte)
			switch len(vb) {
			case 0:
				// No further action here, is a nil.
				continue
			case 32:
				// Add this child to the stack
				log.Printf("\t\tAdding 0x%x (idx: %x) to the stack\n", vb, 17-1-i)
				_, err = ts.PushString(fmt.Sprintf("0x%x", vb))
				if err != nil {
					panic(err)
				}
			default:
				// Is either nothing or a proper 32 bytes element.
				return fmt.Errorf("unrecognized object in trie: %x", vb)
			}
		}
	default:
		return fmt.Errorf("unknown trie node type")
	}

	return nil
}

/*
  UTILS
*/

// keccak256ToCid takes the string representation of a keccak-256 hash
// and its codec to deliver its IPLD cid
func keccak256ToCid(codec uint64, s string) *cid.Cid {
	h, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}

	mhash, err := mh.Encode(h[:], mh.KECCAK_256)
	if err != nil {
		panic(err)
	}

	return cid.NewCidV1(codec, mhash)
}

// fetchTrieNodeData will try to get the rawdata of the state trie node
// from the given parity IPFS JSON RPC.
// You can also fetch this very data from an IPFS node (if available),
// as the API is the same.
func fetchTrieNodeData(c *cid.Cid) ([]byte, error) {
	client := &http.Client{}

	// TODO
	// Should be gotten from the command options, or default.
	ethapi := "http://localhost:15001/"

	url := ethapi + "api/v0/block/get?arg=" + c.String()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// importTrieNode stores your node in the selected medium
func importTrieNode(c *cid.Cid, rawdata []byte) error {
	// TODO
	// IMPLEMENT!
	return nil
}
