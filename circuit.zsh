#!/bin/zsh

# Define the repository path
REPO_PATH="/home/zk/pkgs/factory"

# Function to get the latest version tag from Git
get_latest_version() {
  local latest_tag
  latest_tag=$(git describe --tags $(git rev-list --tags --max-count=1) 2>/dev/null)
  
  if [ -z "$latest_tag" ]; then
    echo "v0.0.0" # Default to v0.0.0 if no tags exist
  else
    echo "$latest_tag"
  fi
}

# Function to increment the version (patch version by default)
increment_version() {
  local current_version=$1
  local major minor patch
  IFS='.' read -r major minor patch <<< "${current_version#v}"
  patch=$((patch + 1))
  echo "v$major.$minor.$patch"
}

# Function to stage, commit, and push changes
commit_and_push() {
  local commit_message=$1
  git add .
  git commit -m "$commit_message" || { echo "Commit failed. Exiting."; exit 1; }
  git push || { echo "Push failed. Exiting."; exit 1; }
}

# Function to create a GitHub release
create_release() {
  local version=$1
  gh release create "$version" --title "$version" --notes "Release $version" || {
    echo "Release creation failed. Exiting."
    exit 1
  }
  echo "Release $version created successfully!"
}

# Main script logic
main() {
  local commit_message=$1

  # Navigate to the repository
  if [ ! -d "$REPO_PATH" ]; then
    echo "Repository path $REPO_PATH does not exist. Exiting."
    exit 1
  fi
  cd "$REPO_PATH" || { echo "Failed to navigate to $REPO_PATH. Exiting."; exit 1; }

  # Get the latest version and increment it
  local latest_version new_version
  latest_version=$(get_latest_version)
  echo "Latest version: $latest_version"
  new_version=$(increment_version "$latest_version")
  echo "New version: $new_version"

  # Commit changes and push
  commit_and_push "$commit_message"

  # Create a new release
  create_release "$new_version"
}

# Check if a commit message is provided
if [ -z "$1" ]; then
  echo "Usage: ./circuit.zsh <commit-message>"
  exit 1
fi

# Run the main function
main "$1"