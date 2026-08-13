#!/usr/bin/env bash
#
# Regenerates the go-swagger client in ./go-openapiv2 from ./openapiv2/*.swagger.json.
#
# The generator version is pinned. go-swagger's output is not stable across
# releases, so an unpinned generator turns the drift check in CI into a coin
# toss and produces spurious diffs for whoever happens to have a different
# binary on PATH. If you upgrade the pin, regenerate and commit in the same
# change.
#
# Run ./gen-proto-go.sh first when the .proto has changed: this script consumes
# the swagger JSON that buf produces, so it can only be as current as that file.

set -euo pipefail

SWAGGER_VERSION="v0.33.0"
TARGET="go-openapiv2"
SPEC="openapiv2/restcol.swagger.json"

swagger_bin="$(command -v swagger || true)"
if [[ -z "${swagger_bin}" ]]; then
    echo "swagger not found on PATH. Install the pinned version:" >&2
    echo "  curl -sSL -o /usr/local/bin/swagger \\" >&2
    echo "    https://github.com/go-swagger/go-swagger/releases/download/${SWAGGER_VERSION}/swagger_linux_amd64" >&2
    echo "  chmod +x /usr/local/bin/swagger" >&2
    exit 1
fi

have="$(swagger version 2>/dev/null | awk '/^version:/ {print $2}')"
if [[ "${have}" != "${SWAGGER_VERSION}" ]]; then
    echo "swagger ${have:-unknown} on PATH, but this repo pins ${SWAGGER_VERSION}." >&2
    echo "Generated output differs between releases; install the pinned version." >&2
    exit 1
fi

# Remove the previous output rather than generating over it. go-swagger does
# not delete files it no longer emits, so generating in place leaves orphans
# behind - which is how client/collections outlived the tag it was named after.
rm -rf "${TARGET}"
mkdir -p "${TARGET}"

swagger generate client -f "${SPEC}" --target "${TARGET}"
