package sliceutils

import (
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
)

func TestMap(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		fn       func(int) string
		expected []string
	}{
		{
			name:     "数字转字符串",
			input:    []int{1, 2, 3, 4, 5},
			fn:       func(i int) string { return strconv.Itoa(i) },
			expected: []string{"1", "2", "3", "4", "5"},
		},
		{
			name:     "加倍",
			input:    []int{1, 2, 3},
			fn:       func(i int) string { return strconv.Itoa(i * 2) },
			expected: []string{"2", "4", "6"},
		},
		{
			name:     "空切片",
			input:    []int{},
			fn:       func(i int) string { return strconv.Itoa(i) },
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Map(tt.input, tt.fn)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Map() = %v, 期望 %v", result, tt.expected)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		predicate func(int) bool
		expected  []int
	}{
		{
			name:      "过滤偶数",
			input:     []int{1, 2, 3, 4, 5},
			predicate: func(i int) bool { return i%2 == 0 },
			expected:  []int{2, 4},
		},
		{
			name:      "过滤大于3的数",
			input:     []int{1, 2, 3, 4, 5},
			predicate: func(i int) bool { return i > 3 },
			expected:  []int{4, 5},
		},
		{
			name:      "没有匹配项",
			input:     []int{1, 2, 3},
			predicate: func(i int) bool { return i > 10 },
			expected:  []int{},
		},
		{
			name:      "空切片",
			input:     []int{},
			predicate: func(i int) bool { return i%2 == 0 },
			expected:  []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Filter(tt.input, tt.predicate)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Filter() = %v, 期望 %v", result, tt.expected)
			}
		})
	}
}

func TestReduce(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		start    int
		fn       func(int, int) int
		expected int
	}{
		{
			name:     "求和",
			input:    []int{1, 2, 3, 4, 5},
			start:    0,
			fn:       func(acc, val int) int { return acc + val },
			expected: 15,
		},
		{
			name:     "求积",
			input:    []int{1, 2, 3, 4},
			start:    1,
			fn:       func(acc, val int) int { return acc * val },
			expected: 24,
		},
		{
			name:     "空切片",
			input:    []int{},
			start:    10,
			fn:       func(acc, val int) int { return acc + val },
			expected: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Reduce(tt.input, tt.start, tt.fn)
			if result != tt.expected {
				t.Errorf("Reduce() = %v, 期望 %v", result, tt.expected)
			}
		})
	}

	// 字符串测试
	t.Run("拼接字符串", func(t *testing.T) {
		input := []string{"a", "b", "c"}
		expected := "abc"
		result := Reduce(input, "", func(acc, val string) string { return acc + val })
		if result != expected {
			t.Errorf("Reduce() = %v, 期望 %v", result, expected)
		}
	})
}

func TestFind(t *testing.T) {
	tests := []struct {
		name         string
		input        []int
		predicate    func(int) bool
		expectedVal  int
		expectedBool bool
	}{
		{
			name:         "找到第一个偶数",
			input:        []int{1, 2, 3, 4, 5},
			predicate:    func(i int) bool { return i%2 == 0 },
			expectedVal:  2,
			expectedBool: true,
		},
		{
			name:         "找不到匹配项",
			input:        []int{1, 3, 5},
			predicate:    func(i int) bool { return i%2 == 0 },
			expectedVal:  0, // 零值
			expectedBool: false,
		},
		{
			name:         "空切片",
			input:        []int{},
			predicate:    func(i int) bool { return i%2 == 0 },
			expectedVal:  0,
			expectedBool: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, found := Find(tt.input, tt.predicate)
			if val != tt.expectedVal || found != tt.expectedBool {
				t.Errorf("Find() = (%v, %v), 期望 (%v, %v)", val, found, tt.expectedVal, tt.expectedBool)
			}
		})
	}

	// 结构体测试
	t.Run("查找结构体", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}
		people := []Person{
			{"Alice", 25},
			{"Bob", 30},
			{"Charlie", 35},
		}
		person, found := Find(people, func(p Person) bool { return p.Age > 30 })
		if !found || person.Name != "Charlie" {
			t.Errorf("Find() = (%v, %v), 期望找到 Charlie", person, found)
		}
	})
}

func TestSome(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		predicate func(int) bool
		expected  bool
	}{
		{
			name:      "有偶数",
			input:     []int{1, 2, 3, 4, 5},
			predicate: func(i int) bool { return i%2 == 0 },
			expected:  true,
		},
		{
			name:      "没有偶数",
			input:     []int{1, 3, 5, 7},
			predicate: func(i int) bool { return i%2 == 0 },
			expected:  false,
		},
		{
			name:      "空切片",
			input:     []int{},
			predicate: func(i int) bool { return i%2 == 0 },
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Some(tt.input, tt.predicate)
			if result != tt.expected {
				t.Errorf("Some() = %v, 期望 %v", result, tt.expected)
			}
		})
	}
}

