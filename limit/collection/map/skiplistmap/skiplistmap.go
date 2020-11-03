package skiplistmap

import (
	"LimitGo/limit/collection"
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
)

const (
	EQ = 1
	LT = 2
	GT = 0
)

type SkipListMap struct {
	head    *Index
	precede func(p1 *collection.Object, p2 *collection.Object) bool
	kt      reflect.Type
	vt      reflect.Type
	size    int

	keySet   *KeySet
	entrySet *EntrySet
	values   *Values
}

type EntryIterator struct {
	sm      *SkipListMap
	next    *Node
	lastRet *Node
}

func New(kt reflect.Type, vt reflect.Type, precede func(p1 *collection.Object, p2 *collection.Object) bool) *SkipListMap {
	head := &Index{1, &Node{nil, nil, nil}, nil, nil}
	return &SkipListMap{head, precede, kt, vt, 0, nil, nil, nil}
}

// Size returns the number of elements in this collection.
func (sm *SkipListMap) Size() int {
	return sm.size
}

// Empty returns true if this collection contains no element.
func (sm *SkipListMap) Empty() bool {
	return sm.size == 0
}

// String returns a string representation of this collection.
func (sm *SkipListMap) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	it := sm.GetEntryIterator()
	for it.HashNext() {
		entry := (*it.Next()).(collection.Entry)
		if buf.Len() > len("{") {
			buf.WriteByte(',')
		}
		key := entry.GetKey()
		value := entry.GetValue()
		var s string
		k, err1 := json.Marshal(*key)
		v, err2 := json.Marshal(*value)
		if err1 == nil && err2 == nil {
			s = string(k) + "=" + string(v)
		} else {
			s = "nil"
		}
		_, _ = fmt.Fprint(&buf, s)
	}
	buf.WriteByte('}')
	return buf.String()
}

// Removes all of the elements from this collection.
func (sm *SkipListMap) Clear() bool {
	sm.head = &Index{1, &Node{nil, nil, nil}, nil, nil}
	sm.size = 0
	sm.keySet = &KeySet{sm}
	sm.entrySet = &EntrySet{sm}
	sm.values = &Values{sm}
	return true
}

// GetEntryIterator returns iterator of entry.
func (sm *SkipListMap) GetEntryIterator() collection.EntryItr {
	it := EntryIterator{sm, sm.findFirst(), nil}
	return &it
}

// ContainsKey returns true if this map contains a mapping for the specified key.
func (sm *SkipListMap) ContainsKey(key *collection.Object) bool {
	if sm.checkNil(key) || !sm.checkKeyType(key) {
		return false
	}
	return sm.doGet(key) != nil
}

// ContainsValue returns true if this map maps one or more keys to the
// specified value.
func (sm *SkipListMap) ContainsValue(value *collection.Object) bool {
	if sm.checkNil(value) || !sm.checkValueType(value) {
		return false
	}
	for node := sm.findFirst(); node != nil; node = node.next {
		p := value
		if reflect.DeepEqual(*value, *p) {
			return true
		}
	}
	return false
}

// Get returns the value to which the specified key is mapped, or null
// if this map contains no mapping for the key.
func (sm *SkipListMap) Get(key *collection.Object) *collection.Object {
	if sm.checkNil(key) || !sm.checkKeyType(key) {
		return nil
	}
	return sm.doGet(key)
}

// Put associates the specified value with the specified key, returns old value
// if the specified key has been in this map.
func (sm *SkipListMap) Put(key *collection.Object, value *collection.Object) (bool, *collection.Object) {
	if sm.checkNil(key) || !sm.checkKeyType(key) {
		return false, nil
	}
	if sm.checkNil(value) || !sm.checkValueType(value) {
		return false, nil
	}
	return true, sm.doPut(key, value)
}

// Remove removes the mapping for a key from this map if it is present.
func (sm *SkipListMap) Remove(key *collection.Object) *collection.Object {
	if sm.checkNil(key) || !sm.checkKeyType(key) {
		return nil
	}
	node := sm.doRemove(key)
	if node == nil {
		return nil
	} else {
		return node.value
	}
}

