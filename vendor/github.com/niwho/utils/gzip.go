package utils

import (
	"compress/gzip"
	"io"
	"io/ioutil"
)

func Ungzip(r io.Reader) (io.ReadCloser, error) {
	rb, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return rb, nil
}

func Ungzip2Bytes(r io.Reader) ([]byte, error) {
	rb, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	dat, err := ioutil.ReadAll(rb)
	return dat, err
}
