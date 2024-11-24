#!/bin/bash

# Read input from stdin
while IFS= read -r line || [ -n "$line" ]; do
    # Process the input
    # Example: Convert input to uppercase
    output=$(echo -n "$line" | tr a-z A-Z)

    # Output the result to stdout
    echo -n "$output"
done

