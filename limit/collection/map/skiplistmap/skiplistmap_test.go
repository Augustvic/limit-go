package skiplistmap

import (
	"LimitGo/limit/collection"
	"LimitGo/limit/collection/linear/hashset"
	"reflect"
	"strconv"
	"testing"
)

type Student struct {
	Id int
	Name string
}

type Teacher struct {
	Id int
	Name string
	Sex int
}

var precede = func(p1 *collection.Object, p2 *collection.Object) bool {
	s1 := (*p1).(Teacher)
	s2 := (*p2).(Teacher)
	return s1.Id < s2.Id
}

var t2 collection.Object
var t4 collection.Object
var t6 collection.Object
var t8 collection.Object
var st11 collection.Object
var st12 collection.Object
var st2 collection.Object
var st3 collection.Object
var st4 collection.Object
var s1 collection.Object
var s2 collection.Object
var s3 collection.Object
var s4 collection.Object
var m *SkipListMap

func RestartSkipListMap() {
	t2 = Teacher{2, "t2", 0}
	t4 = Teacher{4, "t4", 0}
	t6 = Teacher{6, "t6", 0}
	t8 = Teacher{8, "t8", 0}
	st11 = Student{100, "st11"}
	st12 = Student{101, "st12"}
	st2 = Student{102, "st2"}
	st3 = Student{103, "st3"}
	st4 = Student{104, "st4"}
	s1 = *hashset.New()
	s2 = *hashset.New()
	s3 = *hashset.New()
	s4 = *hashset.New()
	s1s := s1.(hashset.HashSet)
	s2s := s2.(hashset.HashSet)
	s3s := s3.(hashset.HashSet)
	s4s := s4.(hashset.HashSet)
	s1s.Add(&st11)
	s1s.Add(&st12)
	s2s.Add(&st2)
	s3s.Add(&st3)
	s4s.Add(&st4)
	m = New(precede)
	m.Put(&t2, &s1)
	m.Put(&t4, &s2)
	m.Put(&t6, &s3)
	m.Put(&t8, &s4)
}

func TestSkipListMapAll(t *testing.T) {
	TestSkipListMap_CeilingEntry(t)
	TestSkipListMap_Clear(t)
	TestSkipListMap_ContainsKey(t)
	TestSkipListMap_ContainsValue(t)
	TestSkipListMap_Empty(t)
	TestSkipListMap_EntrySet(t)
	TestSkipListMap_Equals(t)
	TestSkipListMap_FirstEntry(t)
	TestSkipListMap_FloorEntry(t)
	TestSkipListMap_Get(t)
	TestSkipListMap_GetEntryIterator(t)
	TestSkipListMap_HeadMap(t)
	TestSkipListMap_HigherEntry(t)
	TestSkipListMap_KeySet(t)
	TestSkipListMap_LastEntry(t)
	TestSkipListMap_LowerEntry(t)
	TestSkipListMap_PollFirstEntry(t)
	TestSkipListMap_PollLastEntry(t)
	TestSkipListMap_Put(t)
	TestSkipListMap_PutAll(t)
	TestSkipListMap_Remove(t)
	TestSkipListMap_String(t)
	TestSkipListMap_SubMap(t)
	TestSkipListMap_TailMap(t)
	TestSkipListMap_Values(t)
}

func TestSkipListMap_CeilingEntry(t *testing.T) {
	RestartSkipListMap()
	var temp1 collection.Object = Teacher{4, "t4", 0}
	var temp2 collection.Object = Teacher{5, "t5", 0}
	var temp3 collection.Object = Teacher{8, "t8", 0}
	var temp4 collection.Object = Teacher{10, "t10", 0}
	k1 := (*m.CeilingEntry(&temp1)).GetKey()
	if *k1 != t4 {
		t.Error("CeilingEntry operation fail!")
	}
	k2 := (*m.CeilingEntry(&temp2)).GetKey()
	if *k2 != t6 {
		t.Error("CeilingEntry operation fail!")
	}
	k3 := (*m.CeilingEntry(&temp3)).GetKey()
	if *k3 != t8 {
		t.Error("CeilingEntry operation fail!")
	}
	k4 := m.CeilingEntry(&temp4)
	if k4 != nil {
		t.Error("CeilingEntry operation fail!")
	}
}

