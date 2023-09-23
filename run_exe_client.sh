#!/bin/bash

# Number of execute processes to start
execute_num=10
shard_num=2

# Array to store the process IDs
pids=()

# Start order processes for each shard in the background
for ((shard = 1; shard <= shard_num; shard++)); do
    for ((i = 1; i <= execute_num; i++)); do
        ./execute $shard  shard_num 10 20 &
        pids+=($!)  # Store the process ID in the array
    done
done

# Wait for all order processes to complete
for pid in "${pids[@]}"; do
    wait $pid
done

echo "All processes have completed."
