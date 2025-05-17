#!/bin/bash
# This is a basic shell script with comments

# Configuration
NAME="World"
COUNT=5

# shellcheck disable=SC2034
UNUSED_VAR="This is unused"

# Regular comment
echo "Hello, $NAME!"

# shellcheck source=./utilities.sh
# Loop comment
for i in $(seq 1 $COUNT); do
    # Print iteration
    echo "Iteration $i"
done

# Final comment
exit 0

# Simple line comment
echo "Hello"  # End of line comment

# Comments inside string literals
echo "This is not a # comment"
echo 'This is not a # comment either'
echo "Hello"  # This is a comment that should be removed

# Empty comment lines
#
# 
#    

# Multiple adjacent comment lines
# First comment
# Second comment
# Third comment
echo "Code after comment block"

# Heredoc - comments inside should remain
cat << EOF
This is a here document
# This is not a comment
EOF
echo "After heredoc"

# Variable assignments with comments
NAME="John" # User name
AGE=30 # User age
echo "$NAME is $AGE years old"

# String handling
echo "Escaped \# symbol is not a comment"
echo "Regular # inside string also not a comment"
echo 'Hash # in single quotes'
echo "Hash # in double quotes"
VAR='value with # symbol' # This is a comment

# --- Directives testing section ---
# Regular comment
# shellcheck disable=SC2034
VAR_DIRECTIVE_1="unused variable for directive test"
# shellcheck disable=SC2154
# shellcheck source=./lib.sh
# shellcheck shell=bash
echo "Testing directives" 