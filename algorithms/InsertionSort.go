package algorithms

import (
	"strconv"
	"time"
)

//InsertionSort ...
func InsertionSort(items []int, c chan [][]int, stopCh chan struct{}, msgCh chan string) {
	swaps := 0
	comparisons := 0
	loops := 0
	t := time.Now()

	var n = len(items)
	for i := 1; i < n; i++ {
		j := i
		for j > 0 {
			loops++
			comparisons++
			if items[j-1] > items[j] {
				items[j-1], items[j] = items[j], items[j-1]
				swaps++
				select {
				case <-stopCh:
					close(c)
					return
				case c <- [][]int{[]int{items[j-1], j - 1}, []int{items[j], j}}:
				}
			}
			j = j - 1
		}
	}
	hi, mi, si := t.Clock()
	hf, mf, sf := time.Now().Clock()
	msgCh <- "\nInsertionSort:" + "\n  Tiempo inicio = " + strconv.Itoa(hi) + ":" + strconv.Itoa(mi) + ":" + strconv.Itoa(si) + "\n  Tiempo final = " + strconv.Itoa(hf) + ":" + strconv.Itoa(mf) + ":" + strconv.Itoa(sf) + "\n  Tiempo total = " + time.Since(t).String() + "\n  Intercambio de valores = " + strconv.Itoa(swaps) + "\n  Comparación entre valores = " + strconv.Itoa(comparisons) + "\n  Condición de un ciclo = " + strconv.Itoa(loops)

	close(c)
}
