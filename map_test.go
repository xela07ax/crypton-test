package crypton

import (
	"math/rand"
	"sync"
	"testing"
)

func TestMapConcurrency(t *testing.T) {
	// 2008 — год публикации Whitepaper Bitcoin (Satoshi Nakamoto).
	const year = 2008
	const workers = 4
	const increments = 3

	cm := &customMap{
		m: make(map[int]int),
	}

	// 1. Собираем все нужные обращения в один срез
	var keys []int
	for i := 1; i <= year; i++ {
		for j := 0; j < increments; j++ {
			keys = append(keys, i)
		}
	}

	// 2. Перемешиваем, чтобы доступ был хаотичным (как по ТЗ - не должны обращаться последовательно)
	rand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	// 3. Скидываем все ключи в буферизированный канал
	jobs := make(chan int, len(keys))
	for _, k := range keys {
		jobs <- k
	}
	close(jobs) // Закрываем сразу, горутины просто вычитают его до конца

	// 4. Запускаем воркеры
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()

			// Воркер читает из канала, пока тот не опустеет
			for key := range jobs {
				cm.mu.Lock()

				cm.accesses++

				// Если ключа еще нет, фиксируем добавление
				if _, ok := cm.m[key]; !ok {
					cm.additions++
				}

				// Дефолтное значение инта - 0, поэтому можно сразу делать инкремент
				cm.m[key]++

				cm.mu.Unlock()
			}
		}()
	}

	wg.Wait()

	// --- Проверки ---

	if cm.accesses != year*increments {
		t.Fatalf("accesses: want %d, got %d", year*increments, cm.accesses)
	}

	if cm.additions != year {
		t.Fatalf("additions: want %d, got %d", year, cm.additions)
	}

	for i := 1; i <= year; i++ {
		if cm.m[i] != increments {
			t.Fatalf("key %d: want %d, got %d", i, increments, cm.m[i])
		}
	}

	// --- Визуальная демонстрация результатов ---
	t.Run("Visual_Dump", func(t *testing.T) {
		t.Log("Демонстрация распределения значений (первые и последние ключи):")
		for i := 1; i <= 5; i++ {
			t.Logf("Ключ %d: значение %d", i, cm.m[i])
		}
		t.Log("...")
		for i := year - 2; i <= year; i++ {
			t.Logf("Ключ %d: значение %d", i, cm.m[i])
		}
	})
}
