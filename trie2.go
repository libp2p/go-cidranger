package cidranger

import (
	"fmt"
	"net"

	rnet "github.com/libp2p/go-cidranger/net"
	t2 "github.com/libp2p/go-libp2p-xor/trie2"
)

type cidrEntry struct {
	net.IPNet
}

func (e cidrEntry) Network() net.IPNet {
	return e.IPNet
}

type rangerKey struct {
	IP    net.IP
	Size  int
	Entry RangerEntry
}

func rangerEntryToKey(e RangerEntry) rangerKey {
	s, _ := e.Network().Mask.Size()
	return rangerKey{IP: e.Network().IP, Size: s, Entry: e}
}

func ipToKey(ip net.IP) rangerKey {
	return rangerKey{IP: ip, Size: len(ip) * 8}
}

func ipNetToKey(ipNet net.IPNet) rangerKey {
	return rangerEntryToKey(NewBasicRangerEntry(ipNet))
}

func (k rangerKey) String() string {
	return fmt.Sprintf("%v/%v", k.IP, k.Size)
}

func (k rangerKey) Equal(r t2.Key) bool {
	if k2, ok := r.(rangerKey); ok {
		if k.Len() != k2.Len() {
			return false
		} else {
			return commonPrefixLen(k.IP, k2.IP) >= k.Len()
		}
	} else {
		return false
	}
}

func commonPrefixLen(a, b []byte) (cpl int) {
	if len(a) > len(b) {
		a = a[:len(b)]
	}
	if len(b) > len(a) {
		b = b[:len(a)]
	}
	for len(a) > 0 {
		if a[0] == b[0] {
			cpl += 8
			a = a[1:]
			b = b[1:]
			continue
		}
		bits := 8
		ab, bb := a[0], b[0]
		for {
			ab >>= 1
			bb >>= 1
			bits--
			if ab == bb {
				cpl += bits
				return
			}
		}
	}
	return
}

func (k rangerKey) BitAt(i int) byte {
	b := []byte(k.IP)
	// the most significant byte in an IP address is the first one
	d := b[i/8] & (byte(1) << (7 - (i % 8)))
	if d == 0 {
		return 0
	} else {
		return 1
	}
}

func (k rangerKey) Len() int {
	return k.Size
}

type trie2Ranger struct {
	trie *t2.Trie
}

func newTrie2Ranger(v rnet.IPVersion) Ranger {
	return &trie2Ranger{trie: &t2.Trie{}}
}

func (r *trie2Ranger) Insert(entry RangerEntry) error {
	r.trie.Add(rangerEntryToKey(entry))
	return nil
}

func (r *trie2Ranger) Remove(network net.IPNet) (RangerEntry, error) {
	panic("not supported")
}

func (r *trie2Ranger) Contains(ip net.IP) (bool, error) {
	if c, err := r.ContainingNetworks(ip); err != nil {
		return false, err
	} else {
		return len(c) > 0, nil
	}
}

func (r *trie2Ranger) ContainingNetworks(ip net.IP) ([]RangerEntry, error) {
	_, found := r.trie.FindSubKeys(ipToKey(ip))
	q := []RangerEntry{}
	for _, f := range found {
		e := f.(rangerKey).Entry
		pn := e.Network()
		if pn.Contains(ip) {
			q = append(q, e)
		}
	}
	return q, nil
}

func (r *trie2Ranger) CoveredNetworks(network net.IPNet) ([]RangerEntry, error) {
	panic("not supported")
}

func (r *trie2Ranger) Len() int {
	return r.trie.Size()
}
