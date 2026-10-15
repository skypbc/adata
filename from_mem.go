package adata

import (
	"io"
	"math"

	"github.com/skypbc/goutils/gerrors"
)

func NewDataFromMem(data []byte) IDataFromMem {
	return &dataFromMem{
		data:   data,
		offset: 0,
	}
}

type dataFromMem struct {
	data   []byte
	offset int
}

func (d *dataFromMem) Path() string {
	return ""
}

func (d *dataFromMem) Read(buff []byte) (size int, err error) {
	ds := len(d.data)
	if d.offset >= ds {
		return 0, io.EOF
	}

	size = copy(buff, d.data[d.offset:])
	d.offset += size

	if d.offset >= ds {
		d.offset = 0
		return size, io.EOF
	}
	return size, nil
}

func (d *dataFromMem) Raw() []byte {
	return d.data
}

func (d *dataFromMem) Seek(offset int64, whence int) (int64, error) {
	var base int64

	switch whence {
	case io.SeekStart:
		base = 0

	case io.SeekCurrent:
		base = int64(d.offset)

	case io.SeekEnd:
		base = int64(len(d.data))

	default:
		return 0, gerrors.NewIncorrectParamsError().
			SetTemplate(`invalid whence: "{whence}"`).
			AddInt("whence", whence)
	}

	// Проверяем переполнение при сложении
	if offset > 0 && base > math.MaxInt64-offset {
		return 0, gerrors.NewIncorrectParamsError().
			SetTemplate("seek position overflow")
	}
	if offset < 0 && base < math.MinInt64-offset {
		return 0, gerrors.NewIncorrectParamsError().
			SetTemplate("seek position overflow")
	}

	position := base + offset

	// Позиция перед началом данных недопустима.
	if position < 0 {
		return 0, gerrors.NewIncorrectParamsError().
			SetTemplate(`negative seek position: "{position}"`).
			AddInt64("position", position)
	}

	// Проверяем, что значение помещается в int.
	if uint64(position) > uint64(^uint(0)>>1) {
		return 0, gerrors.NewIncorrectParamsError().
			SetTemplate("seek position does not fit into int")
	}

	d.offset = int(position)

	return position, nil
}

func (d *dataFromMem) Tell() (offset int64, err error) {
	return int64(d.offset), nil
}

func (d *dataFromMem) Size() (size int64, err error) {
	return int64(len(d.data)), nil
}
