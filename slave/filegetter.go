package main

type filegetter struct {
}

func (filegetter) Get(k string) ([]byte, error) {
	return []byte(k), nil
	return nil, nil
}
