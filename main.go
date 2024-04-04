package main

import (
	"crypto/sha1"
	"fmt"
	"hash"
	"slices"
	"strconv"
)

func main() {
	fmt.Println("hello, this is my attempt at implementing consistent hashing")
	// new server
	masterServer := NewMasterServer()
	// add nodes to the server
	// db nodes a map of node names to data
	nodeCount := 10
	for i := 0; i < nodeCount; i++ {
		node := Node{
			ID:          strconv.Itoa(i) + " Node",
			Name:        strconv.Itoa(i) + " Node Name",
			PrimaryData: make(map[string]string),
			ReplicaData: make(map[string]string),
		}
		masterServer.AddNode(node)
		masterServer.DBClient.AddNode(node)
	}

	// test, TODO move to test file
	// add data
	for i := 0; i < 10; i++ {
		masterServer.AddData(strconv.Itoa(i)+" something key", strconv.Itoa(i)+" something value")
	}
	for i := 0; i < 10; i++ {
		key := strconv.Itoa(i) + " something key"
		val := masterServer.GetData(key)
		fmt.Printf("key: %v, vale: %v\n", key, val)
	}

}

func NewMasterServer() *MasterServer {
	dbClient := DbClient{
		DataDB: make(map[string]Node),
	}
	master := MasterServer{
		HashRing:        make([]string, 0),
		HashIDNodeNames: map[string]string{},
		Hasher:          sha1.New(),
		DBClient:        dbClient,
	}

	return &master
}

type DbClient struct {
	DataDB map[string]Node // node names to Node
}

func (dbc *DbClient) AddNode(node Node) {
	dbc.DataDB[node.Name] = node
}

func (dbc *DbClient) AddPrimaryData(key, nodeName, value string) {
	// get the node
	node := dbc.DataDB[nodeName]
	node.PrimaryData[key] = value
}

func (dbc *DbClient) GetPrimaryData(key, nodeName string) string {
	node := dbc.DataDB[nodeName]
	return node.PrimaryData[key]
}

// master server contains all context to a node
type MasterServer struct {
	HashRing        []string          // as sorted slice of hashes that we can perform binary search on.
	HashIDNodeNames map[string]string // hashID to nodeNames
	Hasher          hash.Hash
	DBClient        DbClient
}

// AddNode adds a node to the hashRing which is a map of hashes to node names.
func (ms *MasterServer) AddNode(node Node) {
	// hash the name

	hashID := ms.Hasher.Sum([]byte(node.ID))
	// store the hash in a hash ring, this is sorted
	ms.HashRing = append(ms.HashRing, string(hashID))

	// a very inefficient way to sort the hash ring
	slices.Sort(ms.HashRing)
	// add it to the map
	ms.HashIDNodeNames[string(hashID)] = node.Name

	ms.Hasher.Reset() // reset the hash
}

func (ms *MasterServer) GetNodeOwner(key string) string {
	// hash the key
	hashedKey := ms.Hasher.Sum([]byte(key))
	ms.Hasher.Reset()
	// perform binary search on the sorted slice i.e. HashRing
	pos, found := slices.BinarySearch(ms.HashRing, string(hashedKey))
	if found {
		fmt.Println("this is a server!")
		return key
	} else if pos == len(ms.HashRing) {
		// if you are at the end the owner is actually the first node
		return ms.HashIDNodeNames[ms.HashRing[0]]
	} else {
		node := ms.HashIDNodeNames[ms.HashRing[pos-1]]
		// it would have been in position so the previous index is the owner
		fmt.Printf("it would have been in this position %v \ndata would be in node:%v\n", pos, node)
		return node
	}
}

func (ms *MasterServer) AddData(key string, value string) {
	// hash the key, get the owner of the key
	nodeName := ms.GetNodeOwner(key)
	// call the dbclient with the appropriate location and save it to that node
	ms.DBClient.AddPrimaryData(key, nodeName, value)
}

func (ms *MasterServer) GetData(key string) string {
	nodeName := ms.GetNodeOwner(key)
	val := ms.DBClient.GetPrimaryData(key, nodeName)
	return val
}

func (ms *MasterServer) ExpireNode(node Node) {
	// find the next node on the list
	//
}

type Node struct {
	ID          string
	Name        string
	PrimaryData map[string]string
	ReplicaData map[string]string
}

// add node
