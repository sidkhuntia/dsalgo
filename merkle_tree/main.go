package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
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
	Left  *MerkleNode `json:"left,omitempty"`
	Right *MerkleNode `json:"right,omitempty"`
	Hash  []byte      `json:"hash"`
}

type MerkleTree struct {
	Root      *MerkleNode `json:"root"`
	CreatedAt time.Time   `json:"created_at"`
	FileCount int         `json:"file_count"`
	RootHash  string      `json:"root_hash"`
}

func (m *MerkleTree) Print() {
	fmt.Printf("Merkle Tree Root Hash: %s\n", hex.EncodeToString(m.Root.Hash))
}

func (m *MerkleTree) ToJSON() ([]byte, error) {
	if m.Root != nil {
		m.RootHash = hex.EncodeToString(m.Root.Hash)
	}
	return json.MarshalIndent(m, "", "  ")
}

func (m *MerkleTree) SaveToFile(filename string) error {
	jsonData, err := m.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to serialize tree: %v", err)
	}

	return os.WriteFile(filename, jsonData, 0644)
}

func LoadMerkleTreeFromFile(filename string) (*MerkleTree, error) {
	jsonData, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	var tree MerkleTree
	err = json.Unmarshal(jsonData, &tree)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return &tree, nil
}

func (m *MerkleTree) isEqual(other *MerkleTree) bool {
	return m.RootHash == other.RootHash
}

func compareNodes(a, b *MerkleNode) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Compare hashes
	if len(a.Hash) != len(b.Hash) {
		return false
	}
	for i, v := range a.Hash {
		if v != b.Hash[i] {
			return false
		}
	}

	// Recursively compare children
	return compareNodes(a.Left, b.Left) && compareNodes(a.Right, b.Right)
}

func (m *MerkleTree) Compare(other *MerkleTree) {
	fmt.Println("=== Merkle Tree Comparison ===")

	if m.isEqual(other) {
		fmt.Println("✅ Trees are IDENTICAL")
		fmt.Printf("Root Hash: %s\n", hex.EncodeToString(m.Root.Hash))
		return
	}

	fmt.Println("❌ Trees are DIFFERENT")

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

	tree := &MerkleTree{
		Root:      nodes[0],
		CreatedAt: time.Now(),
		FileCount: len(data),
	}

	// Set the root hash string
	if tree.Root != nil {
		tree.RootHash = hex.EncodeToString(tree.Root.Hash)
	}

	return tree
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
	// Define flags
	var (
		compareJSON = flag.String("compare", "", "Path to JSON file containing previous Merkle tree for comparison")
		saveJSON    = flag.String("save", "", "Path to save current Merkle tree as JSON")
		loadJSON    = flag.String("load", "", "Path to load Merkle tree from JSON file")
		showHelp    = flag.Bool("h", false, "Show help message")
	)

	flag.Parse()

	if *showHelp {
		fmt.Println("Merkle Tree CLI Tool")
		fmt.Println("Usage:")
		fmt.Println("  Build from files:     go run main.go [files...]")
		fmt.Println("  Build from directory: go run main.go [directory]")
		fmt.Println("  Compare with JSON:    go run main.go -compare=old.json [files...]")
		fmt.Println("  Save to JSON:         go run main.go -save=tree.json [files...]")
		fmt.Println("  Load from JSON:       go run main.go -load=tree.json")
		fmt.Println("")
		fmt.Println("Flags:")
		flag.PrintDefaults()
		return
	}

	args := flag.Args()

	// Handle load JSON case
	if *loadJSON != "" {
		tree, err := LoadMerkleTreeFromFile(*loadJSON)
		if err != nil {
			fmt.Printf("Error loading JSON: %v\n", err)
			return
		}

		fmt.Println("=== Loaded Merkle Tree ===")
		tree.Print()
		fmt.Printf("File Count: %d\n", tree.FileCount)
		fmt.Printf("Created At: %s\n", tree.CreatedAt.Format(time.RFC3339))
		return
	}

	// Build new tree from files
	var data [][]byte
	var err error

	if len(args) > 1 {
		data, err = hashDirectFilePaths(args)
		if err != nil {
			fmt.Printf("Error getting direct filepaths: %v\n", err)
			return
		}
	} else if len(args) == 1 {
		data, err = hashFilesInDirectory(args[0])
		if err != nil {
			fmt.Printf("Error hashing files: %v\n", err)
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

	fmt.Println("=== New Merkle Tree ===")
	tree.Print()
	fmt.Printf("File Count: %d\n", tree.FileCount)
	fmt.Printf("Created At: %s\n", tree.CreatedAt.Format(time.RFC3339))

	// Save to JSON if requested
	if *saveJSON != "" {
		err := tree.SaveToFile(*saveJSON)
		if err != nil {
			fmt.Printf("Error saving JSON: %v\n", err)
		} else {
			fmt.Printf("✅ Saved tree to %s\n", *saveJSON)
		}
	}

	// Compare with existing JSON if requested
	if *compareJSON != "" {
		fmt.Println()
		oldTree, err := LoadMerkleTreeFromFile(*compareJSON)
		if err != nil {
			fmt.Printf("Error loading comparison JSON: %v\n", err)
			return
		}

		tree.Compare(oldTree)

		// Show detailed comparison
		if !tree.isEqual(oldTree) {
			fmt.Println("\n=== Detailed Analysis ===")
			if tree.FileCount != oldTree.FileCount {
				fmt.Printf("📊 File count changed: %d → %d\n", oldTree.FileCount, tree.FileCount)
			}

		}
	}
}
