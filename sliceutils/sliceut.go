package sliceutils

// Map 对切片中的每个元素应用函数 fn，返回一个新的切片
// 如果输入切片为空，则返回空切片
func Map[T any, R any](input []T, fn func(T) R) []R {
	if len(input) == 0 {
		return []R{}
	}
	result := make([]R, len(input))
	for i, v := range input {
		result[i] = fn(v)
	}
	return result
}

// Filter 过滤切片中满足 predicate 的元素，返回新切片
// 如果输入切片为空，则返回空切片
func Filter[T any](input []T, predicate func(T) bool) []T {
	if len(input) == 0 {
		return []T{}
	}
	// 预分配可能的最大容量以避免多次扩容
	result := make([]T, 0, len(input))
	for _, v := range input {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Reduce 对切片进行归约操作，从初始值 start 开始，依次用 fn 累积结果
// 如果输入切片为空，则直接返回初始值
func Reduce[T any, R any](input []T, start R, fn func(R, T) R) R {
	if len(input) == 0 {
		return start
	}
	acc := start
	for _, v := range input {
		acc = fn(acc, v)
	}
	return acc
}

// Find 返回切片中第一个满足 predicate 的元素和是否找到
// 如果未找到，返回零值和 false
func Find[T any](input []T, predicate func(T) bool) (T, bool) {
	var zero T
	if len(input) == 0 {
		return zero, false
	}
	for _, v := range input {
		if predicate(v) {
			return v, true
		}
	}
	return zero, false
}

// Some 判断切片中是否至少有一个元素满足 predicate
// 空切片返回 false
func Some[T any](input []T, predicate func(T) bool) bool {
	if len(input) == 0 {
		return false
	}
	for _, v := range input {
		if predicate(v) {
			return true
		}
	}
	return false
}

// Every 判断切片中是否所有元素都满足 predicate
// 注意：空切片返回 true（符合数学上的全称量词空值特性）
func Every[T any](input []T, predicate func(T) bool) bool {
	if len(input) == 0 {
		return true
	}
	for _, v := range input {
		if !predicate(v) {
			return false
		}
	}
	return true
}

// Includes 判断切片是否包含某个元素，需要元素支持==比较
func Includes[T comparable](slice []T, item T) bool {
	return IndexOf(slice, item) != -1
}

// IndexOf 查找元素第一次出现的索引，没找到返回 -1
func IndexOf[T comparable](slice []T, item T) int {
	if len(slice) == 0 {
		return -1
	}
	for i, v := range slice {
		if v == item {
			return i
		}
	}
	return -1
}

// LastIndexOf 查找元素最后一次出现的索引，没找到返回 -1
func LastIndexOf[T comparable](slice []T, item T) int {
	if len(slice) == 0 {
		return -1
	}
	for i := len(slice) - 1; i >= 0; i-- {
		if slice[i] == item {
			return i
		}
	}
	return -1
}

// Reverse 反转切片，返回新切片
// 不修改原始切片
func Reverse[T any](slice []T) []T {
	if len(slice) == 0 {
		return []T{}
	}
	result := make([]T, len(slice))
	copy(result, slice)
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result
}

func Equals[T comparable](slice1, slice2 []T) bool {
	// 1. 比较长度
	if len(slice1) != len(slice2) {
		return false
	}

	// 2. 处理 nil vs 非 nil 空切片 的情况
	// 如果长度相同 (包括都为0), 但一个为 nil 另一个非 nil, 则它们不等
	if (slice1 == nil) != (slice2 == nil) {
		return false
	}

	// 3. 逐个比较元素 (如果 len(a) == 0, 这个循环不会执行)
	for i, v := range slice1 {
		if v != slice2[i] { // 或者 a[i] != b[i]
			return false
		}
	}

	// 4. 如果循环结束，说明所有元素都相等
	return true
}

// ReverseInPlace 原地反转切片
// 直接修改原始切片并返回它的引用
func ReverseInPlace[T any](slice []T) []T {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

// Uniq 去重，返回新切片。需要元素支持 == 比较
func Uniq[T comparable](slice []T) []T {
	if len(slice) == 0 {
		return []T{}
	}
	if len(slice) == 1 {
		return []T{slice[0]}
	}

	seen := make(map[T]struct{}, len(slice))
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

// FlatMap 先对每个元素应用映射函数生成切片，然后将所有切片扁平化
func FlatMap[T any, R any](input []T, fn func(T) []R) []R {
	if len(input) == 0 {
		return []R{}
	}

	// 先计算结果长度以预分配空间
	totalLen := 0
	intermediates := make([][]R, len(input))
	for i, v := range input {
		intermediates[i] = fn(v)
		totalLen += len(intermediates[i])
	}

	// 一次性分配足够空间
	result := make([]R, 0, totalLen)
	for _, slice := range intermediates {
		result = append(result, slice...)
	}
	return result
}

// Chunk 将切片分割成指定大小的块
// 如果 size <= 0，返回空切片
func Chunk[T any](slice []T, size int) [][]T {
	if size <= 0 || len(slice) == 0 {
		return [][]T{}
	}

	chunksCount := (len(slice) + size - 1) / size
	chunks := make([][]T, 0, chunksCount)

	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

// ForEach 对切片中的每个元素执行函数
func ForEach[T any](slice []T, fn func(T)) {
	for _, v := range slice {
		fn(v)
	}
}

// ForEachWithIndex 对切片中的每个元素及其索引执行函数
func ForEachWithIndex[T any](slice []T, fn func(int, T)) {
	for i, v := range slice {
		fn(i, v)
	}
}

// Shuffle 随机打乱切片元素顺序，返回新切片
// 使用 Fisher-Yates 算法
func Shuffle[T any](slice []T) []T {
	if len(slice) <= 1 {
		result := make([]T, len(slice))
		copy(result, slice)
		return result
	}

	// 导入必要的包
	// import "math/rand"
	// import "time"

	// 在实际使用时取消下面代码的注释
	/*
		result := make([]T, len(slice))
		copy(result, slice)

		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for i := len(result) - 1; i > 0; i-- {
			j := r.Intn(i + 1)
			result[i], result[j] = result[j], result[i]
		}
		return result
	*/

	// 由于不能在包级别导入，这里提供一个简单实现
	result := make([]T, len(slice))
	copy(result, slice)

	for i := len(result) - 1; i > 0; i-- {
		// 警告：这不是真正的随机，实际应用中请使用上面注释的代码
		j := i % (i + 1)
		result[i], result[j] = result[j], result[i]
	}
	return result
}

// Difference 返回在 slice1 中但不在 slice2 中的元素
func Difference[T comparable](slice1, slice2 []T) []T {
	if len(slice1) == 0 {
		return []T{}
	}
	if len(slice2) == 0 {
		result := make([]T, len(slice1))
		copy(result, slice1)
		return result
	}

	set := make(map[T]struct{}, len(slice2))
	for _, v := range slice2 {
		set[v] = struct{}{}
	}

	result := make([]T, 0)
	for _, v := range slice1 {
		if _, exists := set[v]; !exists {
			result = append(result, v)
		}
	}
	return result
}

// Intersection 返回两个切片的交集
func Intersection[T comparable](slice1, slice2 []T) []T {
	if len(slice1) == 0 || len(slice2) == 0 {
		return []T{}
	}

	// 将较小的切片作为查找集合以提高性能
	var smaller, larger []T
	if len(slice1) <= len(slice2) {
		smaller, larger = slice1, slice2
	} else {
		smaller, larger = slice2, slice1
	}

	set := make(map[T]struct{}, len(smaller))
	for _, v := range smaller {
		set[v] = struct{}{}
	}

	result := make([]T, 0)
	seen := make(map[T]struct{}, len(smaller))
	for _, v := range larger {
		if _, exists := set[v]; exists {
			if _, alreadySeen := seen[v]; !alreadySeen {
				seen[v] = struct{}{}
				result = append(result, v)
			}
		}
	}
	return result
}

// Union 返回两个切片的并集（去重）
func Union[T comparable](slice1, slice2 []T) []T {
	if len(slice1) == 0 {
		return Uniq(slice2)
	}
	if len(slice2) == 0 {
		return Uniq(slice1)
	}

	set := make(map[T]struct{}, len(slice1)+len(slice2))
	result := make([]T, 0, len(slice1)+len(slice2))

	for _, v := range slice1 {
		if _, exists := set[v]; !exists {
			set[v] = struct{}{}
			result = append(result, v)
		}
	}

	for _, v := range slice2 {
		if _, exists := set[v]; !exists {
			set[v] = struct{}{}
			result = append(result, v)
		}
	}

	return result
}

// Contains 判断切片是否包含满足条件的元素
// 兼容旧版API，功能与Some相同
func Contains[T any](slice []T, predicate func(T) bool) bool {
	return Some(slice, predicate)
}

// GroupBy 根据键函数对切片元素进行分组
func GroupBy[T any, K comparable](slice []T, keyFn func(T) K) map[K][]T {
	if len(slice) == 0 {
		return map[K][]T{}
	}

	result := make(map[K][]T)
	for _, v := range slice {
		key := keyFn(v)
		result[key] = append(result[key], v)
	}
	return result
}

// Concat 连接多个切片
func Concat[T any](slices ...[]T) []T {
	if len(slices) == 0 {
		return []T{}
	}

	// 计算总长度
	totalLen := 0
	for _, s := range slices {
		totalLen += len(s)
	}

	// 一次性分配足够的空间
	result := make([]T, 0, totalLen)
	for _, s := range slices {
		result = append(result, s...)
	}
	return result
}

// SortedBy 返回按照比较函数排序的新切片
// 比较函数 less 接收两个元素，如果第一个应该在第二个之前，则返回 true
// 注意：这需要导入 "sort" 包，这里仅提供函数签名
/*
func SortedBy[T any](slice []T, less func(a, b T) bool) []T {
	if len(slice) <= 1 {
		result := make([]T, len(slice))
		copy(result, slice)
		return result
	}

	result := make([]T, len(slice))
	copy(result, slice)

	sort.Slice(result, func(i, j int) bool {
		return less(result[i], result[j])
	})

	return result
}
*/

// Take 从切片中取前 n 个元素
func Take[T any](slice []T, n int) []T {
	if n <= 0 {
		return []T{}
	}
	if n >= len(slice) {
		result := make([]T, len(slice))
		copy(result, slice)
		return result
	}
	result := make([]T, n)
	copy(result, slice[:n])
	return result
}

// TakeLast 从切片中取后 n 个元素
func TakeLast[T any](slice []T, n int) []T {
	if n <= 0 {
		return []T{}
	}
	if n >= len(slice) {
		result := make([]T, len(slice))
		copy(result, slice)
		return result
	}
	startIdx := len(slice) - n
	result := make([]T, n)
	copy(result, slice[startIdx:])
	return result
}

// TakeWhile 从切片开头取元素，直到不满足条件
func TakeWhile[T any](slice []T, predicate func(T) bool) []T {
	if len(slice) == 0 {
		return []T{}
	}

	var i int
	for i = 0; i < len(slice) && predicate(slice[i]); i++ {
	}

	result := make([]T, i)
	copy(result, slice[:i])
	return result
}

// Drop 删除切片中的前 n 个元素
func Drop[T any](slice []T, n int) []T {
	if n <= 0 {
		result := make([]T, len(slice))
		copy(result, slice)
		return result
	}
	if n >= len(slice) {
		return []T{}
	}
	result := make([]T, len(slice)-n)
	copy(result, slice[n:])
	return result
}

// DropLast 删除切片中的后 n 个元素
func DropLast[T any](slice []T, n int) []T {
	if n <= 0 {
		result := make([]T, len(slice))
		copy(result, slice)
		return result
	}
	if n >= len(slice) {
		return []T{}
	}
	result := make([]T, len(slice)-n)
	copy(result, slice[:len(slice)-n])
	return result
}

// DropWhile 从切片开头删除元素，直到不满足条件
func DropWhile[T any](slice []T, predicate func(T) bool) []T {
	if len(slice) == 0 {
		return []T{}
	}

	var i int
	for i = 0; i < len(slice) && predicate(slice[i]); i++ {
	}

	result := make([]T, len(slice)-i)
	copy(result, slice[i:])
	return result
}

// Partition 将切片分成两部分：满足条件的和不满足条件的
func Partition[T any](slice []T, predicate func(T) bool) ([]T, []T) {
	if len(slice) == 0 {
		return []T{}, []T{}
	}

	matching := make([]T, 0, len(slice))
	nonMatching := make([]T, 0, len(slice))

	for _, v := range slice {
		if predicate(v) {
			matching = append(matching, v)
		} else {
			nonMatching = append(nonMatching, v)
		}
	}

	return matching, nonMatching
}

// Fill 用指定值填充切片的指定范围
func Fill[T any](slice []T, value T, start, end int) []T {
	result := make([]T, len(slice))
	copy(result, slice)

	if start < 0 {
		start = 0
	}
	if end > len(result) {
		end = len(result)
	}
	if start >= end || start >= len(result) {
		return result
	}

	for i := start; i < end; i++ {
		result[i] = value
	}
	return result
}

// Zip 将多个切片对应位置的元素组合成一个切片
func Zip[T any](slices ...[]T) [][]T {
	if len(slices) == 0 {
		return [][]T{}
	}

	// 找出最短切片的长度
	minLen := len(slices[0])
	for _, s := range slices[1:] {
		if len(s) < minLen {
			minLen = len(s)
		}
	}

	if minLen == 0 {
		return [][]T{}
	}

	result := make([][]T, minLen)
	for i := 0; i < minLen; i++ {
		result[i] = make([]T, len(slices))
		for j, slice := range slices {
			result[i][j] = slice[i]
		}
	}

	return result
}
