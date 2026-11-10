# 关于本项目

本项目完全由 AI 维护。代码源自 github.com/ssgo/u 的重构。

# @go/file

`@go/file` 是一个为“高性能、内存虚拟化、消除摩擦”设计的 IO 工具库。它支持原生文件系统操作、内存文件系统（支持 Gzip 压缩存储）、以及智能对象序列化（JSON/YAML）。

## 🎯 设计哲学

*   **内存文件系统 (Memory FS)**：支持将文件加载到内存中进行读写，极大提升了测试场景与高频小文件读写的性能。支持安全加载（基于 `go/safe`）与自动压缩存储。
*   **消除摩擦 (Frictionless)**：废除 `Must` 前缀函数，全面结合 `go/cast` 的 `As` 函数。
*   **智能序列化**：提供 `UnmarshalFile` 与 `MarshalFile`，自动处理路径创建与序列化格式；YAML 与 JSON 共享 `cast` 的字段命名和模糊匹配规则。

## 🛠 API Reference

### 基础文件操作 (Frictionless with cast.As)
- `func ReadBytes(filename string) ([]byte, error)`
- `func Read(filename string) (string, error)`
- `func ReadLines(filename string) ([]string, error)`
- `func WriteBytes(filename string, content []byte) error`
- `func Write(filename string, content string) error`
- `func Exists(filename string) bool`

### 压缩与归档
- `func Compress(data []byte, cType string) ([]byte, error)`
- `func Decompress(data []byte, cType string) ([]byte, error)`
- `func Archive(srcPath, destFile string) error`
- `func Extract(srcFile, destDir string, stripRoot bool) error`

### 内存文件系统 (Memory FS)
- `func AddFileToMemory(mf MemFile)`
- `func ReadFileFromMemory(name string) *MemFile`
- `func SafeLoadFileToMemory(filename string)`
- `func LoadFilesToMemoryFromB64(b64File *MemFileB64)`

## 📦 安装

```bash
go get github.com/api-go/file
```

## 💡 示例

```go
import (
    "github.com/api-go/file"
    "github.com/api-go/cast"
)

// 读取文件内容，消除错误摩擦
content := cast.As(file.Read("config.txt"))
```
