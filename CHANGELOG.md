# Changelog

## [1.0.5] - 2026-05-05
- **性能优化**: 优化 `EnsureParentDir`，减少冗余的系统调用。
- **基础设施对齐**: `UnmarshalFile` 迁移至 `cast.To` 语义对齐（内部仍使用 `Convert` 确保指针更新）。
- **健壮性**: `RunCommand` 现在支持过滤空行，并提供更详细的执行错误信息。
- **命名规范**: `EnsureDirSuffix` 内部变量 `spe` 重命名为 `sep`。
- **测试增强**: 补全了 `Marshal/Unmarshal` 的 Benchmark 测试。

## [1.0.4] - 2026-05-01
...

- 修复了依赖引用问题（convert 模块函数命名不一致导致编译失败）。
