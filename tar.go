package reachtar

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type TarArchive struct {
	buffer    *bytes.Buffer
	tarWriter *tar.Writer
}

func NewArchive() *TarArchive {
	t := &TarArchive{}
	t.buffer = bytes.NewBuffer(nil)
	t.tarWriter = tar.NewWriter(t.buffer)
	return t
}

func (this *TarArchive) ArchiveFile(src string, path string, fi os.FileInfo) error {
	hdr, err := tar.FileInfoHeader(fi, path) //what should replace the path?
	if err != nil {
		return fmt.Errorf("convert file header error: %v", err)
	}
	relname, err := filepath.Rel(filepath.Clean(src+"/../"), filepath.Clean(path)) //get the rel path
	if err != nil {
		return fmt.Errorf("transform relpath error: %v", err)
	}
	hdr.Name = relname
	path = filepath.Clean(path)
	err = this.tarWriter.WriteHeader(hdr)
	if err != nil {
		return fmt.Errorf("write file header error: %v", err)
	}

	if !fi.IsDir() {
		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("read file \"%v\" error: %v", path, err)
		}
		defer f.Close()
		_, err = io.Copy(this.tarWriter, f)
		if err != nil {
			return err
		}
	}

	err = this.tarWriter.Flush()
	if err != nil {
		return fmt.Errorf("flush buffer error: %v", err)
	}
	return nil
}
func TarByte(src string) ([]byte, error) {
	a := NewArchive()
	src = filepath.Clean(src)

	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk file error: %v", err)
		}
		err = a.ArchiveFile(src, path, info)
		if err != nil {
			return fmt.Errorf("archive file \"%v\" error: %v", path, err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	a.tarWriter.Close()
	return a.buffer.Bytes(), nil
}

func Tarit(src string, target string) error {
	byt, err := TarByte(src)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	f.Write(byt)
	return nil
}
