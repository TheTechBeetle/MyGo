#!/bin/sh
set -eu

project_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
web_dir="$project_dir/web"
dist_dir="$web_dir/dist"
goroot=$(go env GOROOT)
wasm_exec="$goroot/lib/wasm/wasm_exec.js"

if [ ! -f "$wasm_exec" ]; then
	wasm_exec="$goroot/misc/wasm/wasm_exec.js"
fi

if [ ! -f "$wasm_exec" ]; then
	printf 'Could not find wasm_exec.js under %s\n' "$goroot" >&2
	exit 1
fi

mkdir -p "$dist_dir"
cp "$web_dir/index.html" "$dist_dir/index.html"
cp "$wasm_exec" "$dist_dir/wasm_exec.js"
(cd "$project_dir" && GOOS=js GOARCH=wasm go build -o "$dist_dir/pongo2.wasm" .)

printf 'Web build created at %s\n' "$dist_dir"
