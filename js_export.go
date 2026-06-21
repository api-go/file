package file

import (
	"context"

	"apigo.cc/go/jsmod"
)

func init() {
	jsmod.Register("file", map[string]any{
		// 读操作 (映射到私有包装器)
		"Exists":            jsExists,
		"Read":              jsRead,
		"ReadBytes":         jsReadBytes,
		"ReadLines":         jsReadLines,
		"ReadDir":           jsReadDir,
		"GetFileInfo":       jsGetFileInfo,

		// 写操作
		"Write":             jsWrite,
		"WriteBytes":        jsWriteBytes,
		"Remove":            jsRemove,
		"Mkdir":             jsMkdir,
		"Copy":              jsCopy,
		"Move":              jsMove,
		"Replace":           jsReplace,

		// 序列化
		"UnmarshalFile":     jsUnmarshalFile,
		"MarshalFile":       jsMarshalFile,
		"MarshalFilePretty": jsMarshalFilePretty,

		// 归档
		"Archive":           jsArchive,
		"Extract":           jsExtract,

		// 压缩工具
		"Compress":          jsCompress,
		"Decompress":        jsDecompress,
	})
}

func jsExists(ctx context.Context, path string) bool {
	p, err := VerifyPathForSafeMode(ctx, path)
	if err != nil {
		return false
	}
	return Exists(p)
}

func jsRead(ctx context.Context, path string) (string, error) {
	p, err := VerifyPathForSafeMode(ctx, path)
	if err != nil {
		return "", jsmod.MakeError(err)
	}
	res, err := Read(p)
	if err != nil {
		return "", jsmod.MakeError(err)
	}
	return res, nil
}

func jsReadBytes(ctx context.Context, path string) ([]byte, error) {
	p, err := VerifyPathForSafeMode(ctx, path)
	if err != nil {
		return nil, jsmod.MakeError(err)
	}
	res, err := ReadBytes(p)
	if err != nil {
		return nil, jsmod.MakeError(err)
	}
	return res, nil
}

func jsReadLines(ctx context.Context, path string) ([]string, error) {
	p, err := VerifyPathForSafeMode(ctx, path)
	if err != nil {
		return nil, jsmod.MakeError(err)
	}
	res, err := ReadLines(p)
	if err != nil {
		return nil, jsmod.MakeError(err)
	}
	return res, nil
}

func jsReadDir(ctx context.Context, path string) ([]FileInfo, error) {
	p, err := VerifyPathForSafeMode(ctx, path)
	if err != nil {
		return nil, jsmod.MakeError(err)
	}
	res, err := ReadDir(p)
	if err != nil {
		return nil, jsmod.MakeError(err)
	}
	return res, nil
}

func jsGetFileInfo(ctx context.Context, path string) *FileInfo {
	p, err := VerifyPathForSafeMode(ctx, path)
	if err != nil {
		return nil
	}
	return GetFileInfo(p)
}

func jsWrite(ctx context.Context, path string, content string) error {
	p, err := VerifyPathForSafeMode(ctx, path)
	if err != nil {
		return jsmod.MakeError(err)
	}
	err = Write(p, content)
	if err != nil {
		return jsmod.MakeError(err)
	}
	return nil
}

func jsWriteBytes(ctx context.Context, path string, content []byte) error {
	p, err := VerifyPathForSafeMode(ctx, path)
	if err != nil {
		return jsmod.MakeError(err)
	}
	err = WriteBytes(p, content)
	if err != nil {
		return jsmod.MakeError(err)
	}
	return nil
}

func jsRemove(ctx context.Context, path string) error {
	p, err := VerifyPathForSafeMode(ctx, path)
	if err != nil {
		return jsmod.MakeError(err)
	}
	err = Remove(p)
	if err != nil {
		return jsmod.MakeError(err)
	}
	return nil
}

func jsMkdir(ctx context.Context, path string) error {
	p, err := VerifyPathForSafeMode(ctx, path)
	if err != nil {
		return jsmod.MakeError(err)
	}
	err = Mkdir(p)
	if err != nil {
		return jsmod.MakeError(err)
	}
	return nil
}

func jsCopy(ctx context.Context, from, to string) error {
	pFrom, err := VerifyPathForSafeMode(ctx, from)
	if err != nil {
		return jsmod.MakeError(err)
	}
	pTo, err := VerifyPathForSafeMode(ctx, to)
	if err != nil {
		return jsmod.MakeError(err)
	}
	err = Copy(pFrom, pTo)
	if err != nil {
		return jsmod.MakeError(err)
	}
	return nil
}

func jsMove(ctx context.Context, from, to string) error {
	pFrom, err := VerifyPathForSafeMode(ctx, from)
	if err != nil {
		return jsmod.MakeError(err)
	}
	pTo, err := VerifyPathForSafeMode(ctx, to)
	if err != nil {
		return jsmod.MakeError(err)
	}
	err = Move(pFrom, pTo)
	if err != nil {
		return jsmod.MakeError(err)
	}
	return nil
}

func jsReplace(ctx context.Context, path, old, new string) error {
	p, err := VerifyPathForSafeMode(ctx, path)
	if err != nil {
		return jsmod.MakeError(err)
	}
	err = Replace(p, old, new)
	if err != nil {
		return jsmod.MakeError(err)
	}
	return nil
}

func jsUnmarshalFile(ctx context.Context, path string, to any) error {
	p, err := VerifyPathForSafeMode(ctx, path)
	if err != nil {
		return jsmod.MakeError(err)
	}
	err = UnmarshalFile(p, to)
	if err != nil {
		return jsmod.MakeError(err)
	}
	return nil
}

func jsMarshalFile(ctx context.Context, path string, data any) error {
	p, err := VerifyPathForSafeMode(ctx, path)
	if err != nil {
		return jsmod.MakeError(err)
	}
	err = MarshalFile(p, data)
	if err != nil {
		return jsmod.MakeError(err)
	}
	return nil
}

func jsMarshalFilePretty(ctx context.Context, path string, data any) error {
	p, err := VerifyPathForSafeMode(ctx, path)
	if err != nil {
		return jsmod.MakeError(err)
	}
	err = MarshalFilePretty(p, data)
	if err != nil {
		return jsmod.MakeError(err)
	}
	return nil
}

func jsArchive(ctx context.Context, src, dest string) error {
	pSrc, err := VerifyPathForSafeMode(ctx, src)
	if err != nil {
		return jsmod.MakeError(err)
	}
	pDest, err := VerifyPathForSafeMode(ctx, dest)
	if err != nil {
		return jsmod.MakeError(err)
	}
	err = Archive(pSrc, pDest)
	if err != nil {
		return jsmod.MakeError(err)
	}
	return nil
}

func jsExtract(ctx context.Context, src, dest string, strip bool) error {
	pSrc, err := VerifyPathForSafeMode(ctx, src)
	if err != nil {
		return jsmod.MakeError(err)
	}
	pDest, err := VerifyPathForSafeMode(ctx, dest)
	if err != nil {
		return jsmod.MakeError(err)
	}
	err = Extract(pSrc, pDest, strip)
	if err != nil {
		return jsmod.MakeError(err)
	}
	return nil
}

func jsCompress(data []byte, cType string) ([]byte, error) {
	res, err := Compress(data, cType)
	if err != nil {
		return nil, jsmod.MakeError(err)
	}
	return res, nil
}

func jsDecompress(data []byte, cType string) ([]byte, error) {
	res, err := Decompress(data, cType)
	if err != nil {
		return nil, jsmod.MakeError(err)
	}
	return res, nil
}
