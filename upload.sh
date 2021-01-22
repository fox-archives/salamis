#!/bin/bash

set -euo pipefail
source .env

owner="eankeen"
repo="salamis"
tag="v0.2.3"
file="salamis.deb"

repoData="$(
	curl \
		--silent \
		--show-error \
		--header "Authorization: token $GITHUB_TOKEN" \
		"https://api.github.com/repos/eankeen/salamis/releases/tags/$tag"
)"

releaseId="$(jq '.id' <<< $repoData)"
: ${releaseId:?"releaseId must be valid. Exiting"}

curl \
	--header "Authorization: token $GITHUB_TOKEN" \
	--header "Content-Type: application/octet-stream" \
	--data-binary @"$file" \
	"https://uploads.github.com/repos/$owner/$repo/releases/$releaseId/assets?name=$file"
