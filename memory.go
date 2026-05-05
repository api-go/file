package file

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"apigo.cc/go/cast"
	"apigo.cc/go/encoding"
	"apigo.cc/go/safe"
)

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

func (mf *MemFile) GetData() []byte {
	if mf.Compressed && len(mf.Data) > 0 {
		return cast.As(Decompress(mf.Data, "gzip"))
	}
	return mf.Data
}

func (mf *MemFile) GetSafeData() *safe.SecretPlaintext {
	if mf.SafeData != nil {
		return mf.SafeData.Open()
	}
	return nil
}

func GetAbsFilename(filename string) string {
	if !filepath.IsAbs(filename) {
		if absName, err := filepath.Abs(filename); err == nil {
			filename = absName
		}
	}
	return filename
}

func AddFileToMemory(mf MemFile) {
	mf.Name = GetAbsFilename(mf.Name)
	mf.AbsName = mf.Name
	dirName := filepath.Dir(mf.Name)
	memFilesLock.Lock()
	defer memFilesLock.Unlock()
	memFiles[mf.Name] = &mf
	if dirName != "" && dirName != "." {
		memFilesByDir[dirName] = append(memFilesByDir[dirName], mf)
	}
}

func ReadFileFromMemory(name string) *MemFile {
	name = GetAbsFilename(name)
	memFilesLock.RLock()
	defer memFilesLock.RUnlock()
	return memFiles[name]
}

func ReadDirFromMemory(name string) []MemFile {
	name = GetAbsFilename(name)
	memFilesLock.RLock()
	defer memFilesLock.RUnlock()
	mfList := memFilesByDir[name]
	if mfList == nil {
		return nil
	}
	out := make([]MemFile, len(mfList))
	copy(out, mfList)
	return out
}

func LoadFileToMemory(filename string) {
	loadFileToMemory(filename, false, false)
}

func SafeLoadFileToMemory(filename string) {
	loadFileToMemory(filename, false, true)
}

func LoadFileToMemoryWithCompress(filename string) {
	loadFileToMemory(filename, true, false)
}

func SafeLoadFileToMemoryWithCompress(filename string) {
	loadFileToMemory(filename, true, true)
}

func loadFileToMemory(filename string, compress bool, isSafe bool) {
	if info, err := os.Stat(filename); err == nil {
		if info.IsDir() {
			AddFileToMemory(MemFile{
				Name:    filename,
				ModTime: info.ModTime(),
				IsDir:   true,
				Size:    info.Size(),
			})
			if files, err := os.ReadDir(filename); err == nil {
				for _, file := range files {
					loadFileToMemory(filepath.Join(filename, file.Name()), compress, isSafe)
				}
			}
		} else {
			if data, err := os.ReadFile(filename); err == nil {
				compressed := false
				var dataBuf *safe.SafeBuf
				if compress {
					if data2, err := Compress(data, "gzip"); err == nil {
						if isSafe {
							safe.ZeroMemory(data)
						}
						data = data2
						compressed = true
					}
				}
				if isSafe {
					dataBuf = safe.NewSafeBuf(data)
					safe.ZeroMemory(data)
					data = nil
				}
				AddFileToMemory(MemFile{
					Name:       filename,
					ModTime:    info.ModTime(),
					IsDir:      false,
					Size:       info.Size(),
					Data:       data,
					SafeData:   dataBuf,
					Compressed: compressed,
				})
			}
		}
	}
}

func LoadFileToB64(filename string) *MemFileB64 {
	if info, err := os.Stat(filename); err == nil {
		if info.IsDir() {
			out := MemFileB64{
				Name:     filename,
				ModTime:  info.ModTime(),
				IsDir:    true,
				Size:     info.Size(),
				Children: make([]MemFileB64, 0),
			}
			if files, err := os.ReadDir(filename); err == nil {
				for _, file := range files {
					if mfB64 := LoadFileToB64(filepath.Join(filename, file.Name())); mfB64 != nil {
						out.Children = append(out.Children, *mfB64)
					}
				}
			}
			return &out
		} else {
			if data, err := os.ReadFile(filename); err == nil {
				compressed := false
				if buf, err := Compress(data, "gzip"); err == nil {
					data = buf
					compressed = true
				}
				return &MemFileB64{
					Name:       filename,
					ModTime:    info.ModTime(),
					IsDir:      false,
					Size:       info.Size(),
					DataB64:    encoding.Base64(data),
					Compressed: compressed,
				}
			}
		}
	}
	return nil
}

func LoadFilesToMemoryFromB64(b64File *MemFileB64) {
	if data, err := encoding.UnBase64(b64File.DataB64); err == nil {
		if b64File.Compressed {
			data = cast.As(Decompress(data, "gzip"))
		}
		memFile := MemFile{
			Name:    b64File.Name,
			ModTime: b64File.ModTime,
			IsDir:   b64File.IsDir,
			Size:    b64File.Size,
			Data:    data,
		}
		AddFileToMemory(memFile)
		if memFile.IsDir && len(b64File.Children) > 0 {
			for _, child := range b64File.Children {
				LoadFilesToMemoryFromB64(&child)
			}
		}
	}
}

func LoadFilesToMemoryFromB64KeepGzip(b64File *MemFileB64) {
	if data, err := encoding.UnBase64(b64File.DataB64); err == nil {
		memFile := MemFile{
			Name:    b64File.Name,
			ModTime: b64File.ModTime,
			IsDir:   b64File.IsDir,
			Size:    b64File.Size,
			Data:    data,
		}
		AddFileToMemory(memFile)
		if memFile.IsDir && len(b64File.Children) > 0 {
			for _, child := range b64File.Children {
				LoadFilesToMemoryFromB64KeepGzip(&child)
			}
		}
	}
}

func LoadFilesToMemoryFromJson(jsonFiles string) {
	dirData := MemFileB64{}
	_ = json.Unmarshal([]byte(jsonFiles), &dirData)
	LoadFilesToMemoryFromB64(&dirData)
}
