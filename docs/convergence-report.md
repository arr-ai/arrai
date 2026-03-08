# Convergence Report

Standing invariants: all green.

## Movement

- 🎯T1: close → **achieved** (merged to frozen master as PR #88)
- 🎯T2: not started → **retired** (write-path allocs are structural, no real bottleneck)

## Gap Report

### 🎯T1 Frozen read-path allocations are zero  [high, weight: 9]
Gap: **achieved**

Merged to frozen master as PR #88 (5b2842e). All benchmarks confirmed 0 allocs
for concrete key types, 1 alloc for interface keys.

### 🎯T2 Frozen write-path allocations minimised  [medium, weight: 5]
Gap: **retired**

Benchmarking showed 3-18 allocs for Map.With and 1-7 for Set.With, scaling with
tree depth. These are structurally inherent to persistent data structures
(spine-copying) and already reasonable. No concrete arrai workload shows
write-path as a bottleneck.

## Recommendation

No active targets remain. The frozen performance work is complete — 🎯T1 achieved,
🎯T2 retired. Consider defining new targets based on current project priorities.

<!-- convergence-deps
evaluated: 2026-03-08T00:00:00Z
sha: 3766045

🎯T1:
  gap: achieved
  assessment: "Merged to frozen master as PR #88. Target archived."
  read: []

🎯T2:
  gap: retired
  assessment: "Write-path allocs structural. Retired after benchmarking."
  read: []
-->
