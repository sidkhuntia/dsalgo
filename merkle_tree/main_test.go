package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"
)

// Benchmark configuration
type BenchmarkConfig struct {
	NumFiles    int
	FileSize    int // in bytes
	WorkerCount int
	Timeout     time.Duration
}

// Helper function to create benchmark test files
func createBenchmarkFiles(b *testing.B, config BenchmarkConfig) ([]string, func()) {
	tempDir, err := ioutil.TempDir("", "merkle_bench_")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}

	var files []string
	for i := 0; i < config.NumFiles; i++ {
		filename := filepath.Join(tempDir, fmt.Sprintf("file_%d.dat", i))

		// Create file with random data
		data := make([]byte, config.FileSize)
		if _, err := rand.Read(data); err != nil {
			b.Fatalf("Failed to generate random data: %v", err)
		}

		if err := ioutil.WriteFile(filename, data, 0644); err != nil {
			b.Fatalf("Failed to create test file: %v", err)
		}

		files = append(files, filename)
	}

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return files, cleanup
}

// Helper function to create deterministic data
func createDeterministicData(numHashes int, dataSize int) [][]byte {
	data := make([][]byte, numHashes)
	for i := 0; i < numHashes; i++ {
		chunk := make([]byte, dataSize)
		// Fill with deterministic pattern
		for j := range chunk {
			chunk[j] = byte((i + j) % 256)
		}
		data[i] = chunk
	}
	return data
}

// Parse benchmark arguments from test flags
func parseBenchmarkArgs() BenchmarkConfig {
	// Default configuration
	config := BenchmarkConfig{
		NumFiles:    10,
		FileSize:    1024 * 1024, // 1MB
		WorkerCount: 4,
		Timeout:     30 * time.Second,
	}

	// You can pass arguments like: go test -bench=. -benchtime=10s -args 100 2048 8 60
	// This would set: 100 files, 2KB each, 8 workers, 60 second timeout
	if len(os.Args) > 5 {
		if numFiles, err := strconv.Atoi(os.Args[5]); err == nil {
			config.NumFiles = numFiles
		}
	}
	if len(os.Args) > 6 {
		if fileSize, err := strconv.Atoi(os.Args[6]); err == nil {
			config.FileSize = fileSize
		}
	}
	if len(os.Args) > 7 {
		if workers, err := strconv.Atoi(os.Args[7]); err == nil {
			config.WorkerCount = workers
		}
	}
	if len(os.Args) > 8 {
		if timeout, err := strconv.Atoi(os.Args[8]); err == nil {
			config.Timeout = time.Duration(timeout) * time.Second
		}
	}

	return config
}

// Benchmark: File hashing operations
func BenchmarkHashFiles(b *testing.B) {
	config := parseBenchmarkArgs()

	b.Logf("Config: %d files, %d bytes each, %d workers, %v timeout",
		config.NumFiles, config.FileSize, config.WorkerCount, config.Timeout)

	files, cleanup := createBenchmarkFiles(b, config)
	defer cleanup()

	// Reset timer after setup
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := hashFilesWithTimeout(files, config.Timeout)
		if err != nil {
			b.Fatalf("hashFiles failed: %v", err)
		}
	}
}

// Benchmark: Single file hashing
func BenchmarkHashSingleFile(b *testing.B) {
	config := parseBenchmarkArgs()

	// Create just one file for this benchmark
	tempFile, err := ioutil.TempFile("", "single_bench_*.dat")
	if err != nil {
		b.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write test data
	data := make([]byte, config.FileSize)
	rand.Read(data)
	tempFile.Write(data)
	tempFile.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		_, err := hashFile(ctx, tempFile.Name())
		if err != nil {
			b.Fatalf("hashFile failed: %v", err)
		}
	}
}

// Benchmark: Merkle tree construction
func BenchmarkBuildMerkleTree(b *testing.B) {
	config := parseBenchmarkArgs()

	// Create deterministic data for consistent benchmarking
	data := createDeterministicData(config.NumFiles, 32) // 32 bytes per hash (SHA256)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree := buildMerkleTree(data)
		if tree == nil {
			b.Fatal("buildMerkleTree returned nil")
		}
	}
}

// Benchmark: End-to-end pipeline
func BenchmarkEndToEnd(b *testing.B) {
	config := parseBenchmarkArgs()

	files, cleanup := createBenchmarkFiles(b, config)
	defer cleanup()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Hash files
		hashes, err := hashFilesWithTimeout(files, config.Timeout)
		if err != nil {
			b.Fatalf("hashFiles failed: %v", err)
		}

		// Build Merkle tree
		tree := buildMerkleTree(hashes)
		if tree == nil {
			b.Fatal("buildMerkleTree returned nil")
		}
	}
}

