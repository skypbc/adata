package adata

import (
	"github.com/skypbc/goutils/gerrors"
	"github.com/skypbc/goutils/gslice"
	"io"
)

func NewDataFromMem(data []byte) IDataFromMem {
	return &dataFromMem{
		data:   data,
		offset: 0,
	}
}

type dataFromMem struct {
	data   []byte
	offset int64
}

func (d *dataFromMem) Path() string {
	return ""
}

func (d *dataFromMem) Read(buff []byte) (size int, err error) {
	n, err := d.Read64(buff)
	return int(n), err
}

func (d *dataFromMem) Read64(buff []byte) (size int64, err error) {
	ds := int64(len(d.data))
	if d.offset >= ds {
		return size, io.EOF
	}

	size = int64(len(d.data[d.offset:]))
	if int64(len(buff)) < size {
		size = int64(len(buff))
	}

	copy(buff, d.data[d.offset:])
	d.offset += size

	if d.offset >= ds {
		d.offset = 0
		return size, io.EOF
	}
	return size, nil
}

func (d *dataFromMem) ReadAll() (data []byte, err error) {
	return gslice.Clone(d.data), nil
}

func (d *dataFromMem) Raw() []byte {
	return d.data
}

func (d *dataFromMem) Seek(offset int64, whence int) (ret int64, err error) {
	switch whence {
	case io.SeekStart:
		d.offset = offset
		return offset, nil
	case io.SeekEnd:
		ret = int64(len(d.data) - int(offset))
		return ret, nil
	case io.SeekCurrent:
		d.offset = d.offset + offset
		return d.offset, nil
	}
	return 0, gerrors.NewIncorrectParamsError().
		SetTemplate(`invalid whence: "{whence}"`).
		AddInt("whence", whence)
}

func (d *dataFromMem) Tell() (offset int64, err error) {
	return d.offset, nil
}

func (d *dataFromMem) Size() (size int64, err error) {
	return int64(len(d.data)), nil
}
