package file_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/api-go/cast"
	"github.com/api-go/file"
)

func TestFileBasic(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "file_test")
	defer os.RemoveAll(tmpDir)

	fileA := filepath.Join(tmpDir, "a.txt")
	fileB := filepath.Join(tmpDir, "b.txt")

	_ = file.Write(fileA, "hello world")
	if !file.Exists(fileA) {
		t.Error("File should exist")
	}

	_ = file.Move(fileA, fileB)
	if file.Exists(fileA) {
		t.Error("File A should be moved")
	}

	if cast.As(file.Read(fileB)) != "hello world" {
		t.Error("Read content mismatch")
	}

	info := file.GetFileInfo(fileB)
	if info == nil || info.Name != fileB {
		t.Error("GetFileInfo failed")
	}

	_ = file.Replace(fileB, "hello", "hi")
	if cast.As(file.Read(fileB)) != "hi world" {
		t.Error("Replace failed")
	}

	_ = file.Remove(fileB)
	if file.Exists(fileB) {
		t.Error("Remove failed")
	}
}

func TestReadLines(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "file_test_lines")
	defer os.RemoveAll(tmpDir)

	path := filepath.Join(tmpDir, "lines.txt")
	content := "line1\nline2\nline3"
	_ = file.Write(path, content)

	lines := cast.As(file.ReadLines(path))
	if len(lines) != 3 || lines[1] != "line2" {
		t.Errorf("ReadLines failed: got %v", lines)
	}
}

func TestReadDir(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "file_test_dir")
	defer os.RemoveAll(tmpDir)

	_ = file.Write(filepath.Join(tmpDir, "f1.txt"), "1")
	_ = file.Write(filepath.Join(tmpDir, "f2.txt"), "2")

	files := cast.As(file.ReadDir(tmpDir))
	if len(files) != 2 {
		t.Errorf("ReadDir failed: got %d files", len(files))
	}
}
