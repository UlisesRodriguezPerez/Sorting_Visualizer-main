package algorithms

import (
	"math/rand"
	"strconv"
	"time"
)

var closed bool

// QuickSort ...
func QuickSort(data []int, c chan [][]int, stopCh chan struct{}, msgCh chan string) {
	t := time.Now()
	swaps := 0
	comparations := 0
	loops := 0

	closed = false
	quickSortAux(data, c, stopCh, msgCh, &swaps, &comparations, &loops)
	if !closed {
		hi, mi, si := t.Clock()
		hf, mf, sf := time.Now().Clock()
		msgCh <- "\nQuickSort" + "\n  Tiempo inicio = " + strconv.Itoa(hi) + ":" + strconv.Itoa(mi) + ":" + strconv.Itoa(si) + "\n  Tiempo final = " + strconv.Itoa(hf) + ":" + strconv.Itoa(mf) + ":" + strconv.Itoa(sf) + "\n  Tiempo total = " + time.Since(t).String() + "\n  Intercambio de valores = " + strconv.Itoa(swaps) + "\n  Comparación entre valores = " + strconv.Itoa(comparations) + "\n  Condición de un ciclo = " + strconv.Itoa(loops)
		close(c)
	}
}

func quickSortAux(a []int, c chan [][]int, stopCh chan struct{}, msgCh chan string, swaps, comparations, loops *int) []int {
	if len(a) < 2 {
		return a
	}

	left, right := 0, len(a)-1

	pivot := rand.Int() % len(a)

	a[pivot], a[right] = a[right], a[pivot]

	for i := range a {
		*loops++
		*comparations++
		if a[i] < a[right] {
			*swaps++
			a[left], a[i] = a[i], a[left]
			left++
		}
	}

	*swaps++
	a[left], a[right] = a[right], a[left]

	quickSortAux(a[:left], c, stopCh, msgCh, swaps, comparations, loops)
	quickSortAux(a[left+1:], c, stopCh, msgCh, swaps, comparations, loops)

	var listIndex [][]int
	for i := 0; i < len(a); i++ {
		listIndex = append(listIndex, []int{a[i], i})
	}
	if !closed {
		select {
		case <-stopCh:
			close(c)
			closed = true
			return a
		case c <- listIndex:
		}
	}
	return a
}