// PutAll copies all of the mappings from the specified map to this map.
func (sm *SkipListMap) PutAll(m *collection.Map) {
	if m == nil || (*m) == nil || (*m).Size() == 0 {
		return
	}
	it := (*m).GetEntryIterator()
	for it.HashNext() {
		entry := (*it.Next()).(collection.Entry)
		sm.Put(entry.GetKey(), entry.GetValue())
	}
}

// KeySet returns a Set view of the keys contained in this map.
func (sm *SkipListMap) KeySet() *collection.Set {
	return nil
}

// Values returns a List view of the values contained in this map.
func (sm *SkipListMap) Values() *collection.Linear {
	return nil
}

// EntrySet returns a Set view of the mappings contained in this map.
func (sm *SkipListMap) EntrySet() *collection.Set {
	return nil
}

// Equals returns true only if the corresponding pairs of the elements
//in the two maps are equal.
func (sm *SkipListMap) Equals(m *collection.Map) bool {
	if m == nil || (*m) == nil || (*m).Size() == 0 || sm.Size() != (*m).Size() {
		return false
	}
	it := (*m).GetEntryIterator()
	for it.HashNext() {
		entry := (*it.Next()).(collection.Entry)
		if !reflect.DeepEqual(*sm.Get(entry.GetKey()), *entry.GetValue()) {
			return false
		}
	}
	return true
}

// SubMap returns a view of the portion of this map whose keys range
// from "fromKey" to "toKey".  If "fromKey" and "toKey" are equal,
// the returned map is empty.)
func (sm *SkipListMap) SubMap(fromKey *collection.Object, fromInclusive bool, toKey *collection.Object, toInclusive bool) *collection.SortedMap {
	if fromKey != nil && !sm.checkKeyType(fromKey) {
		return nil
	}
	if toKey != nil && !sm.checkKeyType(fromKey) {
		return nil
	}
	p := &SubMap{sm, fromKey, toKey, fromInclusive, toInclusive, false, sm.size, nil, nil, nil}
	p.InitSize()
	var ret collection.SortedMap = p
	return &ret
}

// HeadMap returns a view of the portion of this map whose keys are strictly
// less than toKey.
func (sm *SkipListMap) HeadMap(toKey *collection.Object, inclusive bool) *collection.SortedMap {
	if !sm.checkKeyType(toKey) {
		return nil
	}
	return sm.SubMap(nil, false, toKey, inclusive)
}

// TailMap returns a view of the portion of this map whose keys are greater than
// or equal to fromKey.
func (sm *SkipListMap) TailMap(fromKey *collection.Object, inclusive bool) *collection.SortedMap {
	if !sm.checkKeyType(fromKey) {
		return nil
	}
	return sm.SubMap(fromKey, inclusive, nil, false)
}

// SortedKeySet returns a SortedSet view of the keys contained in this map.
func (sm *SkipListMap) SortedKeySet() *collection.SortedSet {
	return nil
}

// LowerEntry returns a key-value mapping associated with the greatest key
// strictly less than the given key, or nil if there is no such key.
func (sm *SkipListMap) LowerEntry(key *collection.Object) *collection.Entry {
	if sm.checkNil(key) || !sm.checkKeyType(key) {
		return nil
	}
	var entry collection.Entry = sm.findNear(key, LT)
	return &entry
}

// FloorEntry returns a key-value mapping associated with the greatest key
// less than or equal to the given key, or nil if there is no such key.
func (sm *SkipListMap) FloorEntry(key *collection.Object) *collection.Entry {
	if sm.checkNil(key) || !sm.checkKeyType(key) {
		return nil
	}
	var entry collection.Entry = sm.findNear(key, LT|EQ)
	return &entry
}

