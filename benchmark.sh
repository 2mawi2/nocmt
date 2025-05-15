#!/bin/bash
set -e

cd "$(dirname "$0")"

BENCH_DIR="./benchmarks"
mkdir -p "$BENCH_DIR"

SCRIPT_DIR="$(pwd)"
BENCH_DIR="$SCRIPT_DIR/benchmarks"

BASELINE="$BENCH_DIR/baseline.txt"
TODAY=$(date +"%Y-%m-%d")
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
NEW_RESULT="$BENCH_DIR/bench_${TIMESTAMP}.txt"
LATEST="$BENCH_DIR/latest.txt"
HISTORY="$BENCH_DIR/history.txt"

run_benchmarks() {
    echo "Running benchmarks..."
    cd "$SCRIPT_DIR/processor" && go test -run=^$ -bench=. -benchmem -count=3 -benchtime=0.5s > "$1"
    echo "Benchmarks saved to $1"
    cd "$SCRIPT_DIR"
}

run_quick_benchmarks() {
    echo "Running quick benchmarks..."
    cd "$SCRIPT_DIR/processor" && go test -run=^$ -bench="BenchmarkStripComments|BenchmarkParallelProcessing" -benchmem -count=1 -benchtime=0.2s > "$1"
    echo "Quick benchmarks saved to $1"
    cd "$SCRIPT_DIR"
}

add_to_history() {
    echo "=== Benchmark run: $(date) ===" >> "$HISTORY"
    echo "File: $1" >> "$HISTORY"
    
    echo "Summary:" >> "$HISTORY"
    grep -E "Benchmark(StripComments|ParallelProcessing)" "$1" | awk '
        { 
            if ($1 ~ /BenchmarkStripComments\/LargeCode/) {
                printf "  Large code: %.2f ms/op, %.2f MB memory, %d allocs/op\n", 
                    $3/1000000, $5/1048576, $7
            } else if ($1 ~ /BenchmarkParallelProcessing/) {
                printf "  Parallel: %.2f ms/op, %.2f MB memory, %d allocs/op\n", 
                    $3/1000000, $5/1048576, $7
            }
        }
    ' "$1" >> "$HISTORY"
    
    echo "" >> "$HISTORY"
}

compare_results() {
    base=$1
    new=$2
    
    echo "Comparing benchmark results:"
    echo "============================"
    echo "Baseline: $(basename $base)"
    echo "Current:  $(basename $new)"
    echo ""
    
    benchmarks=$(grep -E "Benchmark" $new | awk '{print $1}' | sort | uniq)
    
    for bench in $benchmarks; do
        echo "Benchmark: $bench"
        
        if grep -q "$bench" $base; then
            base_time=$(grep "$bench" $base | awk '{sum+=$3; count++} END {print sum/count}')
            base_bytes=$(grep "$bench" $base | awk '{sum+=$5; count++} END {print sum/count}')
            base_allocs=$(grep "$bench" $base | awk '{sum+=$7; count++} END {print sum/count}')
            
            new_time=$(grep "$bench" $new | awk '{sum+=$3; count++} END {print sum/count}')
            new_bytes=$(grep "$bench" $new | awk '{sum+=$5; count++} END {print sum/count}')
            new_allocs=$(grep "$bench" $new | awk '{sum+=$7; count++} END {print sum/count}')
            
            time_diff=$(awk "BEGIN {printf \"%.2f\", (($new_time-$base_time)/$base_time)*100}")
            bytes_diff=$(awk "BEGIN {printf \"%.2f\", (($new_bytes-$base_bytes)/$base_bytes)*100}")
            allocs_diff=$(awk "BEGIN {printf \"%.2f\", (($new_allocs-$base_allocs)/$base_allocs)*100}")
            
            echo "  Time:   $(printf "%.2f" $new_time) ns/op   ($(color_diff $time_diff)%)"
            echo "  Memory: $(printf "%.2f" $new_bytes) B/op    ($(color_diff $bytes_diff)%)"
            echo "  Allocs: $(printf "%.2f" $new_allocs) allocs/op ($(color_diff $allocs_diff)%)"
        else
            new_time=$(grep "$bench" $new | awk '{sum+=$3; count++} END {print sum/count}')
            new_bytes=$(grep "$bench" $new | awk '{sum+=$5; count++} END {print sum/count}')
            new_allocs=$(grep "$bench" $new | awk '{sum+=$7; count++} END {print sum/count}')
            
            echo "  Time:   $(printf "%.2f" $new_time) ns/op   (new benchmark)"
            echo "  Memory: $(printf "%.2f" $new_bytes) B/op    (new benchmark)"
            echo "  Allocs: $(printf "%.2f" $new_allocs) allocs/op (new benchmark)"
        fi
        echo ""
    done
}

color_diff() {
    diff=$1
    if awk "BEGIN {exit !($diff < 0)}"; then
        echo -e "\033[32m$diff\033[0m"
    elif awk "BEGIN {exit !($diff > 0)}"; then
        echo -e "\033[31m$diff\033[0m"
    else
        echo -e "\033[33m$diff\033[0m"
    fi
}

case "$1" in
    run)
        run_benchmarks "$NEW_RESULT"
        
        cp "$NEW_RESULT" "$BASELINE"
        cp "$NEW_RESULT" "$LATEST"
        
        add_to_history "$NEW_RESULT"
        
        echo "New baseline established."
        ;;
        
    quick)
        run_quick_benchmarks "$NEW_RESULT"
        
        cp "$NEW_RESULT" "$LATEST"
        
        add_to_history "$NEW_RESULT"
        
        echo ""
        echo "Quick benchmark results:"
        grep -E "Benchmark(StripComments|ParallelProcessing)" "$NEW_RESULT" | awk '
            { 
                if ($1 ~ /BenchmarkStripComments\/LargeCode/) {
                    printf "  Large code: %.2f ms/op, %.2f MB memory, %d allocs/op\n", 
                        $3/1000000, $5/1048576, $7
                } else if ($1 ~ /BenchmarkParallelProcessing/) {
                    printf "  Parallel: %.2f ms/op, %.2f MB memory, %d allocs/op\n", 
                        $3/1000000, $5/1048576, $7
                }
            }
        ' "$NEW_RESULT"
        ;;
        
    compare)
        COMPARE_WITH="$BASELINE"
        if [ ! -z "$2" ]; then
            COMPARE_WITH="$2"
        fi
        
        if [ ! -f "$COMPARE_WITH" ]; then
            echo "Error: Comparison file $COMPARE_WITH not found."
            echo "Run './benchmark.sh run' first to establish a baseline."
            exit 1
        fi
        
        run_benchmarks "$NEW_RESULT"
        cp "$NEW_RESULT" "$LATEST"
        
        add_to_history "$NEW_RESULT"
        
        echo ""
        compare_results "$COMPARE_WITH" "$NEW_RESULT"
        ;;
        
    history)
        if [ -f "$HISTORY" ]; then
            cat "$HISTORY"
        else
            echo "No benchmark history found."
            echo "Run './benchmark.sh run' to establish a baseline."
        fi
        ;;
        
    *)
        echo "Usage:"
        echo "  ./benchmark.sh run              - Run benchmarks and save as new baseline"
        echo "  ./benchmark.sh quick            - Run a quick benchmark (faster but less accurate)"
        echo "  ./benchmark.sh compare [file]   - Compare current performance with baseline"
        echo "  ./benchmark.sh history          - Show history of benchmark results"
        ;;
esac