package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
)

type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Hash  []byte
}

type MerkleTree struct {
	Root *MerkleNode
}

func hashData(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	hash := sha256.New()
	if left != nil {
		hash.Write(left.Hash)
	}
	if right != nil {
		hash.Write(right.Hash)
	}
	if data != nil {
		hash.Write(data)
	}
	hashValue := hash.Sum(nil)

	return &MerkleNode{Left: left, Right: right, Hash: hashValue}
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

func hashFiles(files []string) ([][]byte, error) {
	var data [][]byte
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}
		hashedContent := hashData(content)
		data = append(data, hashedContent)
	}
	return data, nil
}

func main() {

	args := os.Args[1:]

	var data [][]byte
	var err error
	if len(args) > 0 {
		data, err = hashFiles(args)
		if err != nil {
			fmt.Println("Error hashing files:", err)
			return
		}

	} else {
		fmt.Println("No files provided")
		return
	}

	tree := buildMerkleTree(data)

	fmt.Println("Merkle Tree Root Hash:", hex.EncodeToString(tree.Root.Hash))

}
