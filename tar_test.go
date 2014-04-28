package tar

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"log"
	//"os"
	"testing"
)

func Test_tar(te *testing.T) {
	byt, err := TarByte(`/home/xarch/golang/backup/`)
	if err != nil {
		te.Fatal(err)
	}
	r := bytes.NewReader(byt)
	tr := tar.NewReader(r)

	// Iterate through the files in the archive.
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("File %s:\n", hdr.Name)
	}
}
