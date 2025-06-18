package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"sync"
	"time"
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

	if len(nodes) == 0 {
		return nil
	}

	return &MerkleTree{Root: nodes[0]}
}

func getAllFilesInDirectory(directory string) ([]string, error) {
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	filenames := make([]string, 0, len(files))

	for _, file := range files {
		fullPath := filepath.Join(directory, file.Name())
		fullPath = filepath.Clean(fullPath)
		fullPath, err = filepath.Abs(fullPath)
		if err != nil {
			return nil, err
		}

		if file.IsDir() {
			subFiles, err := getAllFilesInDirectory(fullPath)
			if err != nil {
				return nil, err
			}
			filenames = append(filenames, subFiles...)
		} else {
			filenames = append(filenames, fullPath)
		}
	}

	return filenames, nil
}

func hashFilesInDirectory(directory string) ([][]byte, error) {

	filenames, err := getAllFilesInDirectory(directory)
	if err != nil {
		return nil, err
	}

	return hashFiles(filenames)
}

func hashDirectFilePaths(filenames []string) ([][]byte, error) {

	directFilePaths := make([]string, 0, len(filenames))

	for _, filename := range filenames {

		fileInfo, err := os.Stat(filename)
		if err != nil {
			return nil, err
		}

		if fileInfo.IsDir() {
			return nil, fmt.Errorf("cannot hash directories along with filepaths")
		}
		absPath, err := filepath.Abs(filename)
		if err != nil {
			return nil, err
		}
		directFilePaths = append(directFilePaths, absPath)
	}

	return hashFiles(directFilePaths)
}

func hashFiles(files []string) ([][]byte, error) {

	if len(files) == 0 {
		return nil, fmt.Errorf("no files provided")
	}

	// sort the files
	sort.Strings(files)

	fmt.Printf("Hashing %d files\n", len(files))

	for _, file := range files {
		fmt.Printf("Hashing file: %s\n", file)
	}

	return hashFilesWithTimeout(files, 30*time.Second) // 30 second default timeout
}

type HashResult struct {
	File string
	Hash []byte
}

func hashFilesWithTimeout(files []string, timeout time.Duration) ([][]byte, error) {
	workers := min(len(files), runtime.NumCPU())

	jobs := make(chan string, len(files))
	results := make(chan HashResult, len(files))
	errors := make(chan error, len(files))

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	wg := sync.WaitGroup{}

	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case job, ok := <-jobs:
					if !ok {
						return // jobs channel closed
					}
					hash, err := hashFile(ctx, job)
					if err != nil {
						errors <- err
						cancel()
						return
					}
					results <- HashResult{File: job, Hash: hash}
				case <-ctx.Done():
					return // Context cancelled, stop worker
				}
			}
		}()
	}

	go func() {
		defer close(jobs)
		for _, file := range files {
			select {
			case <-ctx.Done():
				return
			case jobs <- file:
			}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	var hashedFiles []HashResult
	var data [][]byte
	expectedResults := len(files)
	receivedResults := 0

	for receivedResults < expectedResults {
		select {
		case result, ok := <-results:
			if !ok {
				// if the results channel is closed, check if we got all results
				if receivedResults < expectedResults {
					return nil, fmt.Errorf("not all files processed successfully")
				}
				break
			}
			hashedFiles = append(hashedFiles, result)
			receivedResults++
		case err := <-errors:
			cancel()
			return nil, err
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	sort.Slice(hashedFiles, func(i, j int) bool {
		return hashedFiles[i].File < hashedFiles[j].File
	})

	for _, file := range hashedFiles {
		data = append(data, file.Hash)
	}

	return data, nil
}

func hashFile(ctx context.Context, file string) ([]byte, error) {

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		break
	}

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

	if stat.Size() <= 5*1024*1024 { // if file is less than 5MB, read the whole file
		content, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}
		hash.Write(content)
	} else if stat.Size() <= 50*1024*1024 { // if file is less than 50MB, read the file in chunks of 1MB
		buffer := make([]byte, 1024*1024)
		for {
			// Check if context is cancelled before each read
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}

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
	if len(args) > 1 {
		data, err = hashDirectFilePaths(args) // get the direct filepaths of the files
		if err != nil {
			fmt.Println("Error getting direct filepaths:", err)
			return
		}

	} else if len(args) == 1 {
		// if there is one arg, treat it as directory and hash all files in the directory
		data, err = hashFilesInDirectory(args[0])
		if err != nil {
			fmt.Println("Error hashing files:", err)
			return
		}

	} else {
		fmt.Println("No files provided")
		return
	}

	tree := buildMerkleTree(data)

	if tree == nil {
		fmt.Println("Could not build Merkle Tree")
		return
	}

	fmt.Println("Merkle Tree Root Hash:", hex.EncodeToString(tree.Root.Hash))

}
