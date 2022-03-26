#!/bin/bash
APP_NAME=${1:?"app name"}


if BUILD_GIT_REVISION=$(git rev-parse HEAD 2> /dev/null); then
  if [[ -n "$(git status --porcelain 2>/dev/null)" ]]; then
    BUILD_GIT_REVISION=${BUILD_GIT_REVISION}"-dirty"
  fi
else
  BUILD_GIT_REVISION=unknown
fi

# Check for local changes
if git diff-index --quiet HEAD --; then
  tree_status="Clean"
else
  tree_status="Modified"
fi

# XXX This needs to be updated to accomodate tags added after building, rather than prior to builds
RELEASE_TAG=$(git describe --match 'v[0-9]*\.[0-9]*\.[0-9]*' --exact-match 2> /dev/null || echo "")

# security wanted VERSION='unknown'
VERSION="${BUILD_GIT_REVISION}"
if [[ -n "${RELEASE_TAG}" ]]; then
  VERSION="${RELEASE_TAG}"
fi

GIT_DESCRIBE_TAG=$(git describe --tags)

# used by scripts/build/gobuild.sh
echo "github.com/hkloudou/xlib/xruntime._appName=${APP_NAME}"
echo "github.com/hkloudou/xlib/xruntime._buildVersion=${VERSION}"
echo "github.com/hkloudou/xlib/xruntime._buildAppVersion=${BUILD_GIT_REVISION}"
echo "github.com/hkloudou/xlib/xruntime._buildStatus=${tree_status}"
echo "github.com/hkloudou/xlib/xruntime._buildTag=${GIT_DESCRIBE_TAG}"
echo "github.com/hkloudou/xlib/xruntime._buildUser=$(whoami)"
echo "github.com/hkloudou/xlib/xruntime._buildHost=$(hostname -f)"
echo "github.com/hkloudou/xlib/xruntime._buildTime=$(date +%s)"