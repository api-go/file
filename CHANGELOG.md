# Changelog

## v1.5.7 (2026-09-05)
- **YAML 统一**: `UnmarshalFile`、`MarshalFile` 和 `MarshalFilePretty` 改用 `cast` 的 YAML 编解码能力，移除对 yaml.v3 的直接调用。
- **字段映射**: YAML 文件到 Struct 的映射现在与 JSON 一致，忽略键名大小写及分隔符差异。
- **依赖更新**: 对齐 `cast v1.5.6` 及其他基础模块当前版本。

## v1.5.6 (2026-07-18)
- **依赖校验修复**: 重新整理模块图并补齐 checksum。

## v1.5.5 (2026-06-21)
- **JS 对齐**: 重构 JS 导出为具名函数，并引入 `jsmod.MakeError` 动态包装错误以获取调用栈。
- **依赖更新**: 升级依赖 `jsmod` 至 `v1.5.3`，`cast` 至 `v1.5.3`，`rand` 至 `v1.5.3`，`encoding` 至 `v1.5.4`，`safe` 至 `v1.5.2`。

## v1.5.4 (2026-06-11)
- **版本对齐**: 基础设施全局对齐 v1.5.4。

## v1.5.3 (2026-06-09)
- **版本对齐**: 基础设施全局对齐 v1.5.3。

## v1.5.2 (2026-06-08)
- **JS 对齐**: 将所有注册到 `jsmod` 的导出方法名统一为 PascalCase（如 `Read`, `Write`, `GetFileInfo`），消除 JS 环境下的调用摩擦。

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
