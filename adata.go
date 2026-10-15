// Amorphous Data
package adata

type IData interface {
	Read(buff []byte) (size int, err error)

	// Метод из IDataFromFile (добавлен в общий интерфейс для удобства). Его отсутствие производит много
	// однотипного кода (приведений) для доступа к нему.
	//
	// Для IDataFromMem метод всегда возвращает пустую строку.
	Path() string

	Seek(offset int64, whence int) (ret int64, err error)
	Tell() (offset int64, err error)
	Size() (size int64, err error)
}

type IDataFromMem interface {
	IData

	Raw() []byte
}

type IDataFromFile interface {
	IData

	Open() error
	IsOpen() bool

	Close() error
	IsClosed() bool
}
