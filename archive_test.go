package file

import (
	"os"
	"path/filepath"
	"testing"
)

func TestArchiveComprehensive(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "test_arch_full")
	defer os.RemoveAll(tmpDir)

	srcDir := filepath.Join(tmpDir, "src")
	os.Mkdir(srcDir, 0755)
	_ = Write(filepath.Join(srcDir, "test.txt"), "hello")
	
	// Test Zip
	zipFile := filepath.Join(tmpDir, "out.zip")
	Archive(srcDir, zipFile)
	
	outDir := filepath.Join(tmpDir, "out_zip")
	Extract(zipFile, outDir, false)
	if !Exists(filepath.Join(outDir, "src/test.txt")) {
		t.Error("Zip extract failed")
	}

    // Test Tar.Gz
    tgzFile := filepath.Join(tmpDir, "out.tar.gz")
    Archive(srcDir, tgzFile)
    
    outDir2 := filepath.Join(tmpDir, "out_tgz")
    Extract(tgzFile, outDir2, false)
    if !Exists(filepath.Join(outDir2, "src/test.txt")) {
        t.Error("Tar.gz extract failed")
    }
}
