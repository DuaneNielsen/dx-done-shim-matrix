// A stand-in for the real dx-done binary.
//
// This repo verifies the launcher SHIM (bin/dx-done) — that it resolves the
// correct per-platform file by `uname -s` and execs it. It deliberately does
// NOT test the real dx-done CLI (a separate workflow already smoke-tests that
// the real binaries boot). So all this stub needs to be is a genuine native
// executable that proves the exec actually happened: it prints the one line
// the smoke test greps for. Compiled per-runner (Go is preinstalled on the
// ubuntu/macos/windows GitHub-hosted images), so on Windows it is a real PE
// `.exe` — which is what faithfully exercises the MSYS `exec` question.
package main

import "fmt"

func main() {
	fmt.Println("dx-done v0.0.0-stub")
}
