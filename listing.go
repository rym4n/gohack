package gohack

// Replacer 接口
type Replacer interface {
	Len() int
	Replace([]int) Replacer
	ToString() string
}

// IntReplacer 适用元素为int类型
type IntReplacer []int

// Len 计算元素数量
func (ir IntReplacer) Len() int {
	return len(ir)
}

// Replace
func (ir IntReplacer) Replace(indices []int) Replacer {
	result := make(IntReplacer, len(indices), len(indices))
	for i, idx := range indices {
		result[i] = ir[idx]
	}
	return result
}

// ToString 输出字符串
func (ir IntReplacer) ToString() string {
	result := ""
	for _, v := range ir {
		result += StrVal(v)
	}
	return result
}

type RuneReplacer []rune

func (rr RuneReplacer) Len() int {
	return len(rr)
}

func (rr RuneReplacer) Replace(indices []int) Replacer {
	result := make(RuneReplacer, len(indices), len(indices))
	for i, idx := range indices {
		result[i] = rr[idx]
	}
	return result
}

func (rr RuneReplacer) ToString() string {
	result := ""
	for _, v := range rr {
		result += StrVal(v)
	}
	return result
}

type StringReplacer []string

func (sr StringReplacer) Len() int {
	return len(sr)
}

func (sr StringReplacer) Replace(indices []int) Replacer {
	result := make(StringReplacer, len(indices), len(indices))
	for i, idx := range indices {
		result[i] = sr[idx]
	}
	return result
}

func (sr StringReplacer) ToString() string {
	result := ""
	for _, v := range sr {
		result += v
	}
	return result
}

type Float64Replacer []string

func (fr Float64Replacer) Len() int {
	return len(fr)
}
func (fr Float64Replacer) Replace(indices []int) Replacer {
	result := make(Float64Replacer, len(indices), len(indices))
	for i, idx := range indices {
		result[i] = fr[idx]
	}
	return result
}

func (fr Float64Replacer) ToString() string {
	result := ""
	for _, v := range fr {
		result += StrVal(v)
	}
	return result
}

/*
Permutations 生成排列数
@Replacer: 元素列表
@selectNum: 排列长度
@repeatable: 是否允许元素重复
@bufSize: 消息通道缓存大小
*/
func Permutations(list Replacer, selectNum int, repeatable bool, bufSize int) (c chan Replacer) {
	c = make(chan Replacer, bufSize)
	go func() {
		defer close(c)
		var permGenerator func([]int, int, int) chan []int
		if repeatable {
			permGenerator = repeatedPermutations
		} else {
			permGenerator = permutations
		}
		indices := make([]int, list.Len(), list.Len())
		for i := 0; i < list.Len(); i++ {
			indices[i] = i
		}
		for perm := range permGenerator(indices, selectNum, bufSize) {
			c <- list.Replace(perm)
		}
	}()
	return
}

func pop(l []int, i int) (v int, sl []int) {
	v = l[i]
	length := len(l)
	sl = make([]int, length-1, length-1)
	copy(sl, l[:i])
	copy(sl[i:], l[i+1:])
	return
}

func permutations(list []int, selectNum, bufSize int) (c chan []int) {
	c = make(chan []int, bufSize)
	go func() {
		defer close(c)
		switch selectNum {
		case 1:
			for _, v := range list {
				c <- []int{v}
			}
			return
		case 0:
			return
		case len(list):
			for i := 0; i < len(list); i++ {
				top, subList := pop(list, i)
				for perm := range permutations(subList, selectNum-1, bufSize) {
					c <- append([]int{top}, perm...)
				}
			}
		default:
			for comb := range combinations(list, selectNum, bufSize) {
				for perm := range permutations(comb, selectNum, bufSize) {
					c <- perm
				}
			}
		}
	}()
	return
}

func repeatedPermutations(list []int, selectNum, bufSize int) (c chan []int) {
	c = make(chan []int, bufSize)
	go func() {
		defer close(c)
		switch selectNum {
		case 1:
			for _, v := range list {
				c <- []int{v}
			}
		default:
			for i := 0; i < len(list); i++ {
				for perm := range repeatedPermutations(list, selectNum-1, bufSize) {
					c <- append([]int{list[i]}, perm...)
				}
			}
		}
	}()
	return
}

/*
Combinations 组合数生成器
@list: 元素列表
@selectNum: 生成的组合数长度
@repeatable: 是否允许重复元素
@bufSize: 消息通道缓存大小
*/
func Combinations(list Replacer, selectNum int, repeatable bool, bufSize int) (c chan Replacer) {
	c = make(chan Replacer, bufSize)
	index := make([]int, list.Len(), list.Len())
	for i := 0; i < list.Len(); i++ {
		index[i] = i
	}

	var combGenerator func([]int, int, int) chan []int
	if repeatable {
		combGenerator = repeatedCombinations
	} else {
		combGenerator = combinations
	}

	go func() {
		defer close(c)
		for comb := range combGenerator(index, selectNum, bufSize) {
			c <- list.Replace(comb)
		}
	}()

	return
}

func combinations(list []int, selectNum, bufSize int) (c chan []int) {
	c = make(chan []int, bufSize)
	go func() {
		defer close(c)
		switch {
		case selectNum == 0:
			c <- []int{}
		case selectNum == len(list):
			c <- list
		case len(list) < selectNum:
			return
		default:
			for i := 0; i < len(list); i++ {
				for subComb := range combinations(list[i+1:], selectNum-1, bufSize) {
					c <- append([]int{list[i]}, subComb...)
				}
			}
		}
	}()
	return
}

func repeatedCombinations(list []int, selectNum, bufSize int) (c chan []int) {
	c = make(chan []int, bufSize)
	go func() {
		defer close(c)
		if selectNum == 1 {
			for v := range list {
				c <- []int{v}
			}
			return
		}
		for i := 0; i < len(list); i++ {
			for subComb := range repeatedCombinations(list[i:], selectNum-1, bufSize) {
				c <- append([]int{list[i]}, subComb...)
			}
		}
	}()
	return
}
