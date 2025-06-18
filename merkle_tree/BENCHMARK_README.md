# Merkle Tree Benchmark Suite

This benchmark suite provides comprehensive performance testing for the Merkle tree implementation with configurable parameters.

## Quick Start

### Basic Benchmark Run
```bash
go test -bench=.
```

### Run with Custom Parameters
```bash
go test -bench=. -benchtime=10s -args 100 2048 8 60
```
**Parameters**: 100 files, 2KB each, 8 workers, 60 second timeout

## Available Benchmarks

### Core Functionality Benchmarks

| Benchmark | Description |
|-----------|-------------|
| `BenchmarkHashFiles` | Overall file hashing performance with configurable parameters |
| `BenchmarkHashSingleFile` | Single file hashing performance |
| `BenchmarkBuildMerkleTree` | Merkle tree construction only (no file I/O) |
| `BenchmarkEndToEnd` | Complete pipeline: hash files + build tree |
| `BenchmarkDirectoryTraversal` | Directory scanning performance |

### Scaling Benchmarks

#### File Size Scaling
```bash
go test -bench=BenchmarkFileSize
```
- `BenchmarkFileSize1KB` - 1KB files
- `BenchmarkFileSize10KB` - 10KB files  
- `BenchmarkFileSize100KB` - 100KB files
- `BenchmarkFileSize1MB` - 1MB files
- `BenchmarkFileSize10MB` - 10MB files

#### File Count Scaling
```bash
go test -bench=BenchmarkFiles
```
- `BenchmarkFiles10` - 10 files
- `BenchmarkFiles50` - 50 files
- `BenchmarkFiles100` - 100 files
- `BenchmarkFiles500` - 500 files
- `BenchmarkFiles1000` - 1000 files

#### Worker Scaling
```bash
go test -bench=BenchmarkWorkers
```
- `BenchmarkWorkers1` - 1 worker (sequential)
- `BenchmarkWorkers2` - 2 workers
- `BenchmarkWorkers4` - 4 workers
- `BenchmarkWorkers8` - 8 workers
- `BenchmarkWorkers16` - 16 workers

### Memory Benchmarks
```bash
go test -bench=BenchmarkMerkleTreeMemory
```
- `BenchmarkMerkleTreeMemory` - Memory allocation tracking

## Command Line Parameters

### Custom Configuration Format
```bash
go test -bench=<pattern> -benchtime=<duration> -args <numFiles> <fileSize> <workers> <timeout>
```

### Parameter Details

| Parameter | Description | Default | Example |
|-----------|-------------|---------|---------|
| `numFiles` | Number of files to process | 10 | 100 |
| `fileSize` | Size of each file in bytes | 1048576 (1MB) | 2048 (2KB) |
| `workers` | Number of worker goroutines | 4 | 8 |
| `timeout` | Timeout in seconds | 30 | 60 |

## Example Usage Scenarios

### 1. Test Small Files Performance
```bash
# Test with many small files
go test -bench=BenchmarkHashFiles -args 1000 1024 8 30
```

### 2. Test Large Files Performance  
```bash
# Test with few large files
go test -bench=BenchmarkHashFiles -args 10 10485760 4 120
```

### 3. Worker Scaling Analysis
```bash
# Compare different worker counts
go test -bench=BenchmarkWorkers -v
```

### 4. Memory Usage Analysis
```bash
# Track memory allocations
go test -bench=BenchmarkMerkleTreeMemory -benchmem
```

### 5. Comprehensive Performance Test
```bash
# Run all benchmarks with detailed output
go test -bench=. -v -benchmem -benchtime=5s
```

## Advanced Usage

### Run Specific Benchmark Categories
```bash
# Only file size benchmarks
go test -bench=BenchmarkFileSize

# Only worker scaling benchmarks  
go test -bench=BenchmarkWorkers

# Only memory benchmarks
go test -bench=Memory
```

### Performance Profiling
```bash
# CPU profiling
go test -bench=BenchmarkHashFiles -cpuprofile=cpu.prof

# Memory profiling
go test -bench=BenchmarkHashFiles -memprofile=mem.prof

# View profiles
go tool pprof cpu.prof
go tool pprof mem.prof
```

### Benchmark Output Analysis

#### Reading Benchmark Results
```
BenchmarkHashFiles-8    	       3	 415389533 ns/op	  12345 B/op	     567 allocs/op
```

**Breakdown**:
- `BenchmarkHashFiles-8`: Benchmark name with GOMAXPROCS
- `3`: Number of iterations run
- `415389533 ns/op`: Nanoseconds per operation
- `12345 B/op`: Bytes allocated per operation  
- `567 allocs/op`: Number of allocations per operation

#### Performance Metrics to Watch
- **Throughput**: Lower ns/op = better performance
- **Memory Efficiency**: Lower B/op and allocs/op = better
- **Scalability**: Performance should improve with more workers (up to CPU count)
- **Context Cancellation**: Operations should stop quickly on timeout

## Example Results Analysis

### Optimal Worker Count
```bash
go test -bench=BenchmarkWorkers -v | grep -E "(BenchmarkWorkers|ns/op)"
```
Look for the worker count with lowest ns/op.

### File Size Performance Characteristics
```bash
go test -bench=BenchmarkFileSize -v
```
Analyze how performance scales with file size.

## Troubleshooting

### Common Issues

1. **Timeout Errors**: Increase timeout parameter for large files
2. **Memory Issues**: Reduce file count or size for memory-constrained systems
3. **Inconsistent Results**: Use `-benchtime=10s` or higher for stable results

### Debug Mode
```bash
# Verbose output with configuration logging
go test -bench=BenchmarkHashFiles -v -args 10 1024 4 30
```

## Integration with CI/CD

### Regression Testing
```bash
# Save baseline performance
go test -bench=. -benchtime=10s > baseline.txt

# Compare against baseline
go test -bench=. -benchtime=10s > current.txt
# Use benchcmp tool to compare
```

### Performance Gates
```bash
# Fail if performance degrades significantly
go test -bench=BenchmarkEndToEnd -benchtime=5s | \
  awk '/BenchmarkEndToEnd/ { if ($3 > 1000000000) exit 1 }'
```

This benchmark suite provides comprehensive performance analysis capabilities for your Merkle tree implementation, allowing you to optimize for different use cases and system configurations. 