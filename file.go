package file

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"apigo.cc/go/cast"
	"apigo.cc/go/jsmod"
)

type FileInfo struct {
	Name     string
	FullName string
	IsDir    bool
	Size     int64
	ModTime  int64
}

func EnsureParentDir(filename string) {
	pos := strings.LastIndexByte(filename, os.PathSeparator)
	if pos < 0 {
		return
	}
	_ = os.MkdirAll(filename[:pos], 0755)
}

func EnsureDirSuffix(path string) string {
	const sep = string(os.PathSeparator)
	if !strings.HasSuffix(path, sep) {
		return path + sep
	}
	return path
}

func Exists(filename string) bool {
	if mf := ReadFileFromMemory(filename); mf != nil {
		return true
	}
	fi, err := os.Stat(filename)
	return err == nil && fi != nil
}

func GetFileInfo(filename string) *FileInfo {
	if mf := ReadFileFromMemory(filename); mf != nil {
		return &FileInfo{
			Name:     mf.Name,
			FullName: mf.AbsName,
			IsDir:    mf.IsDir,
			Size:     mf.Size,
			ModTime:  mf.ModTime.Unix(),
		}
	}
	if fi, err := os.Stat(filename); err == nil {
		fullName := filename
		if !filepath.IsAbs(filename) {
			fullName, _ = filepath.Abs(filename)
		}
		return &FileInfo{
			Name:     filename,
			FullName: fullName,
			IsDir:    fi.IsDir(),
			Size:     fi.Size(),
			ModTime:  fi.ModTime().Unix(),
		}
	}
	return nil
}

func ReadBytes(filename string) ([]byte, error) {
	if mf := ReadFileFromMemory(filename); mf != nil {
		return mf.GetData(), nil
	}
	return os.ReadFile(filename)
}

func Read(filename string) (string, error) {
	buf, err := ReadBytes(filename)
	return string(buf), err
}

func ReadLines(filename string) ([]string, error) {
	if mf := ReadFileFromMemory(filename); mf != nil {
		return strings.Split(string(mf.GetData()), "\n"), nil
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func WriteBytes(filename string, content []byte) error {
	absFilename := GetAbsFilename(filename)
	memFilesLock.RLock()
	mf := memFiles[absFilename]
	memFilesLock.RUnlock()
	if mf != nil {
		memFilesLock.Lock()
		mf.Data = content
		mf.Size = int64(len(content))
		mf.ModTime = time.Now()
		memFilesLock.Unlock()
	}
	EnsureParentDir(filename)
	return os.WriteFile(filename, content, 0644)
}

func Write(filename string, content string) error {
	return WriteBytes(filename, []byte(content))
}

func Append(filename string, content []byte) error {
	EnsureParentDir(filename)
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(content)
	return err
}

func Copy(from, to string) error {
	fromStat, err := os.Stat(from)
	if err != nil {
		return err
	}
	if fromStat.IsDir() {
		entries, err := os.ReadDir(from)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			err := Copy(filepath.Join(from, entry.Name()), filepath.Join(to, entry.Name()))
			if err != nil {
				return err
			}
		}
		return nil
	}

	toStat, err := os.Stat(to)
	if err == nil && toStat.IsDir() {
		to = filepath.Join(to, filepath.Base(from))
	}
	EnsureParentDir(to)

	src, err := os.Open(from)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.OpenFile(to, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func CopyToFile(from io.Reader, to string) error {
	EnsureParentDir(to)
	fp, err := os.OpenFile(to, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer fp.Close()
	_, err = io.Copy(fp, from)
	return err
}

func Remove(path string) error {
	return os.RemoveAll(path)
}

func Mkdir(path string) error {
	return os.MkdirAll(path, 0755)
}

func Open(filename string) (*os.File, error) {
	return os.Open(filename)
}

func Create(filename string) (*os.File, error) {
	EnsureParentDir(filename)
	return os.Create(filename)
}

func Move(src, dst string) error {
	EnsureParentDir(dst)
	return os.Rename(src, dst)
}

func Replace(filename, old, new string) error {
	content, err := Read(filename)
	if err != nil {
		return err
	}
	newContent := strings.ReplaceAll(content, old, new)
	return Write(filename, newContent)
}

func Search(dir, pattern string) []string {
	var matches []string
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			matched, _ := filepath.Match(pattern, filepath.Base(path))
			if matched {
				matches = append(matches, path)
			}
		}
		return nil
	})
	return matches
}

func ReadDir(filename string) ([]FileInfo, error) {
	var out []FileInfo
	if mfList := ReadDirFromMemory(filename); mfList != nil {
		for _, f := range mfList {
			out = append(out, FileInfo{
				Name:     filepath.Base(f.Name),
				FullName: f.Name,
				IsDir:    f.IsDir,
				Size:     f.Size,
				ModTime:  f.ModTime.Unix(),
			})
		}
		return out, nil
	}

	files, err := os.ReadDir(filename)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		info, _ := f.Info()
		if info != nil {
			out = append(out, FileInfo{
				Name:     f.Name(),
				FullName: filepath.Join(filename, f.Name()),
				IsDir:    info.IsDir(),
				Size:     info.Size(),
				ModTime:  info.ModTime().Unix(),
			})
		}
	}
	return out, nil
}

func RunCommand(command string, args ...string) ([]string, error) {
	cmd := exec.Command(command, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("command execution failed: %w, output: %s", err, string(out))
	}
	rawLines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var lines []string
	for _, line := range rawLines {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			lines = append(lines, trimmed)
		}
	}
	return lines, nil
}

// RealPath 返回客观的物理绝对路径，处理软连接和相对路径
func RealPath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}

	// 尝试解析软连接
	real, err := filepath.EvalSymlinks(abs)
	if err == nil {
		return real
	}

	// 如果文件不存在，EvalSymlinks 会失败，此时返回绝对路径即可（后续写入操作会处理）
	return abs
}

// VerifyPathForSafeMode 在安全模式下校验路径是否在 AllowedDirs 白名单内
func VerifyPathForSafeMode(ctx context.Context, path string) (string, error) {
	// 如果不是安全模式，直接返回 RealPath
	if ctx == nil || !jsmod.IsSafeMode(ctx) {
		return RealPath(path), nil
	}

	// 1. 获取白名单 (约定从 ctx 获取 "AllowedDirs")
	var allowedDirs []string
	cast.Convert(&allowedDirs, ctx.Value("AllowedDirs"))
	if len(allowedDirs) == 0 {
		return "", fmt.Errorf("file: access denied, AllowedDirs not found in context")
	}

	// 2. 获取物理路径
	real := RealPath(path)

	// 3. 校验逻辑：判断当前物理路径是否以任何一个白名单目录为前缀
	isAllowed := false
	for _, dir := range allowedDirs {
		if strings.HasPrefix(real, dir) {
			isAllowed = true
			break
		}
	}

	// 4. 写操作保护：如果路径不存在，追溯其存在的父目录权限
	if !isAllowed {
		temp := real
		for temp != "" && temp != string(os.PathSeparator) {
			temp = filepath.Dir(temp)
			if _, err := os.Stat(temp); err == nil {
				// 找到了第一个存在的祖先目录，检查它是否受控
				for _, dir := range allowedDirs {
					if strings.HasPrefix(temp, dir) {
						isAllowed = true
						break
					}
				}
				break
			}
		}
	}

	if !isAllowed {
		return "", fmt.Errorf("file: permission denied for path %s", path)
	}

	return real, nil
}
