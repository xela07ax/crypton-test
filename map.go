package crypton

import "sync"

// customMap представляет собой мапу с защитой от состояния гонки.
// По условиям задачи код не претендует на production-ready.
type customMap struct {
	mu        sync.Mutex
	m         map[int]int
	accesses  int
	additions int
}
