package file

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestMemoryDeep(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "test_mem_deep")
	defer os.RemoveAll(tmpDir)
    absDir, _ := filepath.Abs(tmpDir)

	filePath := filepath.Join(absDir, "data.txt")
	content := []byte("secret_data")
	WriteBytes(filePath, content)
	
	// 1. 验证 SafeBuf 链路 (无压缩)
	SafeLoadFileToMemory(filePath) // 仅加载并加密
	mfSafe := ReadFileFromMemory(filePath)
	if mfSafe == nil || mfSafe.SafeData == nil {
		t.Fatal("SafeData should not be nil")
	}
	
	plain := mfSafe.GetSafeData()
	if !bytes.Equal(plain.Data, content) {
		t.Error("SafeData decryption mismatch")
	}
	plain.Close()
}
