# 测试报告 (TEST.md)

本文件记录了 `@go/file` 的测试覆盖场景、Benchmark 结果及系统鲁棒性测试情况。

## 1. 测试场景覆盖

### 文件系统操作 (TestFileSystemOps)
- **Copy**: 测试了文件从源到目的地的递归复制，验证了文件内容的一致性。
- **Move**: 验证了文件移动操作，确保源文件被删除且目标文件存在。
- **Replace**: 验证了在文件中进行原地文本替换的准确性。
- **Remove**: 验证了递归删除功能的鲁棒性。

### 目录与检索 (TestSearchAndDir)
- **Search**: 测试了在嵌套目录结构中通过通配符 `*.txt` 检索文件的能力。
- **ReadDir**: 验证了目录读取功能，支持内存模拟环境与磁盘真实环境的一致性。

### 内存文件系统 (TestMemoryDeep & TestMemoryLoadAPI)
- **递归加载**: 验证了 `LoadFileToMemory` 对多级子目录的递归处理。
- **B64 & 压缩链路**: 测试了文件 -> B64 序列化 -> JSON 存储 -> 内存恢复的全链路完整性。
- **SafeBuf 安全链路**: 验证了数据被 `SafeBuf` 加密存储后，解密内容的正确性，以及 `Close()` 对内存的擦除能力。
- **并发写入**: 验证了 `MarshalFile` 在高并发写入场景下的数据完整性。

### 序列化与智能映射 (TestUnmarshal & TestMarshal)
- **智能映射 (SmartMapping)**: 验证了 `UnmarshalFile` 能将 `user_name` (snake_case) 自动映射到结构体的 `UserName` (CamelCase) 字段，无需手动配置。
- **pretty 输出**: 验证了 `MarshalFilePretty` 输出带格式的 JSON/YAML。

## 2. 性能基准测试 (Benchmark)

在 Intel(R) Core(TM) i9-9980HK CPU @ 2.40GHz 环境下：

| Benchmark | ops | 耗时 (ns/op) |
| :--- | :--- | :--- |
| MarshalFile/Normal | 10000 | 111,770 |
| MarshalFile/Pretty | 10000 | 109,581 |
| UnmarshalFile/Normal | 52528 | 22,172 |
| UnmarshalFile/X | 46568 | 24,833 |

## 3. 系统鲁棒性测试
- **不存在的文件**: 验证 `UnmarshalFile` 返回 `os.IsNotExist` 类型错误。
- **非法格式**: 验证输入非法格式文件（如 `{invalid}`）时，系统正确返回解析错误，拒绝执行，严禁 Panic。
