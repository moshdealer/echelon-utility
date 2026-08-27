package source

// Интерфейсная реализация нужна для поддержки разных способов считывания содержимого проверяемого файла
// В нашем случае их два: из конфиг-файла и из stdin

type Input struct {
	Name    string
	Content []byte
}

type Loader interface {
	Load() (*Input, error)
}
