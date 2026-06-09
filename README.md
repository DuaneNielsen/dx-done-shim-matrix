# dx-done-shim-matrix

A self-contained, public CI rig that answers one narrow question:

> **Does the `claude-dx-done` launcher shim actually resolve the right
> per-platform binary and exec it on Linux, macOS, and Windows?**

The shim (`bin/dx-done`) is a `#!/usr/bin/env bash` script that picks a binary
by `uname -s` and execs it. Whether that survives **Windows** process-spawning
is the scariest unknown for the plugin's cross-platform rollout — and the reason
this repo exists. GitHub-hosted public runners are free, real clean machines for
all three OSes, so this is where the question gets answered on contact.

## What it tests (and what it doesn't)

- **Tests:** the **shim mechanics** — PATH resolution → shebang → `uname -s`
  dispatch → `exec` of the per-platform file. The shim is invoked the way Claude
  Code invokes it: the `bin/` dir is prepended to `PATH` and `dx-done` is called
  by name.
- **Does NOT test:** the real dx-done CLI. The binary the shim resolves to is a
  tiny **stub** (`stub/main.go`) that just prints `dx-done v0.0.0-stub`. The stub
  is compiled per-runner, so on Windows it is a genuine native PE `.exe` — which
  is what faithfully exercises the MSYS `exec` question. (A separate workflow
  already smoke-tests that the *real* binaries boot; this repo only adds the shim
  layer that one skips.)

The `bin/dx-done` here is a **verbatim copy** of the canonical shim at
`projects/claude-dx-done/plugin/bin/dx-done` in the wild-west monorepo. It is not
forked — equality is checked at author time.

## How to read a run

Three jobs, one per OS:

| Job | Runner | What green means |
|---|---|---|
| `linux` | `ubuntu-latest` | shim resolved `dx-done-linux-x64` and exec'd it via bash |
| `macos` | `macos-latest` (arm64) | shim resolved `dx-done-macos-arm64` and exec'd it via bash |
| `windows` | `windows-latest` | **Path A** (Git Bash) resolved `dx-done-windows-x64.exe` and exec'd it |

The **Windows job runs two paths** (see the job summary for the recorded outcomes):

- **Path A — Git Bash spawn** *(gates the job)*: how Claude Code's Bash tool
  invokes commands on Windows. **Expected: PASS.**
- **Path B — native non-bash spawn** *(characterized, not gated)*: PowerShell
  calling the extension-less shim directly, the way a `cmd`/`pwsh`-only user
  would. **Expected: FAIL** — Windows cannot natively run a `#!/usr/bin/env bash`
  file. This failure is the **informative answer**, not a test failure;
  `continue-on-error` keeps the job green.

**The pass/fail gate is the bash-path mechanics on all three OSes.** A "Path B
fails" result is a successful, informative outcome — it tells us whether the
pilot needs a `.cmd`/`.ps1`/native Windows launcher or whether the Git-Bash path
(which Claude Code uses) is sufficient.

## Verdict

<!-- Filled in after the first green run. -->

- **Linux (bash):** _pending run_
- **macOS arm64 (bash):** _pending run_
- **Windows Path A (Git Bash):** _pending run_
- **Windows Path B (native spawn):** _pending run_
- **Launcher-shape answer:** _pending run — is the bash shim viable across all
  three OSes as invoked by Claude Code? Does Windows need a non-bash launcher?_

## Context

- Parent tracking issue: `IMS/wild-west#94` (the Mac/PC/Linux install matrix).
- This repo covers the **public-runner shim-mechanics** half. The **on-corp
  marketplace-install → shim** half (which public runners can't reach, being
  off-corp) is tracked separately in `IMS/wild-west#185`.
- Binary *delivery* to end users (fetch-on-first-run + a public binary host) is
  `IMS/wild-west#90` — shaped by this repo's launcher-shape verdict.
