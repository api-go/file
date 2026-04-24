package file

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileSystemOps(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "test_fs")
	defer os.RemoveAll(tmpDir)

	fileA := filepath.Join(tmpDir, "a.txt")
	fileB := filepath.Join(tmpDir, "b.txt")
	_ = Write(fileA, "hello world")

	t.Run("Copy", func(t *testing.T) {
		err := Copy(fileA, fileB)
		if err != nil || !Exists(fileB) {
			t.Errorf("Copy failed: %v", err)
		}
		if MustRead(fileB) != "hello world" {
			t.Error("Copy content mismatch")
		}
	})

	t.Run("Move", func(t *testing.T) {
		fileC := filepath.Join(tmpDir, "c.txt")
		err := Move(fileB, fileC)
		if err != nil || !Exists(fileC) || Exists(fileB) {
			t.Errorf("Move failed: %v", err)
		}
	})

	t.Run("Replace", func(t *testing.T) {
		err := Replace(fileA, "hello", "hi")
		if err != nil || MustRead(fileA) != "hi world" {
			t.Errorf("Replace failed: %v", err)
		}
	})

	t.Run("Remove", func(t *testing.T) {
		err := Remove(fileA)
		if err != nil || Exists(fileA) {
			t.Errorf("Remove failed: %v", err)
		}
	})
}

func TestSearchAndDir(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "test_search")
	defer os.RemoveAll(tmpDir)

	_ = Write(filepath.Join(tmpDir, "1.txt"), "data")
	subDir := filepath.Join(tmpDir, "sub")
	_ = os.Mkdir(subDir, 0755)
	_ = Write(filepath.Join(subDir, "2.txt"), "data")

	t.Run("Search", func(t *testing.T) {
		matches := Search(tmpDir, "*.txt")
		if len(matches) != 2 {
			t.Errorf("Expected 2 matches, got %d", len(matches))
		}
	})

	t.Run("ReadDir", func(t *testing.T) {
		files, err := ReadDir(tmpDir)
		if err != nil || len(files) < 2 {
			t.Errorf("ReadDir failed: %v", err)
		}
	})
}

func TestHelpers(t *testing.T) {
    // Test EnsureParentDir (via Write)
    path := "test_helper/file.txt"
    _ = Write(path, "d")
    if !Exists(path) {
        t.Error("Helper Write failed")
    }
    Remove("test_helper")
}
