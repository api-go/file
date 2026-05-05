# go/file (v1.0.4)

`go/file` 是一个为 Go 语言项目设计的轻量级、高性能且具备内存安全特性的文件系统操作与资源管理库。它封装了常见的文件 IO、对象序列化、归档压缩及内存文件映射功能，旨在通过统一的 API 简化开发工作，并提供内置的安全闭环保障。

## 核心特性

*   **内存文件系统 (Memory File System)**: 支持将文件或资源加载到内存，并可选进行压缩与 `SafeBuf` 加密存储。这非常适合处理嵌入式资源、敏感配置或高并发读取的场景。
*   **智能序列化**: 自动处理 YAML 与 JSON 格式判别，并集成了强大的 `cast.To` 智能字段映射引擎，轻松解决复杂结构体与配置文件之间的映射问题。
*   **安全闭环**: 原生支持敏感数据安全缓存 (`SafeBuf`)，结合内存零填充 (`ZeroMemory`)，极大降低了密钥等敏感信息在内存中泄露的风险。
*   **全功能 IO**: 提供从基础 IO 到递归目录操作、文件原地替换及高效压缩归档的一站式解决方案。

## 🤖 开发与 AI 指导 (Developer & AI Guidelines)

1.  **内存安全优先**: 处理敏感文件时，优先使用 `SafeLoad` 系列，并通过 `SafeBuf` 接口操作。
2.  **高性能 IO**: 在文件 IO 高并发场景下，默认采用 `WriteBytes` 以保证写入性能。
3.  **语义一致性**: 序列化/反序列化（`Unmarshal/Marshal`）通过 `cast.To` 处理自动映射，无需关心字段名不一致问题。
4.  **资源管理**: 任何归档提取或 IO 流操作结束后，必须确保对应资源被 `Close()`。

## 快速入门 (Quick Start)

### 1. 高性能文件读写
无需手动处理 `os.OpenFile` 和 `defer Close`，支持物理文件与内存模拟文件透明切换。

```go
import "apigo.cc/go/file"

// 简单读取与写入
content := file.MustRead("config.txt")
file.Write("log.txt", "operation success")

// 按行读取
lines := file.MustReadLines("list.txt")
```

### 2. 配置文件自动序列化
自动识别 YAML/JSON，并支持 snake_case 到 CamelCase 的智能字段映射。

```go
type AppConfig struct {
    UserName string `yaml:"user_name"`
    Port     int    `json:"port"`
}

var cfg AppConfig
// 自动读取、判别并完成映射
err := file.UnmarshalFile("config.yaml", &cfg)
```

### 3. 原子配置补丁更新
先读取文件，将补丁部分通过智能映射覆盖并回写，适合增量修改配置文件。

```go
// 自动读取 -> 智能 Patch -> 写回
file.PatchFile("data.json", map[string]any{"debug": true})
```

## 🛠 API Reference

### 文件操作 (Filesystem)
- `func Exists(filename string) bool`: 检查文件在磁盘或内存中是否存在。
- `func ReadBytes(filename string) ([]byte, error)`: 从磁盘或内存读取原始字节。
- `func MustReadBytes(filename string) []byte`: `ReadBytes` 的封装，忽略错误，直接返回字节。
- `func Read(filename string) (string, error)`: 读取文件并返回 UTF-8 字符串。
- `func MustRead(filename string) string`: 忽略错误的 `Read` 实现。
- `func ReadLines(filename string) ([]string, error)`: 按行读取文件内容。
- `func MustReadLines(filename string) []string`: 忽略错误的 `ReadLines` 实现。
- `func WriteBytes(filename string, content []byte) error`: 将字节写入文件，会自动处理父目录创建，高性能 IO。
- `func Write(filename string, content string) error`: `WriteBytes` 的字符串封装，处理常规文本写入。
- `func Copy(from, to string) error`: 文件或目录的递归复制。
- `func CopyToFile(from io.Reader, to string) error`: 从流拷贝数据到文件。
- `func Remove(path string) error`: 递归删除文件或目录。
- `func Move(src, dst string) error`: 重命名或移动文件。
- `func Replace(filename, old, new string) error`: 在文件中对文本进行原地批量替换。
- `func Search(dir, pattern string) []string`: 在目录下根据通配符模式匹配文件路径。
- `func ReadDir(filename string) ([]FileInfo, error)`: 读取目录信息，包含内存与磁盘兼容层。
- `func MustReadDir(filename string) []FileInfo`: 忽略错误的 `ReadDir` 实现。

### 对象序列化 (Object/Unmarshal/Marshal)
- `func UnmarshalFile(filename string, to any) error`: 自动识别 YAML/JSON，并进行智能字段映射（支持 snake_case 到 CamelCase 映射）。
- `func MarshalFile(filename string, data any) error`: 根据后缀自动判定为 JSON 或 YAML 并写入文件。
- `func MarshalFilePretty(filename string, data any) error`: 同 `MarshalFile`，但输出带缩进的可读格式。
- `func PatchFile(filename string, patch any) error`: 先读取文件，再将 patch 部分通过智能映射覆盖并回写，适合增量修改配置文件。

### 归档与压缩 (Archive/Compress)
- `func Compress(data []byte, cType string) ([]byte, error)`: 支持 gzip/zlib 压缩。
- `func Decompress(data []byte, cType string) ([]byte, error)`: 支持 gzip/zlib 解压。
- `func MustGzip(data []byte) []byte`: 强制压缩（仅 gzip）。
- `func MustGunzip(data []byte) []byte`: 强制解压（仅 gzip）。
- `func Archive(srcPath, destFile string) error`: 将目录压缩为 .zip 或 .tar.gz。
- `func Extract(srcFile, destDir string, stripRoot bool) error`: 自动识别格式并解压，`stripRoot` 可去除归档内的顶层目录。

### 内存文件系统 (Memory)
- `func AddFileToMemory(mf MemFile)`: 直接向内存中插入一个 `MemFile` 对象。
- `func ReadFileFromMemory(name string) *MemFile`: 内存文件读取。
- `func LoadFileToMemory(filename string)`: 物理文件加载到内存。
- `func SafeLoadFileToMemory(filename string)`: 加载到内存并使用 `SafeBuf` 加密存储。
- `func LoadFileToMemoryWithCompress(filename string)`: 加载并压缩存储。
- `func SafeLoadFileToMemoryWithCompress(filename string)`: 加载、压缩并加密存储。
- `func LoadFileToB64(filename string) *MemFileB64`: 将文件转换为 B64 结构（常用于嵌入式资源）。
- `func LoadFilesToMemoryFromB64(b64File *MemFileB64)`: 从 B64 结构恢复文件到内存。
- `func LoadFilesToMemoryFromJson(jsonFiles string)`: 从 JSON 串格式的 B64 资源恢复内存文件。

## 许可证

本项目基于 MIT 许可证开源。
