package file

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"compress/zlib"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Compress(data []byte, cType string) ([]byte, error) {
	var buf bytes.Buffer
	var w io.WriteCloser

	switch strings.ToLower(cType) {
	case "gzip", "gz":
		w = gzip.NewWriter(&buf)
	default:
		w = zlib.NewWriter(&buf)
	}

	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	w.Close()
	return buf.Bytes(), nil
}

func Decompress(data []byte, cType string) ([]byte, error) {
	bufR := bytes.NewReader(data)
	var r io.ReadCloser
	var err error

	switch strings.ToLower(cType) {
	case "gzip", "gz":
		r, err = gzip.NewReader(bufR)
	default:
		r, err = zlib.NewReader(bufR)
	}

	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}

func MustGzip(data []byte) []byte   { b, _ := Compress(data, "gzip"); return b }
func MustGunzip(data []byte) []byte { b, _ := Decompress(data, "gzip"); return b }
func MustZip(data []byte) []byte    { b, _ := Compress(data, "zlib"); return b }
func MustUnzip(data []byte) []byte  { b, _ := Decompress(data, "zlib"); return b }

func Extract(srcFile, destDir string, stripRoot bool) error {
	f, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer f.Close()

	lowerSrc := strings.ToLower(srcFile)
	switch {
	case strings.HasSuffix(lowerSrc, ".zip"):
		return extractZip(srcFile, destDir, stripRoot)
	case strings.HasSuffix(lowerSrc, ".tar.gz"), strings.HasSuffix(lowerSrc, ".tgz"):
		gzr, err := gzip.NewReader(f)
		if err != nil {
			return err
		}
		defer gzr.Close()
		return extractTar(gzr, destDir, stripRoot)
	case strings.HasSuffix(lowerSrc, ".tar"):
		return extractTar(f, destDir, stripRoot)
	case strings.HasSuffix(lowerSrc, ".gz"):
		gzr, err := gzip.NewReader(f)
		if err != nil {
			return err
		}
		defer gzr.Close()
		return extractSingleFile(gzr, destDir, strings.TrimSuffix(filepath.Base(srcFile), ".gz"))
	case strings.HasSuffix(lowerSrc, ".bz2"):
		bzr := bzip2.NewReader(f)
		return extractSingleFile(bzr, destDir, strings.TrimSuffix(filepath.Base(srcFile), ".bz2"))
	default:
		return extractZip(srcFile, destDir, stripRoot)
	}
}

func extractTar(r io.Reader, dest string, strip bool) error {
	tr := tar.NewReader(r)
	for {
		h, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if err := writeEntry(dest, h.Name, h.FileInfo().Mode(), h.Typeflag == tar.TypeDir, h.Linkname, tr, strip); err != nil {
			return err
		}
	}
}

func extractZip(src, dest string, strip bool) error {
	rz, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer rz.Close()
	for _, f := range rz.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		var linkTarget string
		if f.Mode()&os.ModeSymlink != 0 {
			b, _ := io.ReadAll(rc)
			linkTarget = string(b)
		}
		err = writeEntry(dest, f.Name, f.Mode(), f.FileInfo().IsDir(), linkTarget, rc, strip)
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func extractSingleFile(r io.Reader, destDir, fileName string) error {
	os.MkdirAll(destDir, 0755)
	target := filepath.Join(destDir, fileName)
	out, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, r)
	return err
}

func writeEntry(destDir, name string, mode os.FileMode, isDir bool, linkPath string, r io.Reader, strip bool) error {
	if strip {
		parts := strings.SplitN(name, "/", 2)
		if len(parts) < 2 || parts[1] == "" {
			return nil
		}
		name = parts[1]
	}
	target := filepath.Join(destDir, name)
	if isDir {
		return os.MkdirAll(target, 0755)
	}
	os.MkdirAll(filepath.Dir(target), 0755)
	if mode&os.ModeSymlink != 0 {
		os.Remove(target)
		return os.Symlink(linkPath, target)
	}
	out, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, r)
	return err
}

func Archive(srcPath, destFile string) error {
	f, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer f.Close()

	lower := strings.ToLower(destFile)
	if strings.HasSuffix(lower, ".zip") {
		zw := zip.NewWriter(f)
		defer zw.Close()
		return walkAndAdd(srcPath, func(relPath string, info os.FileInfo, fileReader io.Reader) error {
			header, _ := zip.FileInfoHeader(info)
			header.Name = relPath
			if info.IsDir() {
				header.Name += "/"
			} else {
				header.Method = zip.Deflate
			}
			w, err := zw.CreateHeader(header)
			if err != nil || info.IsDir() {
				return err
			}
			_, err = io.Copy(w, fileReader)
			return err
		})
	}
	gw := gzip.NewWriter(f)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	return walkAndAdd(srcPath, func(relPath string, info os.FileInfo, fileReader io.Reader) error {
		header, _ := tar.FileInfoHeader(info, "")
		header.Name = relPath
		if info.Mode()&os.ModeSymlink != 0 {
			link, _ := os.Readlink(filepath.Join(srcPath, "..", relPath))
			header.Linkname = link
		}
		if err := tw.WriteHeader(header); err != nil || info.IsDir() || header.Typeflag == tar.TypeSymlink {
			return err
		}
		_, err = io.Copy(tw, fileReader)
		return err
	})
}

func walkAndAdd(srcPath string, addFn func(string, os.FileInfo, io.Reader) error) error {
	baseDir := filepath.Dir(srcPath)
	return filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, _ := filepath.Rel(baseDir, path)
		if relPath == "." {
			return nil
		}
		var r io.Reader
		if !info.IsDir() && info.Mode()&os.ModeSymlink == 0 {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			r = f
		}
		return addFn(filepath.ToSlash(relPath), info, r)
	})
}
