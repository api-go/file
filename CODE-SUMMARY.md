### file > object.go
```go



var fileLocksLock = sync.Mutex{}
var fileLocks = map[string]*sync.Mutex{}

func getLock(filename string) *sync.Mutex


// UnmarshalFile 从文件读取数据并映射到 to，自动处理 YAML/JSON 格式判别及智能字段映射。
func UnmarshalFile(filename string, to any) error


// MarshalFile 将数据序列化并写入文件，支持 YAML/JSON 自动判别。
func MarshalFile(filename string, data any) error


// MarshalFilePretty 类似 MarshalFile，但会输出带缩进的格式化内容。
func MarshalFilePretty(filename string, data any) error


func marshalFile(filename string, data any, indent bool) error


func PatchFile(filename string, patch any) error

```

### file > memory.go
```go



type MemFile struct {
	Name       string
	AbsName    string
	ModTime    time.Time
	IsDir      bool
	Compressed bool
	Size       int64
	Data       []byte
	SafeData   *safe.SafeBuf
}

type MemFileB64 struct {
	Name       string
	ModTime    time.Time
	IsDir      bool
	DataB64    []byte
	Compressed bool
	Size       int64
	Children   []MemFileB64
}

var (
	memFiles      = make(map[string]*MemFile)
	memFilesByDir = make(map[string][]MemFile)
	memFilesLock  sync.RWMutex
)

func (mf *MemFile) GetData() []byte


func (mf *MemFile) GetSafeData() *safe.SecretPlaintext


func GetAbsFilename(filename string) string


func AddFileToMemory(mf MemFile)


func ReadFileFromMemory(name string) *MemFile


func ReadDirFromMemory(name string) []MemFile


func LoadFileToMemory(filename string)


func SafeLoadFileToMemory(filename string)


func LoadFileToMemoryWithCompress(filename string)


func SafeLoadFileToMemoryWithCompress(filename string)


func loadFileToMemory(filename string, compress bool, isSafe bool)


func LoadFileToB64(filename string) *MemFileB64


func LoadFilesToMemoryFromB64(b64File *MemFileB64)


func LoadFilesToMemoryFromB64KeepGzip(b64File *MemFileB64)


func LoadFilesToMemoryFromJson(jsonFiles string)

```

### file > archive.go
```go



func Compress(data []byte, cType string) ([]byte, error)


func Decompress(data []byte, cType string) ([]byte, error)


func Extract(srcFile, destDir string, stripRoot bool) error


func extractTar(r io.Reader, dest string, strip bool) error


func extractZip(src, dest string, strip bool) error


func extractSingleFile(r io.Reader, destDir, fileName string) error


func writeEntry(destDir, name string, mode os.FileMode, isDir bool, linkPath string, r io.Reader, strip bool) error


func Archive(srcPath, destFile string) error


func walkAndAdd(srcPath string, addFn func(string, os.FileInfo, io.Reader) error) error

```

### file > file.go
```go



type FileInfo struct {
	Name     string
	FullName string
	IsDir    bool
	Size     int64
	ModTime  int64
}

func EnsureParentDir(filename string)


func EnsureDirSuffix(path string) string


func Exists(filename string) bool


func GetFileInfo(filename string) *FileInfo


func ReadBytes(filename string) ([]byte, error)


func Read(filename string) (string, error)


func ReadLines(filename string) ([]string, error)


func WriteBytes(filename string, content []byte) error


func Write(filename string, content string) error


func Copy(from, to string) error


func CopyToFile(from io.Reader, to string) error


func Remove(path string) error


func Move(src, dst string) error


func Replace(filename, old, new string) error


func Search(dir, pattern string) []string


func ReadDir(filename string) ([]FileInfo, error)


func RunCommand(command string, args ...string) ([]string, error)

```

