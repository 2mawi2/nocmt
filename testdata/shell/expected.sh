#!/bin/bash

export PATH="/usr/local/bin:$PATH"
HOME_DIR="/home/user"

# shellcheck disable=SC2034
UNUSED_VAR="This variable is unused but preserved"

function greet() {
    local name="$1"
    echo "Hello, $name!"
}

if [ -f "/etc/passwd" ]; then
    echo "File exists"
elif [ -d "/home" ]; then
    echo "Directory exists"
else
    echo "Neither exists"
fi

# shellcheck source=./config.sh
source "./config.sh"

for i in {1..5}; do
    echo "Iteration: $i"
done

counter=0
while [ $counter -lt 3 ]; do
    ((counter++))
    echo "Counter: $counter"
done

case "$1" in
    start)
        echo "Starting service"
        ;;
    stop)
        echo "Stopping service"
        ;;
    *)
        echo "Usage: $0 {start|stop}"
        ;;
esac

echo "This # is not a comment inside quotes"
echo 'Single quotes # also protect'
VAR='value # with hash'

CURRENT_DATE=$(date)
FILES=`ls -la`

cat << 'EOF'
This is a here document
# This should not be treated as a comment
Some text with # symbols
EOF

grep "pattern" <<< "text with # symbol"

ARRAY=("item1" "item2" "item3")
echo "${ARRAY[0]}"

RESULT=$((5 + 3))
let "VALUE = 10 * 2"

diff <(ls dir1) <(ls dir2)

echo "output" > file.txt
echo "error" 2> error.log
echo "both" &> all.log

sleep 10 &
wait

if [[ "$USER" =~ ^[a-z]+$ ]]; then
    echo "Valid username"
fi

echo "first"; echo "second"
echo "third" && echo "fourth"
echo "fifth" || echo "sixth"

echo "Home: ${HOME:-/default}"
echo "Length: ${#HOME}"

echo "Script name: $0"
echo "First arg: $1"
echo "All args: $@"

echo "Last exit code: $?"

STRING="hello world"
echo "${STRING^^}"
echo "${STRING%% *}"

(
    cd /tmp
    echo "Current dir: $(pwd)"
)

trap 'echo "Exiting..."; exit 0' INT TERM

exit 0 