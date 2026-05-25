#!/bin/sh

# Preview the release notes that would be generated if we were to create a
# release now.

gh api -X POST /repos/jesseduffield/lazygit/releases/generate-notes \
  -f tag_name=v0.99.0 \
  -f target_commitish=master \
  -q .body | code -
