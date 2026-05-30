package file

import "apigo.cc/go/jsmod"

func init() {
	jsmod.Register("file", map[string]any{
		"exists":             Exists,
		"read":               Read,
		"readBytes":          ReadBytes,
		"readLines":          ReadLines,
		"readDir":            ReadDir,
		"getFileInfo":        GetFileInfo,
		"getAbsFilename":     GetAbsFilename,
		"write":              Write,
		"writeBytes":         WriteBytes,
		"remove":             Remove,
		"mkdir":              Mkdir,
		"copy":               Copy,
		"move":               Move,
		"replace":            Replace,
		"unmarshalFile":      UnmarshalFile,
		"marshalFile":        MarshalFile,
		"marshalFilePretty":  MarshalFilePretty,
		"patchFile":          PatchFile,
		"archive":            Archive,
		"extract":            Extract,
		"compress":           Compress,
		"decompress":         Decompress,
		"readFileFromMemory": ReadFileFromMemory,
		"readDirFromMemory":  ReadDirFromMemory,
	},
		"exists", "read", "readBytes", "readLines", "readDir", "getFileInfo", "getAbsFilename",
		"write", "writeBytes", "remove", "mkdir", "copy", "move", "replace",
		"unmarshalFile", "marshalFile", "marshalFilePretty", "patchFile",
		"archive", "extract", "compress", "decompress",
		"readFileFromMemory", "readDirFromMemory",
	)
}