// CeilingEntry returns a key-value mapping associated with the least key
// greater than or equal to the given key, or nil if there is no such key.
func (sm *SkipListMap) CeilingEntry(key *collection.Object) *collection.Entry {
	if sm.checkNil(key) || !sm.checkKeyType(key) {
		return nil
	}
	var entry collection.Entry = sm.findNear(key, GT|EQ)
	return &entry
}

// HigherEntry returns a key-value mapping associated with the least key
// strictly greater than the given key, or nil if there is no such key.
func (sm *SkipListMap) HigherEntry(key *collection.Object) *collection.Entry {
	if sm.checkNil(key) || !sm.checkKeyType(key) {
		return nil
	}
	var entry collection.Entry = sm.findNear(key, GT)
	return &entry
}

// Entry returns a key-value mapping associated with the least key
// in this map, or nil if the map is empty.
func (sm *SkipListMap) FirstEntry() *collection.Entry {
	var entry collection.Entry = sm.findFirst()
	return &entry
}

// LastEntry returns a key-value mapping associated with the greatest
// key in this map, or nil if the map is empty.
func (sm *SkipListMap) LastEntry() *collection.Entry {
	var entry collection.Entry = sm.findLast()
	return &entry
}

// PollFirstEntry removes and returns a key-value mapping associated with
// the least key in this map, or nil if the map is empty.
func (sm *SkipListMap) PollFirstEntry() *collection.Entry {
	var entry collection.Entry = sm.doRemoveFirstEntry()
	return &entry
}

// PollLastEntry removes and returns a key-value mapping associated with
// the greatest key in this map, or null if the map is empty.
func (sm *SkipListMap) PollLastEntry() *collection.Entry {
	var entry collection.Entry = sm.doRemoveLastEntry()
	return &entry
}

func (sm *SkipListMap) checkNil(p *collection.Object) bool {
	return p == nil || (*p) == nil
}

func (sm *SkipListMap) checkKeyType(p *collection.Object) bool {
	return reflect.TypeOf(*p) == sm.kt
}

func (sm *SkipListMap) checkValueType(p *collection.Object) bool {
	return reflect.TypeOf(*p) == sm.vt
}

func (sm *SkipListMap) doGet(key *collection.Object) *collection.Object {
	node := sm.findPredecessor(key)
	for node != nil && !reflect.DeepEqual(*node.key, *key) {
		if !sm.precede(node.key, key) {
			return nil
		}
		node = node.next
	}
	if node == nil {
		return nil
	} else {
		return node.value
	}
}

func (sm *SkipListMap) doPut(key *collection.Object, value *collection.Object) *collection.Object {
	prev := sm.findPredecessor(key)
	node := prev.next
	// Add node
	for node != nil && !sm.precede(key, node.key) {
		if reflect.DeepEqual(*key, *node.key) {
			v := node.value
			node.value = value
			return v
		}
		prev = node
		node = node.next
	}
	z := &Node{key, value, node}
	prev.next = z
	sm.size++
	// Add Index
	isIndex, level := sm.indexAndLevel(rand.Int31())
	if isIndex {
		// build Index list
		var idx *Index
		if level < sm.head.level {
			for i := 1; i <= level; i++ {
				idx = &Index{i, z, idx, nil}
			}
		} else {
			level = sm.head.level + 1
			for i := 1; i <= level; i++ {
				idx = &Index{i, z, idx, nil}
			}
			newHead := &Index{level, sm.head.node, sm.head, idx}
			sm.head = newHead
		}
		// find insertion points and splice in
		h := sm.head
		if sm.head.level == level {
			h = h.down
			idx = idx.down
		} else {
			for curr := sm.head.level; curr > level; curr-- {
				h = h.down
			}
		}
		for h != nil {
			p := h
			r := p.right
			for r != nil && !sm.precede(key, key) {
				p = r
				r = p.right
			}
			// insert
			p.right = idx
			idx.right = r
			// down
			h = h.down
			idx = idx.down
		}
	}
	return nil
}