func TestSkipListMap_Clear(t *testing.T) {
	RestartSkipListMap()
	if m.Size() == 0 {
		t.Error("Start operation fail!")
	}
	m.Clear()
	if m.Size() != 0 {
		t.Error("Clear operation fail!")
	}
}

func TestSkipListMap_ContainsKey(t *testing.T) {
	RestartSkipListMap()
	var temp1 collection.Object = Teacher{4, "t4", 0}
	var temp2 collection.Object = Teacher{5, "t5", 0}
	if !m.ContainsKey(&temp1) {
		t.Error("ContainsKey operation fail!")
	}
	if m.ContainsKey(&temp2) {
		t.Error("ContainsKey operation fail!")
	}
	if m.ContainsKey(nil) {
		t.Error("ContainsKey operation fail!")
	}
}

func TestSkipListMap_ContainsValue(t *testing.T) {
	RestartSkipListMap()
	var st1t collection.Object = Student{100, "st11"}
	var st2t collection.Object = Student{101, "st12"}
	s1t := *hashset.New()
	s1t.Add(&st1t)
	s1t.Add(&st2t)
	var s1to collection.Object = s1t
	if !m.ContainsValue(&s1to) {
		t.Error("ContainsValue operation fail!")
	}
}

func TestSkipListMap_Empty(t *testing.T) {
	RestartSkipListMap()
	if m.Empty() {
		t.Error("Start operation fail!")
	}
	m.Clear()
	if !m.Empty() {
		t.Error("Empty operation fail!")
	}
}

func TestSkipListMap_EntrySet(t *testing.T) {
	RestartSkipListMap()
	es := m.EntrySet()
	if (*es).Size() != 4 {
		t.Error("EntrySet operation fail!")
	}
}

func TestSkipListMap_Equals(t *testing.T) {
	RestartSkipListMap()
	var tt2 collection.Object
	var tt4 collection.Object
	var tt6 collection.Object
	var tt8 collection.Object
	var tst11 collection.Object
	var tst12 collection.Object
	var tst2 collection.Object
	var tst3 collection.Object
	var tst4 collection.Object
	var ts1 collection.Object
	var ts2 collection.Object
	var ts3 collection.Object
	var ts4 collection.Object
	var tm *SkipListMap
	tt2 = Teacher{2, "t2", 0}
	tt4 = Teacher{4, "t4", 0}
	tt6 = Teacher{6, "t6", 0}
	tt8 = Teacher{8, "t8", 0}
	tst11 = Student{100, "st11"}
	tst12 = Student{101, "st12"}
	tst2 = Student{102, "st2"}
	tst3 = Student{103, "st3"}
	tst4 = Student{104, "st4"}
	ts1 = *hashset.New()
	ts2 = *hashset.New()
	ts3 = *hashset.New()
	ts4 = *hashset.New()
	ts1s := ts1.(hashset.HashSet)
	ts2s := ts2.(hashset.HashSet)
	ts3s := ts3.(hashset.HashSet)
	ts4s := ts4.(hashset.HashSet)
	ts1s.Add(&tst11)
	ts1s.Add(&tst12)
	ts2s.Add(&tst2)
	ts3s.Add(&tst3)
	ts4s.Add(&tst4)
	tm = New(precede)
	tm.Put(&tt2, &ts1)
	tm.Put(&tt4, &ts2)
	tm.Put(&tt6, &ts3)
	tm.Put(&tt8, &ts4)
	var ma collection.Map = tm
	if !m.Equals(&ma) {
		t.Error("Equals operation fail!")
	}
}

func TestSkipListMap_FirstEntry(t *testing.T) {
	RestartSkipListMap()
	key := (*m.FirstEntry()).GetKey()
	if *key != t2 {
		t.Error("FirstEntry operation fail!")
	}
}

func TestSkipListMap_FloorEntry(t *testing.T) {
	RestartSkipListMap()
	var temp1 collection.Object = Teacher{1, "t1", 0}
	var temp2 collection.Object = Teacher{5, "t5", 0}
	var temp3 collection.Object = Teacher{8, "t8", 0}
	var temp4 collection.Object = Teacher{11, "t10", 0}
	k1 := m.FloorEntry(&temp1)
	if k1 != nil {
		t.Error("FloorEntry operation fail!")
	}
	k2 := (*m.FloorEntry(&temp2)).GetKey()
	if *k2 != t4 {
		t.Error("FloorEntry operation fail!")
	}
	k3 := (*m.FloorEntry(&temp3)).GetKey()
	if *k3 != t8 {
		t.Error("FloorEntry operation fail!")
	}
	k4 := (*m.FloorEntry(&temp4)).GetKey()
	if *k4 != t8 {
		t.Error("FloorEntry operation fail!")
	}
}

