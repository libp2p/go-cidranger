package cidranger

import (
	"fmt"
	"net"
	"testing"

	rnet "github.com/libp2p/go-cidranger/net"
)

func TestBasicTrie2(t *testing.T) {
	r := newTrie2Ranger(rnet.IPVersion(""))
	_, nn1, _ := net.ParseCIDR("1.2.3.4/1")
	r.Insert(NewBasicRangerEntry(*nn1))

	ip1 := net.IPv4(1, 2, 3, 4)
	if ok, err := r.Contains(ip1); err != nil {
		t.Errorf("contains (%v)", err)
	} else if !ok {
		t.Errorf("expected key missing")
	}

	ip2 := net.IPv4(255, 2, 3, 4)
	if ok, err := r.Contains(ip2); err != nil {
		t.Errorf("contains (%v)", err)
	} else if ok {
		t.Errorf("expected to not have key")
	}

}

func TestIpKeyBitAt(t *testing.T) {
	// ip, _, err := net.ParseCIDR("34.208.0.0/12")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	ip := net.IPv4(34, 208, 0, 0)
	fmt.Println([]byte(ip))
	// XXX "34.208.0.0/12"
	// XXX "34.192.0.0/12"
}