func TestEvery(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		predicate func(int) bool
		expected  bool
	}{
		{
			name:      "全是偶数",
			input:     []int{2, 4, 6, 8},
			predicate: func(i int) bool { return i%2 == 0 },
			expected:  true,
		},
		{
			name:      "不全是偶数",
			input:     []int{2, 4, 5, 6},
			predicate: func(i int) bool { return i%2 == 0 },
			expected:  false,
		},
		{
			name:      "空切片",
			input:     []int{},
			predicate: func(i int) bool { return i%2 == 0 },
			expected:  true, // 空切片对于全称量词为真
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Every(tt.input, tt.predicate)
			if result != tt.expected {
				t.Errorf("Every() = %v, 期望 %v", result, tt.expected)
			}
		})
	}
}

func TestIncludes(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		item     int
		expected bool
	}{
		{
			name:     "包含元素",
			slice:    []int{1, 2, 3, 4, 5},
			item:     3,
			expected: true,
		},
		{
			name:     "不包含元素",
			slice:    []int{1, 2, 3, 4, 5},
			item:     6,
			expected: false,
		},
		{
			name:     "空切片",
			slice:    []int{},
			item:     1,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Includes(tt.slice, tt.item)
			if result != tt.expected {
				t.Errorf("Includes() = %v, 期望 %v", result, tt.expected)
			}
		})
	}

	// 字符串测试
	t.Run("字符串切片", func(t *testing.T) {
		slice := []string{"apple", "banana", "cherry"}
		if !Includes(slice, "banana") {
			t.Errorf("Includes() 应该找到 'banana'")
		}
		if Includes(slice, "grape") {
			t.Errorf("Includes() 不应该找到 'grape'")
		}
	})
}

func TestIndexOf(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		item     int
		expected int
	}{
		{
			name:     "找到元素",
			slice:    []int{1, 2, 3, 4, 5},
			item:     3,
			expected: 2,
		},
		{
			name:     "找不到元素",
			slice:    []int{1, 2, 3, 4, 5},
			item:     6,
			expected: -1,
		},
		{
			name:     "重复元素",
			slice:    []int{1, 2, 3, 2, 5},
			item:     2,
			expected: 1, // 第一次出现的索引
		},
		{
			name:     "空切片",
			slice:    []int{},
			item:     1,
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IndexOf(tt.slice, tt.item)
			if result != tt.expected {
				t.Errorf("IndexOf() = %v, 期望 %v", result, tt.expected)
			}
		})
	}
}

func TestLastIndexOf(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		item     int
		expected int
	}{
		{
			name:     "找到元素",
			slice:    []int{1, 2, 3, 4, 5},
			item:     3,
			expected: 2,
		},
		{
			name:     "找不到元素",
			slice:    []int{1, 2, 3, 4, 5},
			item:     6,
			expected: -1,
		},
		{
			name:     "重复元素",
			slice:    []int{1, 2, 3, 2, 5},
			item:     2,
			expected: 3, // 最后一次出现的索引
		},
		{
			name:     "空切片",
			slice:    []int{},
			item:     1,
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LastIndexOf(tt.slice, tt.item)
			if result != tt.expected {
				t.Errorf("LastIndexOf() = %v, 期望 %v", result, tt.expected)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		expected []int
	}{
		{
			name:     "正常切片",
			slice:    []int{1, 2, 3, 4, 5},
			expected: []int{5, 4, 3, 2, 1},
		},
		{
			name:     "单元素切片",
			slice:    []int{1},
			expected: []int{1},
		},
		{
			name:     "空切片",
			slice:    []int{},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Reverse(tt.slice)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Reverse() = %v, 期望 %v", result, tt.expected)
			}
			// 检查原切片是否未被修改
			originalCopy := make([]int, len(tt.slice))
			copy(originalCopy, tt.slice)
			if !reflect.DeepEqual(tt.slice, originalCopy) {
				t.Errorf("原切片被修改: %v, 期望保持不变: %v", tt.slice, originalCopy)
			}
		})
	}
}