func TestSkipListMap_Get(t *testing.T) {
	RestartSkipListMap()
	oj := m.Get(&t4)
	s := (*oj).(hashset.HashSet)
	if s.Size() != 1 || !s.Contains(&st2) {
		t.Error("Get operation fail!")
	}
}

func TestSkipListMap_GetEntryIterator(t *testing.T) {
	RestartSkipListMap()
	it := m.GetEntryIterator()
	index := 2
	s := ""
	for i := 0; it.HashNext(); i++ {
		if i == index {
			it.Remove()
		} else {
			entry := it.Next()
			teacher := (*(*entry).GetKey()).(Teacher)
			k := strconv.Itoa(teacher.Id)
			s += k
		}
	}
	if m.Size() != 3 || s != "2468" {
		t.Error("GetEntryIterator operation fail!")
	}
}

func TestSkipListMap_HeadMap(t *testing.T) {
	RestartSkipListMap()
	sm := m.HeadMap(&t2, true)
	if (*sm).Size() != 1 || (*(*(*sm).FirstEntry()).GetKey()) != t2 ||  (*(*(*sm).LastEntry()).GetKey()) != t2 {
		t.Error("HeadMap operation fail!")
	}
}

func TestSkipListMap_HigherEntry(t *testing.T) {
	RestartSkipListMap()
	var temp1 collection.Object = Teacher{4, "t4", 0}
	var temp2 collection.Object = Teacher{5, "t5", 0}
	var temp3 collection.Object = Teacher{8, "t8", 0}
	k1 := (*m.HigherEntry(&temp1)).GetKey()
	if *k1 != t6 {
		t.Error("HigherEntry operation fail!")
	}
	k2 := (*m.HigherEntry(&temp2)).GetKey()
	if *k2 != t6 {
		t.Error("HigherEntry operation fail!")
	}
	k3 := m.HigherEntry(&temp3)
	if k3 != nil {
		t.Error("HigherEntry operation fail!")
	}
}

func TestSkipListMap_KeySet(t *testing.T) {
	RestartSkipListMap()
	ks := m.KeySet()
	if (*ks).String() != "{{\"Id\":2,\"Name\":\"t2\",\"Sex\":0},{\"Id\":4,\"Name\":\"t4\",\"Sex\":0},{\"Id\":6,\"Name\":\"t6\",\"Sex\":0},{\"Id\":8,\"Name\":\"t8\",\"Sex\":0}}" {
		t.Error("KeySet operation fail!")
	}
}

func TestSkipListMap_LastEntry(t *testing.T) {
	RestartSkipListMap()
	key := (*m.LastEntry()).GetKey()
	if *key != t8 {
		t.Error("LastEntry operation fail!")
	}
}

func TestSkipListMap_LowerEntry(t *testing.T) {
	RestartSkipListMap()
	var temp1 collection.Object = Teacher{1, "t1", 0}
	var temp2 collection.Object = Teacher{5, "t5", 0}
	var temp3 collection.Object = Teacher{8, "t8", 0}
	var temp4 collection.Object = Teacher{11, "t10", 0}
	k1 := m.LowerEntry(&temp1)
	if k1 != nil {
		t.Error("LowerEntry operation fail!")
	}
	k2 := (*m.LowerEntry(&temp2)).GetKey()
	if *k2 != t4 {
		t.Error("LowerEntry operation fail!")
	}
	k3 := (*m.LowerEntry(&temp3)).GetKey()
	if *k3 != t6 {
		t.Error("LowerEntry operation fail!")
	}
	k4 := (*m.LowerEntry(&temp4)).GetKey()
	if *k4 != t8 {
		t.Error("LowerEntry operation fail!")
	}
}

func TestSkipListMap_PollFirstEntry(t *testing.T) {
	RestartSkipListMap()
	m.PollFirstEntry()
	key := (*m.FirstEntry()).GetKey()
	if m.Size() != 3 || *key != t4 {
		t.Error("PollFirstEntry operation fail!")
	}
}

