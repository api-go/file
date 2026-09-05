# Test Report: @go/file

## v1.5.8 验证
- 独立模块解析、全量测试与基准测试均通过。
- YAML 文件字段映射覆盖大小写与分隔符差异，并验证通过 `cast` 的统一解码路径。

## 📋 测试概览
- **测试时间**: 2026-09-05
- **测试环境**: darwin/arm64
- **Go 版本**: 1.25.0

## ✅ 功能测试 (Functional Tests)
| 场景 | 状态 | 描述 |
| :--- | :--- | :--- |
| `TestFileBasic` | PASS | 文件写入、移动、替换、删除及 `cast.As` 读取测试。 |
| `TestReadLines` | PASS | 文件按行读取验证。 |
| `TestReadDir` | PASS | 目录遍历读取验证。 |
| `TestArchive` | PASS | Tar/Gzip/Zip 归档与解压完整链路测试。 |
| `TestMemory` | PASS | 内存文件系统加载与检索验证。 |

## 🛡️ 鲁棒性防御 (Robustness)
- **摩擦消除**：移除 `Must` 系列 API，极大降低了业务逻辑中的错误处理心智负担。
- **内存虚拟化**：支持透明的内存文件系统操作，与磁盘文件 API 完全一致。
- **JS 错误调用栈**：JS 桥接层改用具名导出并使用 `jsmod.MakeError` 包裹错误，确保 JS 抛出异常时携带准确的 Go 运行时堆栈。

## ⚡ 性能基准 (Benchmarks)
| 函数 | 平均耗时 | 性能分析 |
| :--- | :--- | :--- |
| `MarshalFile` | **35610 ns/op** | 包含 IO 写入的序列化。 |
| `MarshalFilePretty` | **36006 ns/op** | 格式化文件序列化。 |
| `UnmarshalFile` | **11768 ns/op** | 文件反序列化。 |

## 覆盖率
- **整体覆盖率**: 35.5% of statements
