package common

import "fmt"

//returns the marketCap in M of the given dic {symbol, description}
func val(m map[string]string, keys map[string]int) int {
	if i, ok := keys[m["symbol"]]; !ok {
		return -1
	} else {
		return i
	}
}

func maxHeapify(a []map[string]string, heapSize, i int, keys map[string]int) {
	l := i * 2
	r := l + 1

	largest := i

	if l < heapSize && val(a[l], keys) > val(a[i], keys) {
		largest = l
	}

	if r < heapSize && val(a[r], keys) > val(a[largest], keys) {
		largest = r
	}

	if largest != i {
		a[i], a[largest] = a[largest], a[i]
		maxHeapify(a, heapSize, largest, keys)
	}
}

func buildMaxHeap(a []map[string]string, keys map[string]int) {
	heapSize := len(a)

	for i := heapSize / 2; i >= 1; i-- {
		maxHeapify(a, heapSize, i, keys)
	}
}

func HeapSort(a []map[string]string, keys map[string]int) {
	a = append([]map[string]string{nil}, a...)
	buildMaxHeap(a, keys)

	for i := len(a) - 1; i >= 2; i-- {

		a[i], a[1] = a[1], a[i]

		maxHeapify(a, i, 1, keys)
	}

	a = a[1:]
	fmt.Println(a)
}
