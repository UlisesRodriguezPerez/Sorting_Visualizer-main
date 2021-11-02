package algorithms

import (
	"strconv"
	"time"
)

//MergeSort initialize the sorting
func MergeSort(data []int, c chan [][]int, stopCh chan struct{}, msgCh chan string) {
	t := time.Now()
	swaps := 0
	comparisons := 0
	loops := 0

	closed = false
	mergeSortAux(data, c, stopCh, msgCh, &swaps, &comparisons, &loops)
	if !closed {
		hi, mi, si := t.Clock()
		hf, mf, sf := time.Now().Clock()
		msgCh <- "\nMergeSort:" + "\n  Tiempo inicio = " + strconv.Itoa(hi) + ":" + strconv.Itoa(mi) + ":" + strconv.Itoa(si) + "\n  Tiempo final = " + strconv.Itoa(hf) + ":" + strconv.Itoa(mf) + ":" + strconv.Itoa(sf) + "\n  Tiempo total = " + time.Since(t).String() + "\n  Intercambio de valores = " + strconv.Itoa(swaps) + "\n  Comparación entre valores = " + strconv.Itoa(comparisons) + "\n  Condición de un ciclo = " + strconv.Itoa(loops)
		close(c)
	}
}

func mergeSortAux(data []int, c chan [][]int, stopCh chan struct{}, msgCh chan string, swaps, comparisons, loops *int) []int {

	var num = len(data)

	if num == 1 {
		return data
	}

	middle := int(num / 2)
	var (
		left  = make([]int, middle)
		right = make([]int, num-middle)
	)
	for i := 0; i < num; i++ {
		*loops++
		*comparisons++
		if i < middle {
			left[i] = data[i]
		} else {
			right[i-middle] = data[i]
		}
		*swaps++
	}
	result := merge(mergeSortAux(left, c, stopCh, msgCh, swaps, comparisons, loops), mergeSortAux(right, c, stopCh, msgCh, swaps, comparisons, loops), swaps, comparisons, loops)

	var listIndex [][]int
	for i := 0; i < len(result); i++ {
		listIndex = append(listIndex, []int{result[i], i})
	}

	if !closed {
		select {
		case <-stopCh:
			close(c)
			closed = true
			return data
		case c <- listIndex:
		}
	}
	return result
}

func merge(left, right []int, swaps, comparisons, loops *int) (result []int) {
	result = make([]int, len(left)+len(right))

	i := 0
	for len(left) > 0 && len(right) > 0 {
		*loops++
		*comparisons++
		if left[0] < right[0] {
			result[i] = left[0]
			left = left[1:]
		} else {
			result[i] = right[0]
			right = right[1:]
		}
		*swaps++
		i++
	}

	for j := 0; j < len(left); j++ {
		*loops++
		result[i] = left[j]
		*swaps++
		i++
	}
	for j := 0; j < len(right); j++ {
		*loops++
		result[i] = right[j]
		*swaps++
		i++
	}
	return
}
