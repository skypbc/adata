package adata

import (
	"io"

	"github.com/skypbc/goutils/gerrors"
	"github.com/skypbc/goutils/gslice"
)

// ReadAll читает все данные из IData и возвращает их в виде среза байтов. Если IData является файлом, то функция
// автоматически открывает его перед чтением и закрывает после завершения операции. Функция также сохраняет текущую
// позицию указателя данных и возвращает его обратно после чтения.
func ReadAll(data IData) (res []byte, err error) {
	if data == nil {
		return nil, gerrors.NewIncorrectParamsError().
			SetTemplate("data is nil")
	}

	if d, ok := data.(IDataFromMem); ok {
		return gslice.Clone(d.Raw()), nil
	}

	// Если data является файлом, то нужно убедиться, что он открыт
	if d, ok := data.(IDataFromFile); ok {
		if !d.IsOpen() {
			if err := d.Open(); err != nil {
				return nil, gerrors.Join(err, gerrors.NewUnknownError().
					SetTemplate(`can't open data file`))
			}
			// Если мы открыли файл, то нужно его закрыть после чтения
			defer d.Close()
		}
	}

	// Запоминаем текущую позицию указателя, чтобы вернуть его обратно после чтения
	pos, err := data.Tell()
	if err != nil {
		return nil, err
	}

	// Перемещаем указатель на начало данных, чтобы прочитать их полностью
	if _, err = data.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	defer func() {
		// Возвращаем указатель на исходную позицию
		if _, seekErr := data.Seek(pos, io.SeekStart); seekErr != nil {
			if err != nil {
				err = gerrors.Join(err, seekErr)
			} else {
				err = gerrors.Wrap(seekErr)
			}
		}
	}()

	size, err := data.Size()
	if err != nil {
		return nil, err
	}
	res = make([]byte, size)

	n, err := io.ReadFull(data, res)
	if err != nil {
		return nil, gerrors.Join(err, gerrors.NewUnknownError().
			SetTemplate(`can't read full data`))
	}

	return res[:n], err
}
