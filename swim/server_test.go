package main

import "testing"

func TestJoin(t *testing.T) {
	var (
		clientAddr = ":3001"
		serverAddr = ":3000"
	)

	srv1 := NewServer(serverAddr)
	if err := srv1.Start(); err != nil {
		t.Fatal(err)
	}
	defer srv1.listener.Close()

	go func() {
		if err := srv1.Listen(); err != nil {
			t.Fatal(err)
		}
	}()

	srv2 := NewServer(clientAddr)
	if err := srv2.Start(); err != nil {
		t.Fatal(err)
	}
	defer srv2.listener.Close()

	if err := srv2.Join(serverAddr); err != nil {
		t.Fatal(err)
	}

	// Check first member
	if len(srv1.Members.Members) != 2 {
		t.Fatalf("unexpected member count: %d", len(srv1.Members.Members))
	}

	if _, ok := srv1.Members.Members[clientAddr]; !ok {
		t.Error(srv1.BindAddr, "is missing", clientAddr, "in memberlist")
	}

	if _, ok := srv1.Members.Members[serverAddr]; !ok {
		t.Error(srv1.BindAddr, "is missing", serverAddr, "in memberlist")
	}

	// Check second member
	if len(srv2.Members.Members) != 2 {
		t.Fatalf("unexpected member count: %d", len(srv2.Members.Members))
	}

	if _, ok := srv2.Members.Members[clientAddr]; !ok {
		t.Error(srv2.BindAddr, "is missing", clientAddr, "in memberlist")
	}

	if _, ok := srv2.Members.Members[serverAddr]; !ok {
		t.Error(srv2.BindAddr, "is missing", serverAddr, "in memberlist")
	}
}

func TestJoinThird(t *testing.T) {
	var (
		serverAddr       = ":3000"
		firstClientAddr  = ":3001"
		secondClientAddr = ":3002"
	)

	srv1 := NewServer(serverAddr)
	if err := srv1.Start(); err != nil {
		t.Fatal(err)
	}
	defer srv1.listener.Close()

	go func() {
		if err := srv1.Listen(); err != nil {
			t.Fatal(err)
		}
	}()

	srv2 := NewServer(firstClientAddr)
	if err := srv2.Start(); err != nil {
		t.Fatal(err)
	}
	defer srv2.listener.Close()

	if err := srv2.Join(serverAddr); err != nil {
		t.Fatal(err)
	}

	srv3 := NewServer(secondClientAddr)
	if err := srv3.Start(); err != nil {
		t.Fatal(err)
	}
	defer srv3.listener.Close()

	if err := srv3.Join(serverAddr); err != nil {
		t.Fatal(err)
	}

	// Check first member
	if len(srv1.Members.Members) != 3 {
		t.Errorf("unexpected member count: %d", len(srv1.Members.Members))
	}
	if _, ok := srv1.Members.Members[serverAddr]; !ok {
		t.Error(srv1.BindAddr, "is missing", serverAddr, "in memberlist")
	}
	if _, ok := srv1.Members.Members[firstClientAddr]; !ok {
		t.Error(srv1.BindAddr, "is missing", firstClientAddr, "in memberlist")
	}
	if _, ok := srv1.Members.Members[secondClientAddr]; !ok {
		t.Error(srv1.BindAddr, "is missing", secondClientAddr, "in memberlist")
	}

	// Check second member
	if len(srv2.Members.Members) != 2 {
		t.Errorf("unexpected member count: %d", len(srv2.Members.Members))
	}
	if _, ok := srv2.Members.Members[serverAddr]; !ok {
		t.Error(srv2.BindAddr, "is missing", serverAddr, "in memberlist")
	}
	if _, ok := srv2.Members.Members[firstClientAddr]; !ok {
		t.Error(srv2.BindAddr, "is missing", firstClientAddr, "in memberlist")
	}

	// srv2 should not contain srv3 since srv3 joined after last contact between
	// srv1 and srv2.
	if _, ok := srv2.Members.Members[secondClientAddr]; ok {
		t.Error(srv2.BindAddr, "contains unexpected", secondClientAddr, "in memberlist")
	}

	// Check third member
	if len(srv3.Members.Members) != 3 {
		t.Errorf("unexpected member count: %d", len(srv2.Members.Members))
	}
	if _, ok := srv3.Members.Members[serverAddr]; !ok {
		t.Error(srv3.BindAddr, "is missing", serverAddr, "in memberlist")
	}
	if _, ok := srv3.Members.Members[firstClientAddr]; !ok {
		t.Error(srv3.BindAddr, "is missing", firstClientAddr, "in memberlist")
	}
	if _, ok := srv3.Members.Members[secondClientAddr]; !ok {
		t.Error(srv3.BindAddr, "is missing", secondClientAddr, "in memberlist")
	}
}

func TestPing(t *testing.T) {
	var (
		serverAddr      = ":3000"
		firstClientAddr = ":3001"
	)

	srv1 := NewServer(serverAddr)

	if err := srv1.Start(); err != nil {
		t.Fatal(err)
	}
	defer srv1.listener.Close()

	go func() {
		if err := srv1.Listen(); err != nil {
			t.Fatal(err)
		}
	}()

	srv2 := NewServer(firstClientAddr)
	if err := srv2.Start(); err != nil {
		t.Fatal(err)
	}
	defer srv2.listener.Close()

	if err := srv2.Join(serverAddr); err != nil {
		t.Fatal(err)
	}

	srv2.Members.Add(Member{
		Name:    "test",
		Address: ":3003",
	})

	srv2.Ping(Member{
		Name:    serverAddr,
		Address: serverAddr,
	})

	// Check first member
	if len(srv1.Members.Members) != 3 {
		t.Errorf("unexpected member count: %d", len(srv1.Members.Members))
	}
	if _, ok := srv1.Members.Members[serverAddr]; !ok {
		t.Error(srv1.BindAddr, "is missing", serverAddr, "in memberlist")
	}
	if _, ok := srv1.Members.Members[firstClientAddr]; !ok {
		t.Error(srv1.BindAddr, "is missing", firstClientAddr, "in memberlist")
	}
	if len(srv1.Members.Updates) != 3 {
		t.Errorf("unexpected update count: %d", len(srv1.Members.Updates))
	}

	// Check second member
	if len(srv2.Members.Members) != 3 {
		t.Errorf("unexpected member count: %d", len(srv2.Members.Members))
	}
	if _, ok := srv2.Members.Members[serverAddr]; !ok {
		t.Error(srv2.BindAddr, "is missing", serverAddr, "in memberlist")
	}
	if _, ok := srv2.Members.Members[firstClientAddr]; !ok {
		t.Error(srv2.BindAddr, "is missing", firstClientAddr, "in memberlist")
	}
	if len(srv2.Members.Updates) != 0 {
		t.Errorf("unexpected update count: %d", len(srv2.Members.Updates))
	}

}
