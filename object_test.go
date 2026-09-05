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

	yamlFile := filepath.Join(tmpDir, "test.yaml")
	if err := os.WriteFile(yamlFile, []byte("USER-NAME: Alice\n"), 0644); err != nil {
		t.Fatal(err)
	}
	t.Run("YAMLSmartMapping", func(t *testing.T) {
		var value SnakeStruct
		if err := UnmarshalFile(yamlFile, &value); err != nil {
			t.Fatal(err)
		}
		if value.UserName != "Alice" {
			t.Fatalf("expected Alice, got %s", value.UserName)
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
