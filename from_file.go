package adata

import (
	"github.com/skypbc/goutils/gerrors"
	"github.com/skypbc/goutils/gfiles"
	"io"
	"io/fs"
	"os"
)

func NewDataFromFile(filepath string) IDataFromFile {
	return &dataFromFile{
		filepath: filepath,
	}
}

type dataFromFile struct {
	filepath string
	stat     fs.FileInfo
	fd       *os.File
}

func (d *dataFromFile) Open() (err error) {
	if d.fd != nil {
		return gerrors.NewUnknownError().
			SetTemplate("data is already used")
	}
	fromStat, err := os.Stat(d.filepath)
	if err != nil {
		return err
	}
	if !fromStat.Mode().IsRegular() {
		return gerrors.NewIncorrectParamsError().
			SetTemplate(`The "{path}" file isn't regular`).
			AddStr("path", d.filepath)
	}
	if d.fd, err = os.Open(d.filepath); err != nil {
		return err
	}
	return nil
}

func (d *dataFromFile) IsOpen() bool {
	return d.fd != nil
}

func (d *dataFromFile) Read(buff []byte) (size int, err error) {
	n, err := d.Read64(buff)
	return int(n), err
}

func (d *dataFromFile) Read64(buff []byte) (size int64, err error) {
	if d.fd == nil {
		return 0, gerrors.NewUnknownError().
			SetTemplate("data isn't initialized")
	}
	s, e := d.fd.Read(buff)
	return int64(s), e
}

func (d *dataFromFile) ReadAll() (data []byte, err error) {
	if d.fd == nil {
		if err = d.Open(); err != nil {
			return nil, gerrors.Wrap(err)
		}
		defer d.Close()
	}
	size, err := d.Size()
	if err != nil {
		return nil, err
	}
	data = make([]byte, size)
	n, err := io.ReadFull(d.fd, data)
	if err != nil {
		return nil, err
	}
	return data[:n], err
}

func (d *dataFromFile) Path() string {
	return d.filepath
}

func (d *dataFromFile) Size() (size int64, err error) {
	if d.stat == nil {
		// Если файл не открыт
		if d.fd == nil {
			return gfiles.Size2(d.filepath)
		}
		if d.stat, err = d.fd.Stat(); err != nil {
			return 0, err
		}
		return d.stat.Size(), nil
	}
	return d.stat.Size(), nil
}

func (d *dataFromFile) Seek(offset int64, whence int) (ret int64, err error) {
	if d.fd == nil {
		return 0, gerrors.NewUnknownError().
			SetTemplate("data isn't initialized")
	}
	return d.fd.Seek(offset, io.SeekStart)
}

func (d *dataFromFile) Tell() (int64, error) {
	if d.fd == nil {
		return 0, gerrors.NewUnknownError().
			SetTemplate("data isn't initialized")
	}
	offset, err := d.fd.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}
	_, err = d.fd.Seek(offset, io.SeekStart)
	return offset, err
}

func (d *dataFromFile) Close() (err error) {
	if d.fd != nil {
		err = d.fd.Close()
		d.fd = nil
		d.stat = nil
	}
	return err
}

func (d *dataFromFile) IsClosed() bool {
	return d.fd == nil
}
