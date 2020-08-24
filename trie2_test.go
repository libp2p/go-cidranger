package cidranger

import (
	"net"
	"testing"

	rnet "github.com/libp2p/go-cidranger/net"
)

func TestBasicTrie2(t *testing.T) {
	r := newTrie2Ranger(rnet.IPVersion(""))
	_, nn1, _ := net.ParseCIDR("178.236.0.0/20")
	r.Insert(NewBasicRangerEntry(*nn1))

	ip1 := net.IPv4(178, 236, 15, 124)
	if ok, err := r.Contains(ip1); err != nil {
		t.Errorf("contains (%v)", err)
	} else if !ok {
		t.Errorf("expected key missing")
	}

	ip2 := net.IPv4(178, 236, 255, 124)
	if ok, err := r.Contains(ip2); err != nil {
		t.Errorf("contains (%v)", err)
	} else if ok {
		t.Errorf("expected to not have key")
	}
}
