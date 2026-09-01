# Scale suite (S0)

Workloads that measure the large-data regime the 2026-08 overhaul never
saw. The 1000-app reconstruct scenario stays in `perf/reconstruct`; this
directory holds the programs that go bigger.

| Program | What it is | Default N |
|---|---|---|
| generated sysl protobuf → `vendor/run.arrai` | reconstruct at 10k apps | 10000 |
| `join.arrai` | join + `where` + `=>` over a million rows | 1000000 |

They skip under `-short` and unless `ARRAI_SCALE=1` is set, so default
`go test ./...` / CI Test stay on the 1000-app reconstruct. To run them:

```sh
ARRAI_SCALE=1 go test ./perf -count=1 -run 'TestScale'
```
