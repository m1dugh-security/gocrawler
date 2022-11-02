package main

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
    middle := start / 2 + (end - start) / 2
    if end <= start {
        return start
    }

    c := children[middle].chr
    if c == x {
        return middle
    } else if x < c {
        return _binSearch(children, start, middle, x)
    } else {
        return _binSearch(children, middle + 1, end, x)
    }
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
        if pos > children {
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
