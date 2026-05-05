# Test Report: @go/file

## 📋 测试概览
- **测试时间**: 2026-05-06
- **测试环境**: darwin/amd64
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

## ⚡ 性能基准 (Benchmarks)
| 函数 | 平均耗时 | 性能分析 |
| :--- | :--- | :--- |
| `MarshalFile` | **115787 ns/op** | 包含 IO 写入的高性能序列化。 |
| `UnmarshalFile` | **24682 ns/op** | 高效的反序列化。 |
