// MIT License Copyright (C) 2022 Hiroshi Shimamoto
package pptxpack

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type PPTX struct {
	path string
}

func New(path string) (*PPTX, error) {
	info, err := os.Stat(path)
	if err != nil {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return nil, err
		}
		return &PPTX{path}, nil
	}
	if info.IsDir() {
		return &PPTX{path}, nil
	}
	return nil, fmt.Errorf("%s is not directory", path)
}

func Open(path string) (*PPTX, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return &PPTX{path}, nil
	}
	return nil, fmt.Errorf("%s is not directory", path)
}

func (p *PPTX) Unpack(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return err
	}

	z, err := zip.NewReader(f, info.Size())

	list := []string{}
	for _, zfs := range z.File {
		zpath := zfs.Name
		zf, err := zfs.Open()
		if err != nil {
			return err
		}
		dpath := filepath.Join(p.path, zpath)
		os.MkdirAll(filepath.Dir(dpath), 0755)
		df, err := os.Create(dpath)
		if err != nil {
			return err
		}
		if strings.HasSuffix(zpath, ".xml") || strings.HasSuffix(zpath, ".rels") {
			buf := new(bytes.Buffer)
			buf.ReadFrom(zf)
			xml, err := unpackXML(buf.String())
			if err != nil {
				return err
			}
			df.Write([]byte(xml))
		} else {
			io.Copy(df, zf)
		}
		df.Close()
		// keep file list
		list = append(list, zpath)
	}
	listpath := filepath.Join(p.path, "files.list")
	sort.Strings(list)
	return os.WriteFile(listpath, []byte(strings.Join(list, "\n")), 0644)
}

func (p *PPTX) Pack(path string) error {
	listpath := filepath.Join(p.path, "files.list")
	blist, err := os.ReadFile(listpath)
	if err != nil {
		return err
	}
	list := strings.Split(string(blist), "\n")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	z := zip.NewWriter(f)
	for _, zpath := range list {
		if zpath == "" {
			continue
		}
		dpath := filepath.Join(p.path, zpath)
		info, err := os.Stat(dpath)
		if err != nil {
			return err
		}
		blob, err := os.ReadFile(dpath)
		if err != nil {
			return err
		}
		if strings.HasSuffix(zpath, ".xml") || strings.HasSuffix(zpath, ".xml.rels") {
			s, err := packXML(string(blob))
			if err != nil {
				return err
			}
			blob = []byte(s)
		}
		hdr, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		hdr.Name = zpath
		hdr.Method = zip.Deflate
		hdr.Modified = time.Unix(0, 0)
		w, err := z.CreateHeader(hdr)
		if err != nil {
			return err
		}
		w.Write(blob)
	}
	// zip close here to flush directory
	z.Close()
	return nil
}
