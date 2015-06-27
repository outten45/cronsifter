package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"time"
)

// Simple logger for writing to a file. This is heavily
// based on https://github.com/siddontang/go-log.

// SimpleLogger is a basic logging struct that provide log rotating.
type SimpleLogger struct {
	file      *os.File
	filename  string
	maxBytes  int
	fileCount int
}

// NewSimpleLogger takes the filename, maxBytes and fileCount to create a
// SimpleLogger. maxBytes and fileCount are used to determine when to rotate
// the files.
func NewSimpleLogger(filename string, maxBytes, fileCount int) (*SimpleLogger, error) {
	dir := path.Dir(filename)
	os.Mkdir(dir, 0777)

	sl := &SimpleLogger{filename: filename, maxBytes: maxBytes, fileCount: fileCount}

	var err error
	sl.file, err = os.OpenFile(sl.filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return sl, nil
}

// Write the given bytes to the file.
func (l *SimpleLogger) Write(b []byte) (int, error) {
	l.rotate()
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s: ", time.Now().Format("Mon Jan 2 15:04:05 -0700 MST 2006")))
	buf.Write(b)
	buf.WriteString("\n")
	return l.file.Write(buf.Bytes())
}

// Close the SimpleLogger's file handle.
func (l *SimpleLogger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

func (l *SimpleLogger) rotate() error {
	f, err := l.file.Stat()
	if err != nil {
		return err
	}

	if l.maxBytes <= 0 {
		return fmt.Errorf("Max Bytes needs to be greater than 0 instead of: %d", l.maxBytes)
	}
	if f.Size() < int64(l.maxBytes) {
		return nil
	}

	if l.fileCount > 0 {
		l.file.Close()

		for i := l.fileCount - 1; i > 0; i-- {
			prev := fmt.Sprintf("%s.%d", l.filename, i)
			next := fmt.Sprintf("%s.%d", l.filename, i+1)

			os.Rename(prev, next)
		}

		next := fmt.Sprintf("%s.1", l.filename)
		os.Rename(l.filename, next)

		l.file, _ = os.OpenFile(l.filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	}

	return nil
}