// Benchmark: Directory traversal
func BenchmarkDirectoryTraversal(b *testing.B) {
	config := parseBenchmarkArgs()

	// Create a directory structure
	tempDir, err := ioutil.TempDir("", "dir_bench_")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create subdirectories
	for i := 0; i < 5; i++ {
		subDir := filepath.Join(tempDir, fmt.Sprintf("subdir_%d", i))
		os.MkdirAll(subDir, 0755)

		// Create files in each subdirectory
		for j := 0; j < config.NumFiles/5; j++ {
			filename := filepath.Join(subDir, fmt.Sprintf("file_%d.dat", j))
			data := make([]byte, config.FileSize)
			rand.Read(data)
			ioutil.WriteFile(filename, data, 0644)
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := getAllFilesInDirectory(tempDir)
		if err != nil {
			b.Fatalf("getAllFilesInDirectory failed: %v", err)
		}
	}
}

// Benchmark: Different file sizes
func BenchmarkFileSize1KB(b *testing.B)   { benchmarkByFileSize(b, 1024) }
func BenchmarkFileSize10KB(b *testing.B)  { benchmarkByFileSize(b, 10*1024) }
func BenchmarkFileSize100KB(b *testing.B) { benchmarkByFileSize(b, 100*1024) }
func BenchmarkFileSize1MB(b *testing.B)   { benchmarkByFileSize(b, 1024*1024) }
func BenchmarkFileSize10MB(b *testing.B)  { benchmarkByFileSize(b, 10*1024*1024) }

func benchmarkByFileSize(b *testing.B, fileSize int) {
	config := parseBenchmarkArgs()
	config.FileSize = fileSize
	config.NumFiles = 10 // Keep number of files constant

	files, cleanup := createBenchmarkFiles(b, config)
	defer cleanup()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := hashFilesWithTimeout(files, config.Timeout)
		if err != nil {
			b.Fatalf("hashFiles failed: %v", err)
		}
	}
}

// Benchmark: Different number of files
func BenchmarkFiles10(b *testing.B)   { benchmarkByFileCount(b, 10) }
func BenchmarkFiles50(b *testing.B)   { benchmarkByFileCount(b, 50) }
func BenchmarkFiles100(b *testing.B)  { benchmarkByFileCount(b, 100) }
func BenchmarkFiles500(b *testing.B)  { benchmarkByFileCount(b, 500) }
func BenchmarkFiles1000(b *testing.B) { benchmarkByFileCount(b, 1000) }

func benchmarkByFileCount(b *testing.B, numFiles int) {
	config := parseBenchmarkArgs()
	config.NumFiles = numFiles
	config.FileSize = 1024 // Keep file size constant at 1KB

	files, cleanup := createBenchmarkFiles(b, config)
	defer cleanup()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := hashFilesWithTimeout(files, config.Timeout)
		if err != nil {
			b.Fatalf("hashFiles failed: %v", err)
		}
	}
}

// Benchmark: Worker scaling
func BenchmarkWorkers1(b *testing.B)  { benchmarkByWorkerCount(b, 1) }
func BenchmarkWorkers2(b *testing.B)  { benchmarkByWorkerCount(b, 2) }
func BenchmarkWorkers4(b *testing.B)  { benchmarkByWorkerCount(b, 4) }
func BenchmarkWorkers8(b *testing.B)  { benchmarkByWorkerCount(b, 8) }
func BenchmarkWorkers16(b *testing.B) { benchmarkByWorkerCount(b, 16) }

func benchmarkByWorkerCount(b *testing.B, workers int) {
	config := parseBenchmarkArgs()
	config.WorkerCount = workers

	files, cleanup := createBenchmarkFiles(b, config)
	defer cleanup()

	// Create a custom version of hashFiles that uses specific worker count
	hashFilesCustomWorkers := func(files []string, workerCount int) ([][]byte, error) {
		if len(files) == 0 {
			return nil, fmt.Errorf("no files provided")
		}

		jobs := make(chan string, len(files))
		results := make(chan HashResult, len(files))
		errors := make(chan error, len(files))

		ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
		defer cancel()

		wg := sync.WaitGroup{}

		// Use custom worker count
		for range workerCount {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					select {
					case job, ok := <-jobs:
						if !ok {
							return
						}
						hash, err := hashFile(ctx, job)
						if err != nil {
							errors <- err
							cancel()
							return
						}
						results <- HashResult{File: job, Hash: hash}
					case <-ctx.Done():
						return
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

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := hashFilesCustomWorkers(files, workers)
		if err != nil {
			b.Fatalf("hashFiles failed: %v", err)
		}
	}
}

// Memory allocation benchmarks
func BenchmarkMerkleTreeMemory(b *testing.B) {
	config := parseBenchmarkArgs()
	data := createDeterministicData(config.NumFiles, 32)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tree := buildMerkleTree(data)
		if tree == nil {
			b.Fatal("buildMerkleTree returned nil")
		}
	}
}

// Benchmark suite provides comprehensive performance testing
// See BENCHMARK_README.md for usage instructions
