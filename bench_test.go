package file

import (
	"os"
	"path/filepath"
	"testing"
)

func BenchmarkMarshalFile(b *testing.B) {
	tmpDir, _ := os.MkdirTemp("", "bench_marshal")
	defer os.RemoveAll(tmpDir)
	filename := filepath.Join(tmpDir, "test.json")
	data := TestStruct{Name: "Benchmark", Age: 100}

	b.Run("Normal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = MarshalFile(filename, data)
		}
	})

	b.Run("Pretty", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = MarshalFilePretty(filename, data)
		}
	})
}

func BenchmarkUnmarshalFile(b *testing.B) {
	tmpDir, _ := os.MkdirTemp("", "bench_unmarshal")
	defer os.RemoveAll(tmpDir)
	filename := filepath.Join(tmpDir, "test.json")
	data := `{"name": "Benchmark", "age": 100}`
	os.WriteFile(filename, []byte(data), 0644)

	var ts TestStruct
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = UnmarshalFile(filename, &ts)
	}
}
