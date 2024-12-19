package list_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/datngo2sgtech/go-packages/list"
)

func TestChunk(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name   string
		input  []int
		size   int
		expect [][]int
	}{
		{
			name:   "Empty Input",
			input:  []int{},
			size:   1,
			expect: [][]int{},
		},
		{
			name:  "Normal Input, size is 1",
			input: []int{1, 2, 3, 4, 5, 6},
			size:  1,
			expect: [][]int{
				{1},
				{2},
				{3},
				{4},
				{5},
				{6},
			},
		},
		{
			name:  "Normal Input, len divisible by size",
			input: []int{1, 2, 3, 4, 5, 6},
			size:  3,
			expect: [][]int{
				{1, 2, 3},
				{4, 5, 6},
			},
		},
		{
			name:  "Normal Input, len not divisible by size",
			input: []int{1, 2, 3, 4, 5, 6},
			size:  4,
			expect: [][]int{
				{1, 2, 3, 4},
				{5, 6},
			},
		},
		{
			name:  "Normal Input, len is equal size",
			input: []int{1, 2, 3, 4, 5, 6},
			size:  6,
			expect: [][]int{
				{1, 2, 3, 4, 5, 6},
			},
		},
		{
			name:  "Normal Input, len is greater than size",
			input: []int{1, 2, 3, 4, 5, 6},
			size:  7,
			expect: [][]int{
				{1, 2, 3, 4, 5, 6},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actualResult := list.Chunk(tc.input, tc.size)
			fmt.Println(tc.name)
			fmt.Println(actualResult)
			fmt.Println(tc.expect)
			fmt.Println(reflect.DeepEqual(actualResult, tc.expect))
			assert.True(t, reflect.DeepEqual(actualResult, tc.expect))
		})
	}
}
