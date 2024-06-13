package common

//CREDIT TO https://github.com/kevin-wayne/algs4/blob/master/src/main/java/edu/princeton/cs/algs4/TST.java

type Node struct {
	char  rune
	left  *Node
	mid   *Node
	right *Node
	val   string
}

type TST struct {
	root *Node
	n    int
}

// NewNode creates and initializes a new Node with the given character
func NewNode(char rune) *Node {
	return &Node{
		char:  char,
		left:  nil,
		mid:   nil,
		right: nil,
		val:   "",
	}
}

// Size returns the number of key-value pairs in this symbol table
func (t *TST) Size() int {
	return t.n
}

// Contains checks if this symbol table contains the given key
func (t *TST) Contains(key string) bool {
	if key == "" {
		panic("argument to contains() is empty")
	}
	return t.Get(key) != nil
}

// Get returns the value associated with the given key
func (t *TST) Get(key string) interface{} {
	if key == "" {
		panic("calls get() with empty argument")
	}
	if len(key) == 0 {
		panic("key must have length >= 1")
	}
	x := t.get(t.root, key, 0)
	if x == nil {
		return nil
	}
	return x.val
}

func (t *TST) get(x *Node, key string, d int) *Node {
	if x == nil {
		return nil
	}
	if len(key) == 0 {
		panic("key must have length >= 1")
	}
	char := rune(key[d])
	if char < x.char {
		return t.get(x.left, key, d)
	} else if char > x.char {
		return t.get(x.right, key, d)
	} else if d < len(key)-1 {
		return t.get(x.mid, key, d+1)
	} else {
		return x
	}
}

// Put inserts the key-value pair into the symbol table
func (t *TST) Put(key string, val string) {
	if key == "" {
		panic("calls put() with empty key")
	}
	if !t.Contains(key) {
		t.n++
	} else if val == "" {
		t.n-- // delete existing key
	}
	t.root = t.put(t.root, key, val, 0)
}

func (t *TST) put(x *Node, key string, val string, d int) *Node {
	char := rune(key[d])
	if x == nil {
		x = NewNode(char)
	}
	if char < x.char {
		x.left = t.put(x.left, key, val, d)
	} else if char > x.char {
		x.right = t.put(x.right, key, val, d)
	} else if d < len(key)-1 {
		x.mid = t.put(x.mid, key, val, d+1)
	} else {
		x.val = val
	}
	return x
}

// LongestPrefixOf returns the string in the symbol table that is the longest prefix of query
func (t *TST) LongestPrefixOf(query string) string {
	if query == "" {
		panic("calls longestPrefixOf() with empty argument")
	}
	if len(query) == 0 {
		return ""
	}
	length := 0
	x := t.root
	i := 0
	for x != nil && i < len(query) {
		char := rune(query[i])
		if char < x.char {
			x = x.left
		} else if char > x.char {
			x = x.right
		} else {
			i++
			if x.val != "" {
				length = i
			}
			x = x.mid
		}
	}
	return query[:length]
}

// Keys returns all keys in the symbol table as a slice of strings
func (t *TST) Keys() []string {
	queue := []string{}
	t.collect(t.root, "", &queue)
	return queue
}

// KeysWithPrefix returns all of the keys in the set that start with prefix
func (t *TST) KeysWithPrefix(prefix string) []string {
	if prefix == "" {
		panic("calls keysWithPrefix() with empty argument")
	}
	queue := []string{}
	x := t.get(t.root, prefix, 0)
	if x == nil {
		return queue
	}
	if x.val != "" {
		queue = append(queue, prefix)
	}
	t.collect(x.mid, prefix, &queue)
	return queue
}

// collect collects all keys in subtrie rooted at x with given prefix
func (t *TST) collect(x *Node, prefix string, queue *[]string) {
	if x == nil {
		return
	}
	t.collect(x.left, prefix, queue)
	if x.val != "" {
		*queue = append(*queue, prefix+string(x.char))
	}
	t.collect(x.mid, prefix+string(x.char), queue)
	t.collect(x.right, prefix, queue)
}
