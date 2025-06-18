package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
)

type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Hash  []byte
}

type MerkleTree struct {
	Root *MerkleNode
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

	workers := runtime.NumCPU()

	if len(files) < workers {
		workers = len(files)
	}

	jobs := make(chan string, len(files))
	results := make(chan []byte, len(files))

	wg := sync.WaitGroup{}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				hash, err := hashFile(job)
				if err != nil {
					fmt.Println("Error hashing file:", err)
					continue
				}
				results <- hash
			}
		}()
	}

	go func() {
		defer close(jobs)
		for _, file := range files {

			jobs <- file
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	var data [][]byte
	for result := range results {
		data = append(data, result)
	}

	return data, nil
}

func hashFile(file string) ([]byte, error) {
	data, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer data.Close()

	stat, err := data.Stat()
	if err != nil {
		return nil, err
	}

	if stat.IsDir() {
		return nil, fmt.Errorf("is a directory")
	}

	hash := sha256.New()

	if stat.Size() <= 10*1024*1024 { // if file is less than 10MB, read the whole file
		content, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}
		hash.Write(content)
	} else {
		buffer := make([]byte, 1024*1024) // if file is greater than 10MB, read the file in chunks of 1MB
		for {
			n, err := data.Read(buffer)
			if n > 0 {
				hash.Write(buffer[:n])
			}

			if err == io.EOF {
				break
			}

			if err != nil {
				return nil, err
			}
		}
	}

	return hash.Sum(nil), nil
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
