/*
 * Copyright (c) 2020. App Nerds LLC All Rights Reserved
 */

package minifyfs

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"path/filepath"
	"strings"
)

/*
Dir implements http.FileSystem. This is primarily designed to
be used for scenarios where you wish to embed and compile
static assets into your Go binary. The files listed using
this file system are minified if they are CSS or JavaScript.

I wouldn't recommend using this directly when serving from
an HTTP router, as the File methods here aren't terribly
efficient. In fact, I bet there are better ways to do this.
*/
type Dir string

func (d Dir) Open(name string) (http.File, error) {
	var err error
	var f *File

	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return nil, errors.New("http: invalid character in file path")
	}

	dir := string(d)
	if dir == "" {
		dir = "."
	}

	fullName := filepath.Join(dir, filepath.FromSlash(path.Clean("/"+name)))

	if f, err = Open(fullName); err != nil {
		return nil, fmt.Errorf("Error opening file '%s': %w", name, err)
	}

	return f, nil
}