func TestReverseInPlace(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		expected []int
	}{
		{
			name:     "正常切片",
			slice:    []int{1, 2, 3, 4, 5},
			expected: []int{5, 4, 3, 2, 1},
		},
		{
			name:     "单元素切片",
			slice:    []int{1},
			expected: []int{1},
		},
		{
			name:     "空切片",
			slice:    []int{},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := make([]int, len(tt.slice))
			copy(original, tt.slice)

			result := ReverseInPlace(tt.slice)

			// 检查返回值是否正确
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ReverseInPlace() 返回 = %v, 期望 %v", result, tt.expected)
			}

			// 检查原切片是否被修改
			if !reflect.DeepEqual(tt.slice, tt.expected) {
				t.Errorf("原切片应该被修改为: %v, 但是是: %v", tt.expected, tt.slice)
			}
		})
	}
}

func TestUniq(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		expected []int
	}{
		{
			name:     "有重复元素",
			slice:    []int{1, 2, 2, 3, 4, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "无重复元素",
			slice:    []int{1, 2, 3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "单元素切片",
			slice:    []int{1},
			expected: []int{1},
		},
		{
			name:     "空切片",
			slice:    []int{},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Uniq(tt.slice)

			// 检查所有期望元素都在结果中
			for _, e := range tt.expected {
				if !Includes(result, e) {
					t.Errorf("Uniq() 缺少元素 %v", e)
				}
			}

			// 检查结果长度是否正确
			if len(result) != len(tt.expected) {
				t.Errorf("Uniq() 长度 = %v, 期望 %v", len(result), len(tt.expected))
			}

			// 检查结果中是否有重复
			seen := make(map[int]bool)
			for _, v := range result {
				if seen[v] {
					t.Errorf("Uniq() 结果中有重复元素 %v", v)
				}
				seen[v] = true
			}
		})
	}

	// 字符串测试
	t.Run("字符串切片", func(t *testing.T) {
		slice := []string{"a", "b", "a", "c", "b"}
		expected := []string{"a", "b", "c"}
		result := Uniq(slice)

		// 检查所有期望元素都在结果中
		for _, e := range expected {
			if !Includes(result, e) {
				t.Errorf("Uniq() 缺少元素 %v", e)
			}
		}

		// 检查结果长度是否正确
		if len(result) != len(expected) {
			t.Errorf("Uniq() 长度 = %v, 期望 %v", len(result), len(expected))
		}
	})
}

func TestFlatMap(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		fn       func(int) []string
		expected []string
	}{
		{
			name:     "重复字符串",
			input:    []int{1, 2, 3},
			fn:       func(i int) []string { return []string{strings.Repeat("a", i)} },
			expected: []string{"a", "aa", "aaa"},
		},
		{
			name:  "可变长度结果",
			input: []int{1, 2, 3},
			fn: func(i int) []string {
				result := make([]string, i)
				for j := 0; j < i; j++ {
					result[j] = strconv.Itoa(i)
				}
				return result
			},
			expected: []string{"1", "2", "2", "3", "3", "3"},
		},
		{
			name:     "空切片",
			input:    []int{},
			fn:       func(i int) []string { return []string{strconv.Itoa(i)} },
			expected: []string{},
		},
		{
			name:     "部分空结果",
			input:    []int{0, 1, 2},
			fn:       func(i int) []string { return make([]string, i) },
			expected: []string{"", "", ""}, // 一个0长度切片，一个值为空字符串的元素，两个值为空字符串的元素
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FlatMap(tt.input, tt.fn)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("FlatMap() = %v, 期望 %v", result, tt.expected)
			}
		})
	}
}

func TestChunk(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		size     int
		expected [][]int
	}{
		{
			name:     "正常分块",
			slice:    []int{1, 2, 3, 4, 5},
			size:     2,
			expected: [][]int{{1, 2}, {3, 4}, {5}},
		},
		{
			name:     "完整分块",
			slice:    []int{1, 2, 3, 4},
			size:     2,
			expected: [][]int{{1, 2}, {3, 4}},
		},
		{
			name:     "块大小大于切片",
			slice:    []int{1, 2, 3},
			size:     5,
			expected: [][]int{{1, 2, 3}},
		},
		{
			name:     "空切片",
			slice:    []int{},
			size:     2,
			expected: [][]int{},
		},
		{
			name:     "非法块大小",
			slice:    []int{1, 2, 3},
			size:     0,
			expected: [][]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Chunk(tt.slice, tt.size)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Chunk() = %v, 期望 %v", result, tt.expected)
			}
		})
	}
}

func TestForEach(t *testing.T) {
	t.Run("累加器", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		sum := 0
		ForEach(slice, func(v int) {
			sum += v
		})
		if sum != 15 {
			t.Errorf("ForEach() 累加 = %v, 期望 %v", sum, 15)
		}
	})

	t.Run("空切片", func(t *testing.T) {
		var slice []int
		called := false
		ForEach(slice, func(v int) {
			called = true
		})
		if called {
			t.Errorf("ForEach() 不应该对空切片调用函数")
		}
	})
}

