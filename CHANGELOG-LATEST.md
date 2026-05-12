## [1.0.7] - 2026-05-06
- **设计哲学对齐**：全面废除 `Must` 前缀函数（`MustRead`, `MustReadBytes`, `MustReadLines`, `MustReadDir`, `MustGzip`, `MustGunzip`, `MustZip`, `MustUnzip`），改为配合 `go/cast` 的 `As` 函数消除摩擦。
- **内部优化**：重构 `memory.go` 以移除对废弃 `Must` 函数的内部依赖。

