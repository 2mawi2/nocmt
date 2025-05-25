#!/bin/bash
set -e

cd "$(dirname "$0")"

BENCHMARKS_DIRECTORY="./benchmarks"
mkdir -p "$BENCHMARKS_DIRECTORY"

SCRIPT_ROOT_DIRECTORY="$(pwd)"
BENCHMARKS_DIRECTORY="$SCRIPT_ROOT_DIRECTORY/benchmarks"

BASELINE_RESULTS_FILE="$BENCHMARKS_DIRECTORY/baseline.txt"
TODAY_DATE=$(date +"%Y-%m-%d")
CURRENT_TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
NEW_BENCHMARK_RESULTS="$BENCHMARKS_DIRECTORY/bench_${CURRENT_TIMESTAMP}.txt"
LATEST_RESULTS_FILE="$BENCHMARKS_DIRECTORY/latest.txt"
BENCHMARK_HISTORY_FILE="$BENCHMARKS_DIRECTORY/history.txt"

run_full_benchmark_suite() {
    echo "Running benchmarks..."
    cd "$SCRIPT_ROOT_DIRECTORY/internal/processor" && go test -run='^$' -bench=. -benchmem -count=3 -benchtime=0.5s > "$1"
    echo "Benchmarks saved to $1"
    cd "$SCRIPT_ROOT_DIRECTORY"
}

run_quick_benchmark_subset() {
    echo "Running quick benchmarks..."
    cd "$SCRIPT_ROOT_DIRECTORY/internal/processor" && go test -run='^$' -bench='BenchmarkGoStripComments|BenchmarkGoRealWorldCode|BenchmarkGoMemoryUsage|BenchmarkGoVeryLargeFiles' -benchmem -count=1 -benchtime=0.2s > "$1"
    echo "Quick benchmarks saved to $1"
    cd "$SCRIPT_ROOT_DIRECTORY"
}

append_benchmark_results_to_history() {
    benchmark_results_file="$1"
    echo "=== Benchmark run: $(date) ===" >> "$BENCHMARK_HISTORY_FILE"
    echo "File: $benchmark_results_file" >> "$BENCHMARK_HISTORY_FILE"
    
    echo "Summary:" >> "$BENCHMARK_HISTORY_FILE"
    extract_key_benchmark_metrics_for_history "$benchmark_results_file"
    echo "" >> "$BENCHMARK_HISTORY_FILE"
}

extract_key_benchmark_metrics_for_history() {
    benchmark_file="$1"
    grep -E "Benchmark(GoStripComments|GoRealWorldCode|GoMemoryUsage)" "$benchmark_file" | awk '
        { 
            if ($1 ~ /BenchmarkGoStripComments\/LargeCode/) {
                printf "  Large code: %.2f ms/op, %.2f MB memory, %d allocs/op\n", 
                    $3/1000000, $5/1048576, $7
            } else if ($1 ~ /BenchmarkGoRealWorldCode/) {
                printf "  Real world: %.2f ms/op, %.2f MB memory, %d allocs/op\n", 
                    $3/1000000, $5/1048576, $7
            } else if ($1 ~ /BenchmarkGoMemoryUsage/) {
                printf "  Memory usage: %.2f ms/op, %.2f MB memory, %d allocs/op\n", 
                    $3/1000000, $5/1048576, $7
            }
        }
    ' "$benchmark_file" >> "$BENCHMARK_HISTORY_FILE"
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
        run_full_benchmark_suite "$NEW_BENCHMARK_RESULTS"
        
        cp "$NEW_BENCHMARK_RESULTS" "$BASELINE_RESULTS_FILE"
        cp "$NEW_BENCHMARK_RESULTS" "$LATEST_RESULTS_FILE"
        
        append_benchmark_results_to_history "$NEW_BENCHMARK_RESULTS"
        
        echo "New baseline established."
        ;;
        
    quick)
        run_quick_benchmark_subset "$NEW_BENCHMARK_RESULTS"
        
        cp "$NEW_BENCHMARK_RESULTS" "$LATEST_RESULTS_FILE"
        
        append_benchmark_results_to_history "$NEW_BENCHMARK_RESULTS"
        
        echo ""
        echo "Quick benchmark results:"
        grep -E "Benchmark(GoStripComments|GoRealWorldCode|GoMemoryUsage|GoVeryLargeFiles)" "$NEW_BENCHMARK_RESULTS" | awk '
            { 
                if ($1 ~ /BenchmarkGoStripComments\/LargeCode/) {
                    printf "  Large code: %.2f ms/op, %.2f MB memory, %d allocs/op\n", 
                        $3/1000000, $5/1048576, $7
                } else if ($1 ~ /BenchmarkGoRealWorldCode/) {
                    printf "  Real world: %.2f ms/op, %.2f MB memory, %d allocs/op\n", 
                        $3/1000000, $5/1048576, $7
                } else if ($1 ~ /BenchmarkGoMemoryUsage/) {
                    printf "  Memory usage: %.2f ms/op, %.2f MB memory, %d allocs/op\n", 
                        $3/1000000, $5/1048576, $7
                } else if ($1 ~ /BenchmarkGoVeryLargeFiles\/1000Lines/) {
                    printf "  1000 lines: %.2f ms/op, %.2f MB memory, %d allocs/op\n", 
                        $3/1000000, $5/1048576, $7
                } else if ($1 ~ /BenchmarkGoVeryLargeFiles\/5000Lines/) {
                    printf "  5000 lines: %.2f ms/op, %.2f MB memory, %d allocs/op\n", 
                        $3/1000000, $5/1048576, $7
                }
            }
        ' "$NEW_BENCHMARK_RESULTS"
        ;;
        
    compare)
        COMPARE_WITH="$BASELINE_RESULTS_FILE"
        if [ ! -z "$2" ]; then
            COMPARE_WITH="$2"
        fi
        
        if [ ! -f "$COMPARE_WITH" ]; then
            echo "Error: Comparison file $COMPARE_WITH not found."
            echo "Run './benchmark.sh run' first to establish a baseline."
            exit 1
        fi
        
        run_full_benchmark_suite "$NEW_BENCHMARK_RESULTS"
        cp "$NEW_BENCHMARK_RESULTS" "$LATEST_RESULTS_FILE"
        
        append_benchmark_results_to_history "$NEW_BENCHMARK_RESULTS"
        
        echo ""
        compare_results "$COMPARE_WITH" "$NEW_BENCHMARK_RESULTS"
        ;;
        
    history)
        if [ -f "$BENCHMARK_HISTORY_FILE" ]; then
            cat "$BENCHMARK_HISTORY_FILE"
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