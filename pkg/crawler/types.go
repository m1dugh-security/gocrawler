package crawler

import "fmt"

type PrefixTree struct {
    chr byte
    children []*PrefixTree
    occurences int
}


func (parent *PrefixTree) _insertChildrenAt(index int,
            chr byte,
            occurences int) *PrefixTree {
    var child *PrefixTree = &PrefixTree{chr, nil, occurences}
    children := parent.children
    if index == len(children) {
        parent.children = append(children, child)
    } else if index < len(children) {
        parent.children = append(children[:index + 1], children[index:]...)
        parent.children[index] = child
    }

    return child
}


func _binSearch(children []*PrefixTree, start int, end int, x byte) int {

    for start < end {
        middle := start + (end - start) / 2

        c := children[middle].chr
        if c == x {
            return middle
        } else if x < c {
            end = middle
        } else {
            start = middle + 1
        }
    }

    return start
}

func (t *PrefixTree) _searchWord(str string, strlen int, index int) int {
    if index == strlen {
        return t.occurences
    }
    c := str[index]
    children := len(t.children)
    if children == 0 {
        return 0
    } else {
        pos := _binSearch(t.children, 0, children, c)
        if pos >= children {
            return 0
        } else if t.children[pos].chr == c {
            return t.children[pos]._searchWord(str, strlen, index + 1)
        } else {
            return 0
        }
    }
}

func (t *PrefixTree) SearchWord(str string) int {
    return t._searchWord(str, len(str), 0)
}

func (t *PrefixTree) _listWords(prefix string) []string {

    var res []string = nil
    prefix = fmt.Sprintf("%s%c", prefix, t.chr)
    if t.occurences > 0 {
        res = append(res, prefix)
    }
    for _, v := range t.children {
        res = append(res, v._listWords(prefix)...)
    }

    return res
}

func (t *PrefixTree) ListWords() []string {
    return t._listWords("")
}

func CreatePrefixTree() *PrefixTree {
    res := &PrefixTree{0, nil, 0}
    return res
}


/**
 *  \return the (Found Prefix tree or parent if not found,
 the last found index + 1,
 the expected pos of the children)
 */
func (t *PrefixTree) _searchNode(str string, strlen int, index int) (*PrefixTree, int, int) {
    children := len(t.children)
    if index == strlen {
        return t, index, -1
    } else if children == 0 {
        return t, index, 0
    } else {
        c := str[index]
        pos := _binSearch(t.children, 0, children, c)
        if pos >= children || t.children[pos].chr != c {
            return t, index, pos
        } else {
            return t.children[pos]._searchNode(str, strlen, index + 1)
        }
    }
}

func (t *PrefixTree) AddWord(word string) bool {
    wl := len(word)
    node, index, pos := t._searchNode(word, wl, 0)
    if pos == -1 {
        res := node.occurences == 0
        node.occurences++
        return res
    } else {
        node = node._insertChildrenAt(pos, word[index], 0)
        index++
        for index < wl {
            pos = _binSearch(node.children, 0, len(node.children), word[index])
            if pos >= len(node.children) || node.children[pos].chr != word[index] {
                node = node._insertChildrenAt(pos, word[index], 0)
            } else {
                node = node.children[pos]
            }
            index++
        }

        res := node.occurences == 0
        node.occurences++
        return res
    }
}


func _compareStrings(a string, b string) int {
    la := len(a)
    lb := len(b)
    var min int

    if la > lb {
        min = lb
    } else {
        min = la
    }
    var res int = 0

    for i := 0; i < min && res == 0; i++ {
        res = int(a[i]) - int(b[i])
    }

    if res == 0 {
        if la < lb {
            return -1
        } else if la == lb {
            return 0
        } else {
            return 1
        }
    }

    return res
}


type StringSet []string;

func NewStringSet(values []string) *StringSet {
    
    var res *StringSet = &StringSet{}
    *res = make([]string, len(values))

    for _, v := range values {
        res.AddWord(v)
    }

    return res
}

func (set *StringSet) _binsearch(value string) (int, bool) {

    start := 0
    end := len(*set)

    for start < end {
        middle := start + (end - start) / 2
        s := (*set)[middle]
        res := _compareStrings(value, s)
        if res == 0 {
            return middle, true
        } else if res < 0 {
            end = middle
        } else {
            start = middle + 1
        }
    }

    return start, false
}

func (set *StringSet) _insertAt(value string, pos int) {
    *set = append(*set, value)
    for i := len(*set) - 1; i > pos; i-- {
        (*set)[i] = (*set)[i - 1]
    }

    (*set)[pos] = value
}

func (set *StringSet) AddWord(value string) bool {
    pos, found := set._binsearch(value)
    if found {
        return false
    }

    set._insertAt(value, pos)
    return true
}

func (set *StringSet) ContainsWord(value string) bool {
    _, found := set._binsearch(value)
    return found
}

func (set *StringSet) ToArray() []string {
    dest := make([]string, len(*set))
    copy(dest, *set)
    return dest
}

