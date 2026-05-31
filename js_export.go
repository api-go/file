package file

import (
	"context"

	"apigo.cc/go/jsmod"
)

func init() {
	jsmod.Register("file", map[string]any{
		// 读操作 (映射到私有包装器)
		"exists": func(ctx context.Context, path string) bool {
			p, err := VerifyPathForSafeMode(ctx, path)
			if err != nil {
				return false
			}
			return Exists(p)
		},
		"read": func(ctx context.Context, path string) (string, error) {
			p, err := VerifyPathForSafeMode(ctx, path)
			if err != nil {
				return "", err
			}
			return Read(p)
		},
		"readBytes": func(ctx context.Context, path string) ([]byte, error) {
			p, err := VerifyPathForSafeMode(ctx, path)
			if err != nil {
				return nil, err
			}
			return ReadBytes(p)
		},
		"readLines": func(ctx context.Context, path string) ([]string, error) {
			p, err := VerifyPathForSafeMode(ctx, path)
			if err != nil {
				return nil, err
			}
			return ReadLines(p)
		},
		"readDir": func(ctx context.Context, path string) ([]FileInfo, error) {
			p, err := VerifyPathForSafeMode(ctx, path)
			if err != nil {
				return nil, err
			}
			return ReadDir(p)
		},
		"getFileInfo": func(ctx context.Context, path string) *FileInfo {
			p, err := VerifyPathForSafeMode(ctx, path)
			if err != nil {
				return nil
			}
			return GetFileInfo(p)
		},

		// 写操作
		"write": func(ctx context.Context, path string, content string) error {
			p, err := VerifyPathForSafeMode(ctx, path)
			if err != nil {
				return err
			}
			return Write(p, content)
		},
		"writeBytes": func(ctx context.Context, path string, content []byte) error {
			p, err := VerifyPathForSafeMode(ctx, path)
			if err != nil {
				return err
			}
			return WriteBytes(p, content)
		},
		"remove": func(ctx context.Context, path string) error {
			p, err := VerifyPathForSafeMode(ctx, path)
			if err != nil {
				return err
			}
			return Remove(p)
		},
		"mkdir": func(ctx context.Context, path string) error {
			p, err := VerifyPathForSafeMode(ctx, path)
			if err != nil {
				return err
			}
			return Mkdir(p)
		},
		"copy": func(ctx context.Context, from, to string) error {
			pFrom, err := VerifyPathForSafeMode(ctx, from)
			if err != nil {
				return err
			}
			pTo, err := VerifyPathForSafeMode(ctx, to)
			if err != nil {
				return err
			}
			return Copy(pFrom, pTo)
		},
		"move": func(ctx context.Context, from, to string) error {
			pFrom, err := VerifyPathForSafeMode(ctx, from)
			if err != nil {
				return err
			}
			pTo, err := VerifyPathForSafeMode(ctx, to)
			if err != nil {
				return err
			}
			return Move(pFrom, pTo)
		},
		"replace": func(ctx context.Context, path, old, new string) error {
			p, err := VerifyPathForSafeMode(ctx, path)
			if err != nil {
				return err
			}
			return Replace(p, old, new)
		},

		// 序列化
		"unmarshalFile": func(ctx context.Context, path string, to any) error {
			p, err := VerifyPathForSafeMode(ctx, path)
			if err != nil {
				return err
			}
			return UnmarshalFile(p, to)
		},
		"marshalFile": func(ctx context.Context, path string, data any) error {
			p, err := VerifyPathForSafeMode(ctx, path)
			if err != nil {
				return err
			}
			return MarshalFile(p, data)
		},
		"marshalFilePretty": func(ctx context.Context, path string, data any) error {
			p, err := VerifyPathForSafeMode(ctx, path)
			if err != nil {
				return err
			}
			return MarshalFilePretty(p, data)
		},

		// 归档
		"archive": func(ctx context.Context, src, dest string) error {
			pSrc, err := VerifyPathForSafeMode(ctx, src)
			if err != nil {
				return err
			}
			pDest, err := VerifyPathForSafeMode(ctx, dest)
			if err != nil {
				return err
			}
			return Archive(pSrc, pDest)
		},
		"extract": func(ctx context.Context, src, dest string, strip bool) error {
			pSrc, err := VerifyPathForSafeMode(ctx, src)
			if err != nil {
				return err
			}
			pDest, err := VerifyPathForSafeMode(ctx, dest)
			if err != nil {
				return err
			}
			return Extract(pSrc, pDest, strip)
		},

		// 压缩工具 (无路径，不校验)
		"compress":   Compress,
		"decompress": Decompress,
	})
}
