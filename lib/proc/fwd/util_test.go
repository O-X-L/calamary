package fwd

import (
	"fmt"
	"net"
	"testing"

	"github.com/creasty/defaults"
	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/proc/parse"
)

func TestTargetReachable(t *testing.T) {
	cnf.C = &cnf.Config{}
	err := defaults.Set(cnf.C) // probe timeout
	if err != nil {
		fmt.Printf("Failed to set default probe-timeout: %v", err)
	}

	testTarget := net.ParseIP("1.1.1.1")
	if !TargetIsReachable("tcp", testTarget, 53) {
		t.Error("Target reachability check failed #1")
	}
	if TargetIsReachable("tcp", testTarget, 50000) {
		t.Error("Target reachability check failed #2")
	}
}

func TestFirstReachableTarget(t *testing.T) {
	cnf.C = &cnf.Config{}
	err := defaults.Set(cnf.C) // probe timeout
	if err != nil {
		fmt.Printf("Failed to set default probe-timeout: %v", err)
	}

	testTarget1 := net.ParseIP("135.181.170.219")
	testTarget2 := net.ParseIP("1.1.1.1")
	testTarget3 := net.ParseIP("2001:db8::1")
	testTarget4 := net.ParseIP("2606:4700:4700::1001")

	testerHasIPv6 := TargetIsReachable("tcp", testTarget4, 53)

	target := FirstReachableTarget(
		parse.ParsedPacket{},
		[]net.IP{testTarget1},
		[]net.IP{},
		"tcp",
		53,
	)
	if target != nil {
		t.Errorf("Target first-reachable check failed #1 (%v)", target)
	}

	target = FirstReachableTarget(
		parse.ParsedPacket{},
		[]net.IP{testTarget1, testTarget2},
		[]net.IP{},
		"tcp",
		53,
	)
	if target.String() != testTarget2.String() {
		t.Errorf("Target first-reachable check failed #2 (%v)", target)
	}

	target = FirstReachableTarget(
		parse.ParsedPacket{},
		[]net.IP{},
		[]net.IP{testTarget3, testTarget4},
		"tcp",
		53,
	)
	if testerHasIPv6 && target.String() != testTarget4.String() {
		// nil => if testing client has no IPv6
		t.Errorf("Target first-reachable check failed #3 (%v)", target)
	}

	target = FirstReachableTarget(
		parse.ParsedPacket{},
		[]net.IP{testTarget1, testTarget2},
		[]net.IP{testTarget3, testTarget4},
		"tcp",
		53,
	)
	if testerHasIPv6 && target.String() != testTarget4.String() {
		// nil => if testing client has no IPv6
		t.Errorf("Target first-reachable check failed #4 (%v)", target)
	}
}
