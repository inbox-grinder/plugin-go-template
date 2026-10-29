# Go Plugin Template

A minimal Inbox Grinder plugin written in Go that compiles to both a **native binary** (`local-binary` runtime) and a **WebAssembly module** (`wasm-component` runtime via TinyGo).

## Structure

```
go-template/
├── go.mod              — Go module declaration
├── manifest.json       — binary runtime (default)
├── manifest-wasm.json  — WASM runtime
├── build.sh            — builds both artifacts
├── plugin.go           — shared extraction logic (EmailInput → ExtractionResult)
├── main_binary.go      — binary entry (//go:build !wasm)
├── main_wasm.go        — WASM exports (//go:build wasm, alloc / process)
└── fixtures/
    ├── email-input.json
    └── extraction-result.json
```

## Quick start

### Binary

```sh
go build -o plugin .
# Test manually:
./plugin < fixtures/email-input.json | jq .
```

### WASM

Requires [TinyGo](https://tinygo.org/getting-started/install/) ≥ 0.35.

```sh
tinygo build -target=wasip1 -gc=leaking -no-debug -o plugin.wasm .
# Switch manifest to WASM:
cp manifest-wasm.json manifest.json
```

Or run `build.sh` to produce both artifacts at once.

> **Why `-gc=leaking`?**
> TinyGo's WASI target calls `proc_exit(0)` at the end of `_start` (the
> initialisation entry point). The conservative GC tears down its internal
> state at that point, making subsequent calls to exported functions fail.
> The leaking (bump-pointer) GC has no such cleanup, so exported functions
> remain callable after `_start` completes.

## Making your own plugin

1. **Copy this directory** and rename it (e.g., `my-plugin`).
2. **Edit `manifest.json`** — change `id`, `name`, `matches`, and `permissions`.
3. **Edit `plugin.go`** — implement your extraction logic in `run()`.
4. Build and drop the folder into the app.

## WASM ABI

The WASM host calls:

| Export | Signature | Purpose |
|--------|-----------|---------|
| `_start` | `() → ()` | TinyGo runtime init (called once before alloc/process) |
| `alloc` | `(i32) → i32` | Allocate `n` bytes; host writes input JSON there |
| `process` | `(i32, i32) → i64` | Run extraction; return `(output_ptr << 32) | output_len` |
| `memory` | — | Linear memory shared with the host |

`process` returns `0` on error (host treats this as plugin failure).

## Build tags

| File | Build tag | When used |
|------|-----------|-----------|
| `main_binary.go` | `!wasm` | `go build .` (native) |
| `main_wasm.go` | `wasm` | `tinygo build -target=wasip1 ...` |
| `plugin.go` | (none) | always included |

## Notes

- `main_wasm.go` uses a package-level `inputBuf` slice for the alloc buffer and `outputBuf` for the result. Both persist for the lifetime of the WASM module instance.
- The binary variant reads the full `EmailInput` JSON from stdin and writes `ExtractionResult` JSON to stdout. Non-zero exit code signals failure; diagnostic output goes to stderr.
- `encoding/json` is used in both runtimes. TinyGo 0.35 supports it fully with the `wasip1` target.
