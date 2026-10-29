#!/usr/bin/env bash
# Build the Go plugin as both a native binary and a WASM module.
# Run from the go-template directory.
set -euo pipefail

cd "$(dirname "$0")"

echo "==> Building binary (local-binary runtime)..."
go build -o plugin .
echo "    plugin ($(du -sh plugin | cut -f1))"

echo "==> Building WASM (wasm-component runtime)..."
# Requires: brew install tinygo  (or https://tinygo.org/getting-started/)
# -gc=leaking uses a bump allocator so the Go runtime stays valid after
# _start completes; the conservative GC tears down state on proc_exit(0).
tinygo build -target=wasip1 -gc=leaking -no-debug -o plugin.wasm .
echo "    plugin.wasm ($(du -sh plugin.wasm | cut -f1))"

echo ""
echo "Done. To switch between runtimes, copy the matching manifest:"
echo "  cp manifest.json manifest-binary.json      # save current (local-binary)"
echo "  cp manifest-wasm.json manifest.json        # switch to wasm-component"
