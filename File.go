/*
 * Copyright (c) 2020. App Nerds LLC All Rights Reserved
 */

package minifyfs

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/js"
)

/*
Implementation of http.File, os.FileInfo, io.Reader, io.Seeker.
When you open a file using the Open method, the contents are read
whole, then minified if they are CSS or JavaScript.
*/
type File struct {
	Contents []byte
	FileName string
	Reader   *bytes.Reader
	Stats    os.FileInfo
}

func Open(fileName string) (*File, error) {
	var err error
	var originalBytes []byte

	result := &File{
		FileName: fileName,
	}

	if result.Stats, err = os.Stat(fileName); err != nil {
		return result, fmt.Errorf("Error getting file information for '%s': %w", fileName, err)
	}

	if result.Stats.IsDir() {
		return result, nil
	}

	if originalBytes, err = ioutil.ReadFile(fileName); err != nil {
		return result, fmt.Errorf("Error opening file '%s': %w", fileName, err)
	}

	contentType := mime.TypeByExtension(filepath.Ext(fileName))

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)

	if contentType == "text/css" ||
		contentType == "application/x-javascript" ||
		contentType == "application/javascript" ||
		contentType == "application/ecmascript" ||
		contentType == "text/x-javascript" ||
		contentType == "text/javascript" ||
		contentType == "text/ecmascript" {
		if result.Contents, err = m.Bytes(contentType, originalBytes); err != nil {
			return result, fmt.Errorf("Error minifying '%s': %w", fileName, err)
		}
	} else {
		result.Contents = originalBytes
	}

	result.Reader = bytes.NewReader(result.Contents)
	return result, nil
}

func (f *File) Name() string {
	return f.FileName
}

func (f *File) Size() int64 {
	if f.IsDir() {
		return f.Stats.Size()
	} else {
		return int64(len(f.Contents))
	}
}

func (f *File) Mode() os.FileMode {
	return f.Stats.Mode()
}

func (f *File) ModTime() time.Time {
	return f.Stats.ModTime()
}

func (f *File) IsDir() bool {
	return f.Stats.IsDir()
}

func (f *File) Sys() interface{} {
	return f.Stats.Sys()
}

func (f *File) Close() error {
	return nil
}

func (f *File) Read(p []byte) (n int, err error) {
	return f.Reader.Read(p)
}

func (f *File) Readdir(count int) ([]os.FileInfo, error) {
	var err error
	var ff *os.File

	if ff, err = os.Open(f.FileName); err != nil {
		return nil, err
	}

	return ff.Readdir(count)
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	return f.Reader.Seek(offset, whence)
}

func (f *File) Stat() (os.FileInfo, error) {
	return f, nil
}
