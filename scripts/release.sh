#!/usr/bin/env bash
set -e
VERSION=$(git describe --tags --abbrev=0)
git tag -a "$VERSION" -m "Release $VERSION"
git push --tags