func (sm *SkipListMap) doRemove(key *collection.Object) *Node {
	// remove node
	prev := sm.findPredecessor(key)
	node := prev.next
	for node != nil && !reflect.DeepEqual(*node.key, *key) {
		if !sm.precede(node.key, key) {
			return nil
		}
		prev = node
		node = node.next
	}
	if node == nil {
		return nil
	}
	prev.next = node.next
	node.next = node
	sm.size--
	// remove index
	h := sm.head
	for h != nil {
		p := h
		r := p.right
		k := key
		for r != nil && !reflect.DeepEqual(*key, *k) {
			p = r
			r = p.right
		}
		if r != nil {
			// remove Index
			p.right = r.right
			r.right = nil
		}
		// down
		h = h.down
	}
	// try reduce level
	sm.tryReduceLevel()
	return node
}

func (sm *SkipListMap) findFirst() *Node {
	return sm.head.node.next
}

func (sm *SkipListMap) findLast() *Node {
	h := sm.head
	for h.down != nil {
		for h.right != nil {
			h = h.right
		}
		h = h.down
	}
	n := h.node
	if n == nil {
		return nil
	}
	for n.next != nil {
		n = n.next
	}
	return n
}

func (sm *SkipListMap) findNear(key *collection.Object, rel int) *Node {
	prev := sm.findPredecessor(key)
	node := prev.next
	for {
		if node == nil {
			if (rel & LT) == 0 || prev.value == nil {
				return nil
			} else {
				return prev
			}
		}
		c := sm.compare(key, node.key)
		if (c == 0 && (rel & EQ) != 0) || c < 0 && (rel & LT) == 0 {
			return node
		}
		if c <= 0 && (rel & LT) != 0 {
			if prev.value == nil {
				return nil
			} else {
				return prev
			}
		}
		prev = node
		node = node.next
	}
}

func (sm *SkipListMap) doRemoveFirstEntry() *Node {
	node := sm.findFirst()
	if node == nil {
		return nil
	}
	sm.doRemove(node.key)
	return node
}

func (sm *SkipListMap) doRemoveLastEntry() *Node {
	node := sm.findFirst()
	if node == nil {
		return nil
	}
	sm.doRemove(node.key)
	return node
}

// Returns a base-level node with key strictly less than given key, or the
// base-level header if there is no such node.
func (sm *SkipListMap) findPredecessor(key *collection.Object) *Node {
	q := sm.head
	r := q.right
	for {
		if r != nil {
			n := r.node
			if sm.precede(n.key, key) {
				q = r
				r = q.right
				continue
			}
		}
		d := q.down
		if d == nil {
			return q.node
		}
		q = d
		r = q.right
	}
}

func (sm *SkipListMap) indexAndLevel(k int32) (bool, int) {
	index := (k & 1) == 0
	k >>= 1
	var level int
	for k&1 != 0 {
		k = k >> 1
		level++
	}
	return index, level
}

func (sm *SkipListMap) tryReduceLevel() {
	h := sm.head
	for h.down != nil && h.right == nil {
		h, h.down = h.down, nil
	}
}

func (sm *SkipListMap) compare(p1 *collection.Object, p2 *collection.Object) int {
	if reflect.DeepEqual(*p1, *p2) {
		return 0
	}
	if sm.precede(p1, p2) {
		return -1
	} else {
		return 1
	}
}

// HashNext returns true if the iteration has more elements.
func (it *EntryIterator) HashNext() bool {
	return it.next != nil
}

// Next returns the next element in the iteration.
func (it *EntryIterator) Next() *collection.Entry {
	if it.HashNext() {
		it.lastRet = it.next
		it.next = it.next.next
		var t collection.Entry =  it.lastRet
		return &t
	}
	return nil
}

// Remove removes from the underlying collection the last element returned
// by this iterator.
func (it *EntryIterator) Remove() (*collection.Entry, bool) {
	if it.lastRet == nil {
		return nil, false
	}
	var last collection.Entry = it.lastRet
	lastNext := it.lastRet.next
	it.sm.doRemove(last.GetKey())
	if it.next == it.lastRet {
		it.next = lastNext
	}
	it.lastRet = nil
	return &last, true
}