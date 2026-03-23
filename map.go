package crypton

import "sync"

// customMap представляет собой мапу с защитой от состояния гонки.
type customMap struct {
	mu        sync.Mutex
	m         map[int]int
	accesses  int // счётчик обращений к ключам
	additions int // счётчик добавлений новых ключей
}
