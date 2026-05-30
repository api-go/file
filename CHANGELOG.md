# Changelog

## [1.3.3] - 2026-05-30
- **新增**: 注册到 `jsmod`，将所有文件操作方法（Read/Write/Archive 等）注册为高危方法（unsafeList），确保安全沙箱隔离。

## [1.0.7] - 2026-05-06
- **设计哲学对齐**：全面废除 `Must` 前缀函数（`MustRead`, `MustReadBytes`, `MustReadLines`, `MustReadDir`, `MustGzip`, `MustGunzip`, `MustZip`, `MustUnzip`），改为配合 `go/cast` 的 `As` 函数消除摩擦。
- **内部优化**：重构 `memory.go` 以移除对废弃 `Must` 函数的内部依赖。

## [1.0.6] - 2026-05-01
- (同步版本号)

## [1.0.5] - 2026-05-05
- **性能优化**: 优化 `EnsureParentDir`，减少冗余的系统调用。
- **基础设施对齐**: `UnmarshalFile` 迁移至 `cast.To` 语义对齐（内部仍使用 `Convert` 确保指针更新）。
- **健壮性**: `RunCommand` 现在支持过滤空行，并提供更详细的执行错误信息。
- **命名规范**: `EnsureDirSuffix` 内部变量 `spe` 重命名为 `sep`。
- **测试增强**: 补全了 `Marshal/Unmarshal` 的 Benchmark 测试。

## [1.0.4] - 2026-05-01
...

- 修复了依赖引用问题（convert 模块函数命名不一致导致编译失败）。
