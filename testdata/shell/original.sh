#!/bin/bash

# Shell test file with comprehensive comment scenarios

# Environment variables
export PATH="/usr/local/bin:$PATH"
HOME_DIR="/home/user"

# shellcheck disable=SC2034
UNUSED_VAR="This variable is unused but preserved"

# Function definition
function greet() {
    # Local variable
    local name="$1"
    echo "Hello, $name!"
}

# Conditional statements
if [ -f "/etc/passwd" ]; then
    # Reading file
    echo "File exists"
elif [ -d "/home" ]; then
    # Directory check
    echo "Directory exists"
else
    # Default case
    echo "Neither exists"
fi

# shellcheck source=./config.sh
source "./config.sh"

# Loop constructs
for i in {1..5}; do
    # Iteration comment
    echo "Iteration: $i"
done

# While loop
counter=0
while [ $counter -lt 3 ]; do
    # Increment counter
    ((counter++))
    echo "Counter: $counter"
done

# Case statement
case "$1" in
    start)
        # Start command
        echo "Starting service"
        ;;
    stop)
        # Stop command
        echo "Stopping service"
        ;;
    *)
        # Default case
        echo "Usage: $0 {start|stop}"
        ;;
esac

# String handling with comments
echo "This # is not a comment inside quotes"
echo 'Single quotes # also protect'
VAR='value # with hash' # This is a comment

# Command substitution
CURRENT_DATE=$(date) # Get current date
FILES=`ls -la` # Alternative command substitution

# Here document - comments inside should remain
cat << 'EOF'
This is a here document
# This should not be treated as a comment
Some text with # symbols
EOF

# Here string
grep "pattern" <<< "text with # symbol"

# Arrays (bash-specific)
ARRAY=("item1" "item2" "item3") # Array definition
echo "${ARRAY[0]}" # First element

# Arithmetic operations
RESULT=$((5 + 3)) # Addition
let "VALUE = 10 * 2" # Multiplication

# Process substitution (bash-specific)
diff <(ls dir1) <(ls dir2) # Compare directory listings

# Redirections with comments
echo "output" > file.txt # Redirect to file
echo "error" 2> error.log # Redirect stderr
echo "both" &> all.log # Redirect both stdout and stderr

# Background processes
sleep 10 & # Run in background
wait # Wait for background processes

# Regular expressions in conditions
if [[ "$USER" =~ ^[a-z]+$ ]]; then
    # User name validation
    echo "Valid username"
fi

# Multiple commands on one line
echo "first"; echo "second" # Semicolon separated
echo "third" && echo "fourth" # AND condition
echo "fifth" || echo "sixth" # OR condition

# Variable expansion
echo "Home: ${HOME:-/default}" # Parameter expansion with default
echo "Length: ${#HOME}" # String length

# Command line arguments
echo "Script name: $0" # Script name
echo "First arg: $1" # First argument
echo "All args: $@" # All arguments

# Exit status
echo "Last exit code: $?" # Previous command exit status

# Advanced variable operations
STRING="hello world"
echo "${STRING^^}" # Uppercase (bash 4+)
echo "${STRING%% *}" # Remove longest match from end

# Subshell
(
    # This runs in a subshell
    cd /tmp
    echo "Current dir: $(pwd)"
)

# Trap signals
trap 'echo "Exiting..."; exit 0' INT TERM # Signal handling

# Final comment before exit
exit 0 