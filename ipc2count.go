package ipc2count

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"iter"
	"os"

	ai "github.com/apache/arrow-go/v18/arrow/ipc"
)

var ErrIpcFilenameToCnt error = errors.New("unable to get the num records of the file")

type IpcFileReader struct{ *ai.FileReader }

func (f IpcFileReader) Close() error    { return f.FileReader.Close() }
func (f IpcFileReader) NumRecords() int { return f.FileReader.NumRecords() }

type IpcFileLike struct{ ai.ReadAtSeeker }

func (f IpcFileLike) ToReader(opts ...ai.Option) (IpcFileReader, error) {
	frdr, e := ai.NewFileReader(f.ReadAtSeeker, opts...)
	return IpcFileReader{FileReader: frdr}, e
}

func (f IpcFileLike) ToNumRecords(opts ...ai.Option) (int, error) {
	ifr, e := f.ToReader(opts...)
	if nil != e {
		return 0, e
	}

	return ifr.NumRecords(), ifr.Close()
}

type IpcFileName string

func (n IpcFileName) ToNumRecords(opts ...ai.Option) (int, error) {
	ofile, e := os.Open(string(n))
	if nil != e {
		return 0, e
	}
	defer ofile.Close() //nolint:errcheck

	ifl := IpcFileLike{ReadAtSeeker: ofile}
	return ifl.ToNumRecords(opts...)
}

type IpcFilenames iter.Seq[string]

func (i IpcFilenames) ToTotalNumRecords(opts ...ai.Option) (int, error) {
	var tot int = 0

	for fname := range i {
		iname := IpcFileName(fname)
		cnt, e := iname.ToNumRecords(opts...)
		if nil != e {
			return 0, fmt.Errorf("%w: %s", ErrIpcFilenameToCnt, fname)
		}
		tot += cnt
	}

	return tot, nil
}

type FilenameReader struct{ io.Reader }

func (r FilenameReader) ToFilenames() iter.Seq[string] {
	return func(yield func(filename string) bool) {
		var s *bufio.Scanner = bufio.NewScanner(r.Reader)
		for s.Scan() {
			var filename string = s.Text()
			if !yield(filename) {
				return
			}
		}
	}
}

func (r FilenameReader) ToTotalNumRecords(opts ...ai.Option) (int, error) {
	return IpcFilenames(r.ToFilenames()).ToTotalNumRecords(opts...)
}

var FilenameReaderStdin FilenameReader = FilenameReader{Reader: os.Stdin}

type IpcFileReadOpts struct{ Opts []ai.Option }

func (o IpcFileReadOpts) FilenamesToTotalCount(i FilenameReader) (int, error) {
	return i.ToTotalNumRecords(o.Opts...)
}

func (o IpcFileReadOpts) StdinToNamesToTotalCount() (int, error) {
	return o.FilenamesToTotalCount(FilenameReaderStdin)
}

var IpcFileReadOptsEmpty IpcFileReadOpts

func StdinToIpcNamesToTotalCountDefault() (int, error) {
	return IpcFileReadOptsEmpty.StdinToNamesToTotalCount()
}
