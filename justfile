default:
    @just --list

build:
    go build -o nocmt

test:
    go test ./...

lint:
    golangci-lint run ./...

run *args:
    go run main.go {{args}}

clean:
    rm -f nocmt 

bench *args:
    ./benchmark.sh {{args}} 

release branch="release":
    #!/usr/bin/env bash
    set -e
    echo "Starting release process to branch: {{branch}}"
    
    # Check if we're on main branch
    current_branch=$(git branch --show-current)
    if [ "$current_branch" != "main" ]; then
        echo "Error: Must be on main branch to create a release"
        exit 1
    fi
    
    # Pull latest changes from main
    echo "Pulling latest changes from main..."
    git pull origin main
    
    # Check for staged changes
    if [ -n "$(git status --porcelain)" ]; then
        echo "Error: Working directory is not clean. Commit or stash changes before creating a release."
        exit 1
    fi
    
    # Create or checkout the release branch
    if git show-ref --verify --quiet refs/heads/{{branch}}; then
        echo "Checking out existing {{branch}} branch..."
        git checkout {{branch}}
    else
        echo "Creating new {{branch}} branch..."
        git checkout -b {{branch}}
    fi
    
    # Merge changes from main
    echo "Merging changes from main..."
    git merge main
    
    # Push the release branch
    echo "Pushing {{branch}} branch to remote..."
    git push origin {{branch}}
    
    # Switch back to main
    echo "Switching back to main branch..."
    git checkout main
    
    echo "Release process completed successfully!" 