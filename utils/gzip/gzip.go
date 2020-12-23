package gzip

/**
 * @DateTime   : 2020/12/22
 * @Author     : xumamba
 * @Description:
 **/

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"sync"
)

type gZip struct {
	w *gzip.Writer
	b *bytes.Buffer
}

var gZipPool = sync.Pool{
	New: func() interface{} {
		b := new(bytes.Buffer)
		return &gZip{
			w: gzip.NewWriter(b),
			b: b,
		}
	}}

func getGZip() *gZip {
	return gZipPool.Get().(*gZip)
}

func putGZip(gz *gZip) {
	gz.b.Reset()
	gz.w.Reset(gz.b)
	gZipPool.Put(gz)
}

func GZip(data []byte) ([]byte, error) {
	gz := getGZip()
	defer putGZip(gz)

	_, err := gz.w.Write(data)
	if err != nil {
		return nil, err
	}
	err = gz.w.Close()
	if err != nil {
		return nil, err
	}
	result := make([]byte, gz.b.Len())
	copy(result, gz.b.Bytes())
	return result, nil
}

func UnGZip(data []byte) ([]byte, error) {
	reader := bytes.NewReader(data)
	r, err := gzip.NewReader(reader)
	defer r.Close()
	if err != nil {
		return nil, err
	}
	all, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return all, nil
}
