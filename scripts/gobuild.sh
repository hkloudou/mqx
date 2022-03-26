#!/bin/bash

# WARNING: DO NOT EDIT, THIS FILE IS PROBABLY A COPY
#
# The original version of this file is located in the https://github.com/istio/common-files repo.
# If you're looking at this file in a different repo and want to make a change, please go to the
# common-files repo, make the change there and check it in. Then come back to this repo and run
# "make update-common".

# Copyright Istio Authors. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This script builds and version stamps the output
VERBOSE=${VERBOSE:-"0"}
V=""
if [[ "${VERBOSE}" == "1" ]];then
    V="-x"
    set -x
fi

SCRIPTPATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

APP_NAME=${1:?"app name"}
OUT=${2:?"output path"}

shift
set -e
# GOOS=${GOOS:-darwin}
# GOARCH=${GOARCH:-amd64}
GOBINARY=${GOBINARY:-go}
# GOPKG="$(go env GOPATH)/pkg"
BUILDINFO=${BUILDINFO:-""}
# STATIC=${STATIC:-1}
LDFLAGS=${LDFLAGS:-"-s -w"}
GOBUILDFLAGS=${GOBUILDFLAGS:-""}
GCFLAGS=${GCFLAGS:-}

# Split GOBUILDFLAGS by spaces into an array called GOBUILDFLAGS_ARRAY.
IFS=' ' read -r -a GOBUILDFLAGS_ARRAY <<< "$GOBUILDFLAGS"

GCFLAGS=${GCFLAGS:-}
# export CGO_ENABLED=0

# gather buildinfo if not already provided
# For a release build BUILDINFO should be produced
# at the beginning of the build and used throughout
if [[ -z ${BUILDINFO} ]];then
    BUILDINFO=$(mktemp)
    sh "${SCRIPTPATH}/report_build_info.sh"  ${APP_NAME}> "${BUILDINFO}"
fi

# BUILD LD_EXTRAFLAGS
LD_EXTRAFLAGS=""

while read -r line; do
    LD_EXTRAFLAGS="${LD_EXTRAFLAGS} -X ${line}"
done < "${BUILDINFO}"

OPTIMIZATION_FLAGS="-trimpath"
if [ "${DEBUG}" == "1" ]; then
    OPTIMIZATION_FLAGS=""
fi

#goxc -d=${OUT} -pv=${APP_VERSION} -bc='linux,windows,darwin' -build-ldflags="${LD_EXTRAFLAGS}"
echo time ${GOBINARY} build ${OPTIMIZATION_FLAGS} ${V} ${GOBUILDFLAGS} ${GCFLAGS:+-gcflags "${GCFLAGS}"} -o ${OUT} \
       -ldflags "${LDFLAGS} ${LD_EXTRAFLAGS}"
time ${GOBINARY} build ${OPTIMIZATION_FLAGS} ${V} ${GOBUILDFLAGS} ${GCFLAGS:+-gcflags "${GCFLAGS}"} -o ${OUT} \
       -ldflags "${LDFLAGS} ${LD_EXTRAFLAGS}"