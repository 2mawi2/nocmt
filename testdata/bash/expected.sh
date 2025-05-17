#!/bin/bash
NAME="World"
COUNT=5

# shellcheck disable=SC2034
UNUSED_VAR="This is unused"

echo "Hello, $NAME!"

# shellcheck source=./utilities.sh
for i in $(seq 1 $COUNT); do
    echo "Iteration $i"
done

exit 0

echo "Hello"

echo "This is not a # comment"
echo 'This is not a # comment either'
echo "Hello"


echo "Code after comment block"

cat << EOF
This is a here document
# This is not a comment
EOF
echo "After heredoc"

NAME="John"
AGE=30
echo "$NAME is $AGE years old"

echo "Escaped \# symbol is not a comment"
echo "Regular # inside string also not a comment"
echo 'Hash # in single quotes'
echo "Hash # in double quotes"
VAR='value with # symbol'

# shellcheck disable=SC2034
VAR_DIRECTIVE_1="unused variable for directive test"
# shellcheck disable=SC2154
# shellcheck source=./lib.sh
# shellcheck shell=bash
echo "Testing directives" 