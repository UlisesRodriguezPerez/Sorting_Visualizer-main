package algorithms

import (
	"strconv"
	"time"
)

//BubbleSort ...
func BubbleSort(data []int, c chan [][]int, stopCh chan struct{}, msgCh chan string) {
	swaps := 0
	comparisons := 0
	loops := 0
	t := time.Now()

	for i := 0; i < len(data); i++ {
		for j := 1; j < len(data)-i; j++ {
			loops++
			comparisons++
			if data[j] < data[j-1] {
				data[j], data[j-1] = data[j-1], data[j]
				select {
				case <-stopCh:
					close(c)
					return
				case c <- [][]int{[]int{data[j], j}, []int{data[j-1], j - 1}}:
				}
				swaps++
			}
		}
	}
	hi, mi, si := t.Clock()
	hf, mf, sf := time.Now().Clock()
	msgCh <- "\nBubbleSort:" + "\n  Tiempo inicio = " + strconv.Itoa(hi) + ":" + strconv.Itoa(mi) + ":" + strconv.Itoa(si) + "\n  Tiempo final = " + strconv.Itoa(hf) + ":" + strconv.Itoa(mf) + ":" + strconv.Itoa(sf) + "\n  Tiempo total = " + time.Since(t).String() + "\n  Intercambio de valores = " + strconv.Itoa(swaps) + "\n  Comparación entre valores = " + strconv.Itoa(comparisons) + "\n  Condición de un ciclo = " + strconv.Itoa(loops)

	close(c)
}
