package main

import (
	"fmt"
	"strconv"
	"testing"
)

func TestMasterServer(t *testing.T) {
	masterServer := NewMasterServer()

	// add nodes to the server
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

	expectToFind := make(map[string]string)
	// add data
	for i := 0; i < 100; i++ {
		key := strconv.Itoa(i) + "id"
		value := strconv.Itoa(i) + "data"
		masterServer.AddData(key, value)
		expectToFind[key] = value
	}

	// get data
	for key, value := range expectToFind {
		retrieved := masterServer.GetData(key)
		if value != retrieved {
			t.Errorf("unable to find expected value: %v", value)
		}
		fmt.Printf("key: %v, value: %v\n", key, retrieved)
	}

}
