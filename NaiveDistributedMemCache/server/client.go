package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type Client struct {
	mu              sync.Mutex
	consistentHash  HashRing
	peer2AddressMap map[string]string
	callMap         map[string]*Call
}

type Call struct {
	wg    sync.WaitGroup
	value []byte
	err   error
}

func NewClient(virtualReplicas int, hash func(data []byte) uint32) *Client {
	return &Client{
		mu:              sync.Mutex{},
		consistentHash:  *NewHashRing(virtualReplicas, hash),
		peer2AddressMap: make(map[string]string),
		callMap:         make(map[string]*Call),
	}
}

func (client *Client) pickPeerAddress(key string) string {
	if peer := client.consistentHash.GetNode(key); len(peer) != 0 {
		if address, ok := client.peer2AddressMap[peer]; ok {
			return address
		}
	}
	return ""
}

func sendRequest(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Peer has no such key")
	}
	if body, err := ioutil.ReadAll(response.Body); err != nil {
		return nil, err
	} else {
		return body, nil
	}
}

func (client *Client) getFromRemotePeer(group string, key string) ([]byte, error) {
	client.mu.Lock()
	if call, ok := client.callMap[key]; ok {
		client.mu.Unlock()
		call.wg.Wait()
		return call.value, call.err
	}
	call := new(Call)
	call.wg.Add(1)
	defer call.wg.Done()
	client.callMap[key] = call
	peer := client.pickPeerAddress(key)
	log.Printf("Get key %s from peer %s\n", key, peer)
	if len(peer) == 0 {
		return nil, fmt.Errorf("No such peer")
	}
	url := fmt.Sprintf("%v/cache/%v/%v", peer, group, key)
	client.mu.Unlock()
	call.value, call.err = sendRequest(url)

	client.mu.Lock()
	// The pointer of Call has already been gotten by blocked threads.
	// So it is safe to delete the map.
	delete(client.callMap, key)
	client.mu.Unlock()
	return call.value, call.err
}
