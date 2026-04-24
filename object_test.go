package file

import (
	"os"
	"path/filepath"
	"testing"
)

type TestStruct struct {
	Name string `json:"name" yaml:"name"`
	Age  int    `json:"age" yaml:"age"`
}

type SnakeStruct struct {
	UserName string `json:"user_name" yaml:"user_name"`
}

func TestUnmarshalFile(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "test_object")
	defer os.RemoveAll(tmpDir)

	jsonFile := filepath.Join(tmpDir, "test.json")
	// 测试智能映射：输入 user_name，目标 UserName
	jsonData := `{"user_name": "Alice"}`
	os.WriteFile(jsonFile, []byte(jsonData), 0644)

	t.Run("SmartMapping", func(t *testing.T) {
		var ts SnakeStruct
		err := UnmarshalFile(jsonFile, &ts)
		if err != nil {
			t.Errorf("UnmarshalFile failed: %v", err)
		}
		if ts.UserName != "Alice" {
			t.Errorf("Expected Alice, got %s", ts.UserName)
		}
	})
}

func TestMarshalFile(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "test_marshal")
	defer os.RemoveAll(tmpDir)

	data := TestStruct{Name: "Bob", Age: 25}

	t.Run("Pretty", func(t *testing.T) {
		filename := filepath.Join(tmpDir, "pretty.json")
		err := MarshalFilePretty(filename, data)
		if err != nil {
			t.Fatalf("MarshalFilePretty failed: %v", err)
		}
	})
}

func TestPatchFile(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "test_patch")
	defer os.RemoveAll(tmpDir)
	filename := filepath.Join(tmpDir, "config.json")

	t.Run("Update", func(t *testing.T) {
		patch := map[string]any{"name": "Initial"}
		err := PatchFile(filename, patch)
		if err != nil {
			t.Fatalf("PatchFile failed: %v", err)
		}
	})
}
