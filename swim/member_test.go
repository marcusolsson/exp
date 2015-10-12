package main

import "testing"

func TestMemberList_NewList(t *testing.T) {
	rounds := 3
	l := NewList(rounds)

	if len(l.Members) != 0 {
		t.Errorf("len(l.Members) = %d; want = %d", len(l.Members), 0)
	}
	if len(l.Failed) != 0 {
		t.Errorf("len(l.Failed) = %d; want = %d", len(l.Failed), 0)
	}
	if len(l.Updates) != 0 {
		t.Errorf("len(l.Updates) = %d; want = %d", len(l.Updates), 0)
	}
}

func TestMemberList_AddMember(t *testing.T) {
	mem := Member{
		Name:    "test_name",
		Address: "test_addr",
	}

	rounds := 3
	l := NewList(rounds)
	l.Add(mem)

	if len(l.Members) != 1 {
		t.Errorf("len(l.Members) = %d; want = %d", len(l.Members), 1)
	}
	if len(l.Failed) != 0 {
		t.Errorf("len(l.Failed) = %d; want = %d", len(l.Failed), 0)
	}
	if len(l.Updates) != 1 {
		t.Errorf("len(l.Updates) = %d; want = %d", len(l.Updates), 1)
	}

	m, ok := l.Members[mem.Address]
	if !ok {
		t.Errorf("missing member")
	}
	if m != mem {
		t.Errorf("m = %v; want = %v", m, mem)
	}

	want := Update{Member: mem, Round: 0, Type: Joined}
	if l.Updates[0] != want {
		t.Errorf("l.Updates[0] = %v; want = %v", l.Updates[0], want)
	}
}

func TestMemberList_RemoveMember(t *testing.T) {
	mem := Member{
		Name:    "test_name",
		Address: "test_addr",
	}

	rounds := 3
	l := NewList(rounds)
	l.Add(mem)
	l.Remove(mem)

	if len(l.Members) != 0 {
		t.Errorf("len(l.Members) = %d; want = %d", len(l.Members), 0)
	}
	if len(l.Failed) != 1 {
		t.Errorf("len(l.Failed) = %d; want = %d", len(l.Failed), 1)
	}
	if len(l.Updates) != 2 {
		t.Errorf("len(l.Updates) = %d; want = %d", len(l.Updates), 2)
	}

	// TODO(marcusolsson): This update is irrelevant(?) and should probably be removed.
	want1 := Update{Member: mem, Round: 0, Type: Joined}
	if l.Updates[0] != want1 {
		t.Errorf("l.Updates[0] = %v; want = %v", l.Updates[0], want1)
	}
	want2 := Update{Member: mem, Round: 0, Type: Failed}
	if l.Updates[1] != want2 {
		t.Errorf("l.Updates[1] = %v; want = %v", l.Updates[1], want2)
	}
}

func TestMemberList_AddRemovedMember(t *testing.T) {
	mem := Member{
		Name:    "test_name",
		Address: "test_addr",
	}

	rounds := 3
	l := NewList(rounds)
	l.Add(mem)
	l.Remove(mem)
	l.Add(mem)

	if len(l.Members) != 1 {
		t.Errorf("len(l.Members) = %d; want = %d", len(l.Members), 1)
	}
	if len(l.Failed) != 0 {
		t.Errorf("len(l.Failed) = %d; want = %d", len(l.Failed), 0)
	}
	if len(l.Updates) != 3 {
		t.Errorf("len(l.Updates) = %d; want = %d", len(l.Updates), 3)
	}

	// TODO(marcusolsson): The first two updates are irrelevant(?) and should
	// probably be removed.
	want1 := Update{Member: mem, Round: 0, Type: Joined}
	if l.Updates[0] != want1 {
		t.Errorf("l.Updates[0] = %v; want = %v", l.Updates[0], want1)
	}
	want2 := Update{Member: mem, Round: 0, Type: Failed}
	if l.Updates[1] != want2 {
		t.Errorf("l.Updates[1] = %v; want = %v", l.Updates[1], want2)
	}
	want3 := Update{Member: mem, Round: 0, Type: Joined}
	if l.Updates[2] != want3 {
		t.Errorf("l.Updates[2] = %v; want = %v", l.Updates[2], want3)
	}
}

func TestMemberList_MergeAdded(t *testing.T) {
	up := Update{
		Member: Member{
			Name:    "test_name",
			Address: "test_addr",
		},
		Round: 0,
		Type:  Joined,
	}

	rounds := 3
	l := NewList(rounds)

	l.Merge([]Update{up})

	if len(l.Members) != 1 {
		t.Errorf("len(l.Members) = %d; want = %d", len(l.Members), 1)
	}
	if len(l.Failed) != 0 {
		t.Errorf("len(l.Failed) = %d; want = %d", len(l.Failed), 0)
	}
	if len(l.Updates) != 1 {
		t.Errorf("len(l.Updates) = %d; want = %d", len(l.Updates), 1)
	}

	want := up
	if l.Updates[0] != want {
		t.Errorf("l.Updates[0] = %v; want = %v", l.Updates[0], want)
	}
}
