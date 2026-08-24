# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
?   	example.com/childfitness/cmd/fitness	[no test files]
--- FAIL: TestChildFitnessImportIsolatesBadRow (0.01s)
    import_test.go:28: expected partial batch, got complete
FAIL
FAIL	example.com/childfitness	0.011s
ok  	example.com/childfitness/internal/capture	0.001s
ok  	example.com/childfitness/internal/cli	0.010s
ok  	example.com/childfitness/internal/domain	0.001s
ok  	example.com/childfitness/internal/privacy	0.001s
ok  	example.com/childfitness/internal/report	0.001s
ok  	example.com/childfitness/internal/service	0.015s
ok  	example.com/childfitness/internal/store	0.007s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/fitness): exit `0`
- Frontend build (web): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/fitness): exit `0`
- Frontend build (web): exit `0`