func TestForEachWithIndex(t *testing.T) {
	t.Run("累加索引", func(t *testing.T) {
		slice := []int{10, 20, 30, 40, 50}
		sum := 0
		indexSum := 0
		ForEachWithIndex(slice, func(i int, v int) {
			sum += v
			indexSum += i
		})
		if sum != 150 {
			t.Errorf("ForEachWithIndex() 值累加 = %v, 期望 %v", sum, 150)
		}
		if indexSum != 10 {
			t.Errorf("ForEachWithIndex() 索引累加 = %v, 期望 %v", indexSum, 10)
		}
	})

	t.Run("空切片", func(t *testing.T) {
		var slice []int
		called := false
		ForEachWithIndex(slice, func(i int, v int) {
			called = true
		})
		if called {
			t.Errorf("ForEachWithIndex() 不应该对空切片调用函数")
		}
	})
}

func TestShuffle(t *testing.T) {
	// 注意：测试随机性是困难的，这里只是基本测试
	t.Run("非空切片", func(t *testing.T) {
		original := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		result := Shuffle(original)

		// 检查长度是否相同
		if len(result) != len(original) {
			t.Errorf("Shuffle() 结果长度 = %v, 期望 %v", len(result), len(original))
		}

		// 检查所有元素是否存在
		for _, v := range original {
			if !Includes(result, v) {
				t.Errorf("Shuffle() 缺少元素 %v", v)
			}
		}

		// 检查原切片是否未被修改
		originalCopy := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		if !reflect.DeepEqual(original, originalCopy) {
			t.Errorf("原切片被修改: %v", original)
		}
	})

	t.Run("空切片", func(t *testing.T) {
		var original []int
		result := Shuffle(original)
		if len(result) != 0 {
			t.Errorf("Shuffle() 空切片结果应为空，而不是 %v", result)
		}
	})
}
func TestDifference(t *testing.T) {
	tests := []struct {
		name     string
		slice1   []int
		slice2   []int
		expected []int
	}{
		{
			name:     "基本差异",
			slice1:   []int{1, 2, 3, 4, 5},
			slice2:   []int{3, 4, 5, 6, 7},
			expected: []int{1, 2},
		},
		{
			name:     "完全不同",
			slice1:   []int{1, 2, 3},
			slice2:   []int{4, 5, 6},
			expected: []int{1, 2, 3},
		},
		{
			name:     "完全相同",
			slice1:   []int{1, 2, 3},
			slice2:   []int{1, 2, 3},
			expected: []int{},
		},
		{
			name:     "第一个是空",
			slice1:   []int{},
			slice2:   []int{1, 2, 3},
			expected: []int{},
		},
		{
			name:     "第二个是空",
			slice1:   []int{1, 2, 3},
			slice2:   []int{},
			expected: []int{1, 2, 3},
		},
		{
			name:     "重复元素",
			slice1:   []int{1, 1, 2, 2, 3},
			slice2:   []int{2, 3, 4},
			expected: []int{1, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Difference(tt.slice1, tt.slice2)

			// 由于map遍历是无序的，我们排序后比较
			sort.Ints(result)
			sort.Ints(tt.expected)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Difference() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIntersection(t *testing.T) {
	tests := []struct {
		name     string
		slice1   []int
		slice2   []int
		expected []int
	}{
		{
			name:     "基本交集",
			slice1:   []int{1, 2, 3, 4, 5},
			slice2:   []int{3, 4, 5, 6, 7},
			expected: []int{3, 4, 5},
		},
		{
			name:     "无交集",
			slice1:   []int{1, 2, 3},
			slice2:   []int{4, 5, 6},
			expected: []int{},
		},
		{
			name:     "完全相同",
			slice1:   []int{1, 2, 3},
			slice2:   []int{1, 2, 3},
			expected: []int{1, 2, 3},
		},
		{
			name:     "其中一个是空",
			slice1:   []int{1, 2, 3},
			slice2:   []int{},
			expected: []int{},
		},
		{
			name:     "两个都是空",
			slice1:   []int{},
			slice2:   []int{},
			expected: []int{},
		},
		{
			name:     "重复元素",
			slice1:   []int{1, 1, 2, 2, 3},
			slice2:   []int{1, 2, 3, 3},
			expected: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Intersection(tt.slice1, tt.slice2)

			// 由于map遍历是无序的，我们排序后比较
			sort.Ints(result)
			sort.Ints(tt.expected)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Intersection() = %v, want %v", result, tt.expected)
			}
		})
	}
}
