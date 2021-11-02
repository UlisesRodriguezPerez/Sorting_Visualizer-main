package algorithms

import (
	"strconv"
	"time"
)

// SelectionSort ...
func SelectionSort(arr []int, c chan [][]int, stopCh chan struct{}, msgCh chan string) {
	swaps := 0
	comparisons := 0
	loops := 0
	t := time.Now()
	len := len(arr)
	for i := 0; i < len-1; i++ {
		minIndex := i
		for j := i + 1; j < len; j++ {
			loops++
			comparisons++
			if arr[j] <= arr[minIndex] {
				arr[j], arr[minIndex] = arr[minIndex], arr[j]
				swaps++
				select {
				case <-stopCh:
					close(c)
					return
				case c <- [][]int{[]int{arr[j], j}, []int{arr[minIndex], minIndex}}:
				}
			}
		}
	}
	hi, mi, si := t.Clock()
	hf, mf, sf := time.Now().Clock()
	msgCh <- "\nSelectionSort:" + "\n  Tiempo inicio = " + strconv.Itoa(hi) + ":" + strconv.Itoa(mi) + ":" + strconv.Itoa(si) + "\n  Tiempo final = " + strconv.Itoa(hf) + ":" + strconv.Itoa(mf) + ":" + strconv.Itoa(sf) + "\n  Tiempo total = " + time.Since(t).String() + "\n  Intercambio de valores = " + strconv.Itoa(swaps) + "\n  Comparación entre valores = " + strconv.Itoa(comparisons) + "\n  Condición de un ciclo = " + strconv.Itoa(loops)
	close(c)
}
