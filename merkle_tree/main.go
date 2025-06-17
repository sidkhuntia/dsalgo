package main

import (
	"fmt"
	"crypto/sha256"
	"encoding/hex"
)



type MerkleNode struct {
	Left *MerkleNode
	Right *MerkleNode
	Hash []byte
}

type MerkleTree struct {
	Root *MerkleNode
}

func hashData(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	hash := hashData(data)
	return &MerkleNode{Left: left, Right: right, Hash: hash}
}

func buildMerkleTree(data [][]byte) *MerkleTree {
	var nodes []*MerkleNode
	for _, d := range data {
		nodes = append(nodes, NewMerkleNode(nil, nil, d))
	}

	if len(nodes)%2 != 0 {
		nodes = append(nodes, nodes[len(nodes)-1])
	}

	for len(nodes) > 1 {
		var newNodes []*MerkleNode
		for i := 1; i < len(nodes); i += 2 {
			
			newNode := NewMerkleNode(nodes[i], nodes[i-1], nil)
			newNodes = append(newNodes, newNode)
		}
		if len(newNodes)%2 != 0 && len(newNodes) > 1 {
			newNodes = append(newNodes, newNodes[len(newNodes)-1])
		}
		nodes = newNodes
	}

	return &MerkleTree{Root: nodes[0]}
}


func main() {
	data := [][]byte{
		[]byte("node1"),
		[]byte("node2"),
		[]byte("node3"),
		[]byte("node4"),
	}

	tree := buildMerkleTree(data)

	fmt.Println("Merkle Tree Root Hash:", hex.EncodeToString(tree.Root.Hash))

}