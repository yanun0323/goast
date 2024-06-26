package goast

type trie[T comparable] interface {
	AddText(s string, value T) *trie_[T]
	FindText(s string) (T, bool)
	FindTextReversely(s string) (T, bool)
	FindByte(b []byte) (T, bool)
	FindByteReversely(b []byte) (T, bool)
	Insert(key byte, value T) *trie_[T]
	Next(key byte) (*trie_[T], bool)
	Value() T
}

type trie_[T comparable] struct {
	value T
	next  map[byte]*trie_[T]
}

func (t *trie_[T]) AddText(s string, value T) *trie_[T] {
	if t == nil {
		return nil
	}

	var zero T
	n := t
	for i := range s {
		n = n.Insert(s[i], zero)
	}
	n.value = value
	return n
}

func (t *trie_[T]) Insert(key byte, value T) *trie_[T] {
	if t == nil {
		return nil
	}

	if t.next == nil {
		t.next = map[byte]*trie_[T]{}
	}

	if _, ok := t.next[key]; !ok {
		t.next[key] = &trie_[T]{}
	}

	t.next[key].value = value

	return t.next[key]
}

func (t *trie_[T]) FindText(s string) (T, bool) {
	if t == nil {
		return t.Value(), false
	}

	n := t
	ok := false
	for i := range s {
		n, ok = n.Next(s[i])
		if !ok {
			return n.Value(), false
		}
	}
	return n.Value(), true
}

func (t *trie_[T]) FindTextReversely(s string) (T, bool) {
	if t == nil {
		return t.Value(), false
	}

	n := t
	ok := false
	for i := len(s) - 1; i >= 0; i-- {
		n, ok = n.Next(s[i])
		if !ok {
			return n.Value(), false
		}
	}
	return n.Value(), true
}

func (t *trie_[T]) FindByte(b []byte) (T, bool) {
	if t == nil {
		return t.Value(), false
	}

	n := t
	ok := false
	for i := range b {
		n, ok = n.Next(b[i])
		if !ok {
			return n.Value(), false
		}
	}

	return n.Value(), true
}

func (t *trie_[T]) FindByteReversely(b []byte) (T, bool) {
	if t == nil {
		return t.Value(), false
	}

	n := t
	ok := false
	for i := len(b) - 1; i >= 0; i-- {
		n, ok = n.Next(b[i])
		if !ok {
			return n.Value(), false
		}
	}

	return n.Value(), true
}

func (t *trie_[T]) Value() T {
	if t == nil {
		var zero T
		return zero
	}
	return t.value
}

func (t *trie_[T]) Next(key byte) (*trie_[T], bool) {
	if t == nil || t.next == nil {
		return nil, false
	}

	n, ok := t.next[key]
	return n, ok
}

func newTrie[T comparable](set map[string]T) trie[T] {
	root := &trie_[T]{}
	for k, v := range set {
		_ = root.AddText(k, v)
	}

	return root
}
