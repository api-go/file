package file

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
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
	path := filename[0:pos]
	if _, err := os.Stat(path); err != nil {
		_ = os.MkdirAll(path, 0755)
	}
}

func EnsureDirSuffix(path string) string {
	const spe = string(os.PathSeparator)
	if !strings.HasSuffix(path, spe) {
		return path + spe
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

func MustReadBytes(filename string) []byte {
	buf, _ := ReadBytes(filename)
	return buf
}

func Read(filename string) (string, error) {
	buf, err := ReadBytes(filename)
	return string(buf), err
}

func MustRead(filename string) string {
	buf, _ := Read(filename)
	return buf
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

func MustReadLines(filename string) []string {
	lines, _ := ReadLines(filename)
	return lines
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

func MustReadDir(filename string) []FileInfo {
	files, _ := ReadDir(filename)
	return files
}

func RunCommand(command string, args ...string) ([]string, error) {
	cmd := exec.Command(command, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	return lines, nil
}
