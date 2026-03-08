# Convergence Targets

## Active

(none)

## Retired

### 🎯T2 Frozen write-path allocations minimised  [medium, weight: 5]

Retired 2026-03-08. Benchmarking showed write-path allocations (3-18 for Map.With,
1-7 for Set.With) are structurally inherent to persistent data structures
(spine-copying). Already reasonable for tree depth. Read-path fix (🎯T1) addressed
the actual regression. No concrete workload shows write-path as a bottleneck.

## Achieved

### 🎯T4 Agent guide is discoverable from the repo root  [medium, weight: 5]

Achieved 2026-03-08. Root symlink `agents-guide.md` → `cmd/arrai/agents-guide.md`,
README updated with agent guide section.

### 🎯T3 arrai v0.333.0 released  [high, weight: 8]

Achieved 2026-03-08. Release created automatically by generate-tag workflow.
Tag `v0.333.0` on master, GitHub release with binaries for all platforms.

### 🎯T1 Frozen read-path allocations are zero  [high, weight: 9]

Achieved 2026-03-07. Merged to frozen master as PR #88 (5b2842e).
All benchmarks at 0 allocs for concrete key types, 1 alloc for interface keys.
