package file

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"sync"

	"apigo.cc/go/cast"
)

var fileLocksLock = sync.Mutex{}
var fileLocks = map[string]*sync.Mutex{}

func getLock(filename string) *sync.Mutex {
	absName := GetAbsFilename(filename)
	fileLocksLock.Lock()
	defer fileLocksLock.Unlock()
	if fileLocks[absName] == nil {
		fileLocks[absName] = new(sync.Mutex)
	}
	return fileLocks[absName]
}

// UnmarshalFile 从文件读取数据并映射到 to，自动处理 YAML/JSON 格式判别及智能字段映射。
func UnmarshalFile(filename string, to any) error {
	data, err := ReadBytes(filename)
	if err != nil {
		return err
	}

	isYaml := strings.HasSuffix(filename, ".yml") || strings.HasSuffix(filename, ".yaml")
	if isYaml {
		return cast.UnmarshalYAML(data, to)
	}

	var in any
	if err = json.Unmarshal(data, &in); err != nil {
		return err
	}
	cast.Convert(to, in)
	return nil
}

// MarshalFile 将数据序列化并写入文件，支持 YAML/JSON 自动判别。
func MarshalFile(filename string, data any) error {
	return marshalFile(filename, data, false)
}

// MarshalFilePretty 类似 MarshalFile，但会输出带缩进的格式化内容。
func MarshalFilePretty(filename string, data any) error {
	return marshalFile(filename, data, true)
}

func marshalFile(filename string, data any, indent bool) error {
	isYaml := strings.HasSuffix(filename, ".yml") || strings.HasSuffix(filename, ".yaml")

	var buf []byte
	var err error
	if isYaml {
		buf, err = cast.ToYAMLBytes(data)
	} else {
		buffer := bytes.Buffer{}
		enc := json.NewEncoder(&buffer)
		enc.SetEscapeHTML(false)
		if indent {
			enc.SetIndent("", "  ")
		}
		err = enc.Encode(data)
		if err == nil {
			buf = buffer.Bytes()
		}
	}
	if err != nil {
		return err
	}

	lock := getLock(filename)
	lock.Lock()
	defer lock.Unlock()

	return WriteBytes(filename, buf)
}

func PatchFile(filename string, patch any) error {
	var current map[string]any
	err := UnmarshalFile(filename, &current)
	if err != nil {
		if os.IsNotExist(err) {
			current = make(map[string]any)
		} else {
			return err
		}
	}
	cast.Convert(&current, patch)
	return MarshalFilePretty(filename, current)
}