func TestSkipListMap_PollLastEntry(t *testing.T) {
	RestartSkipListMap()
	m.PollLastEntry()
	key := (*m.LastEntry()).GetKey()
	if m.Size() != 3 || *key != t6 {
		t.Error("PollLastEntry operation fail!")
	}
}

func TestSkipListMap_Put(t *testing.T) {
	RestartSkipListMap()
	var t5 collection.Object = Teacher{5, "t5", 0}
	var t10 collection.Object = Teacher{10, "t10", 0}
	var st5 collection.Object = Student{102, "st2"}
	var st6 collection.Object = Student{103, "st3"}
	s5 := *hashset.New()
	s6 := *hashset.New()
	s5.Add(&st5)
	s6.Add(&st6)
	var s5o collection.Object = s5
	var s6o collection.Object = s6
	m.Put(&t5, &s5o)
	m.Put(&t10, &s6o)
	m.Put(&t2, &s5o)
	it := m.GetEntryIterator()
	s := ""
	for i := 0; it.HashNext(); i++ {
		entry := it.Next()
		teacher := (*(*entry).GetKey()).(Teacher)
		k := strconv.Itoa(teacher.Id)
		s += k
	}
	if m.Size() != 6 || s != "2456810" {
		t.Error("Put operation fail!")
	}
	v := *m.Get(&t2)
	if !reflect.DeepEqual(v, s5) {
		t.Error("Put operation fail!")
	}
}

func TestSkipListMap_PutAll(t *testing.T) {
	RestartSkipListMap()
	var t5 collection.Object = Teacher{5, "t5", 0}
	var t10 collection.Object = Teacher{10, "t10", 0}
	var st5 collection.Object = Student{102, "st2"}
	var st6 collection.Object = Student{103, "st3"}
	s5 := *hashset.New()
	s6 := *hashset.New()
	s5.Add(&st5)
	s6.Add(&st6)
	var s5o collection.Object = s5
	var s6o collection.Object = s6
	var m2 collection.Map = New(precede)
	m2.Put(&t5, &s5o)
	m2.Put(&t10, &s6o)
	m.PutAll(&m2)
	it := m.GetEntryIterator()
	s := ""
	for i := 0; it.HashNext(); i++ {
		entry := it.Next()
		teacher := (*(*entry).GetKey()).(Teacher)
		k := strconv.Itoa(teacher.Id)
		s += k
	}
	if m.Size() != 6 || s != "2456810" {
		t.Error("PutAll operation fail!")
	}
}

func TestSkipListMap_Remove(t *testing.T) {
	RestartSkipListMap()
	m.Remove(&t2)
	m.Remove(&t4)
	key := (*m.FirstEntry()).GetKey()
	if m.Size() != 2 || *key != t6 {
		t.Error("Remove operation fail!")
	}
}

func TestSkipListMap_String(t *testing.T) {
	RestartSkipListMap()
	if m.String() != "{{\"Id\":2,\"Name\":\"t2\",\"Sex\":0}={},{\"Id\":4,\"Name\":\"t4\",\"Sex\":0}={},{\"Id\":6,\"Name\":\"t6\",\"Sex\":0}={},{\"Id\":8,\"Name\":\"t8\",\"Sex\":0}={}}" {
		t.Error("String operation fail!")
	}
}

func TestSkipListMap_SubMap(t *testing.T) {
	RestartSkipListMap()
	sm := m.SubMap(nil, false, nil, false)
	if (*sm).Size() != 4 || (*(*(*sm).FirstEntry()).GetKey()) != t2 {
		t.Error("SubMap operation fail!")
	}
}

func TestSkipListMap_TailMap(t *testing.T) {
	RestartSkipListMap()
	sm := m.TailMap(&t4, true)
	if (*sm).Size() != 3 || (*(*(*sm).FirstEntry()).GetKey()) != t4 ||  (*(*(*sm).LastEntry()).GetKey()) != t8 {
		t.Error("TailMap operation fail!")
	}
}

func TestSkipListMap_Values(t *testing.T) {
	RestartSkipListMap()
	vs := m.Values()
	if (*vs).Size() != 4 {
		t.Error("Values operation fail!")
	}
}