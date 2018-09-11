package dms3fsaddr

import (
	"strings"
	"testing"

	peer "github.com/dms3-p2p/go-p2p-peer"
	ma "github.com/dms3-mft/go-multiaddr"
)

var good = []string{
	"/dms3fs/5dru6bJPUM1B7N69528u49DJiWZnok",
	"/dms3fs/kTRX47RthhwNzWdi6ggwqjuX",
	"/dms3fs/QmUCseQWXCSrhf9edzVKTvoj8o8Ts5aXFGNPameZRPJ6uR",
	"/ip4/1.2.3.4/tcp/1234/dms3fs/5dru6bJPUM1B7N69528u49DJiWZnok",
	"/ip4/1.2.3.4/tcp/1234/dms3fs/kTRX47RthhwNzWdi6ggwqjuX",
	"/ip4/1.2.3.4/tcp/1234/dms3fs/QmUCseQWXCSrhf9edzVKTvoj8o8Ts5aXFGNPameZRPJ6uR",
}

var bad = []string{
	"5dru6bJPUM1B7N69528u49DJiWZnok",                                // bad ma
	"kTRX47RthhwNzWdi6ggwqjuX",                                      // bad ma
	"QmUCseQWXCSrhf9edzVKTvoj8o8Ts5aXFGNPameZRPJ6uR",                // bad ma
	"dms3fs/5dru6bJPUM1B7N69528u49DJiWZnok",                           // bad ma
	"dms3fs/kTRX47RthhwNzWdi6ggwqjuX",                                 // bad ma
	"dms3fs/QmUCseQWXCSrhf9edzVKTvoj8o8Ts5aXFGNPameZRPJ6uR",           // bad ma
	"/dms3fs/5dru6bJPUM1B7N69528u49DJiWZno",                           // bad mh
	"/dms3fs/kTRX47RthhwNzWdi6ggwqju",                                 // bad mh
	"/dms3fs/QmUCseQWXCSrhf9edzVKTvj8o8Ts5aXFGNPameZRPJ6uR",           // bad mh
	"/dms3fs/QmUCseQWXCSrhf9edzVKTvoj8o8Ts5aXFGNPameZRPJ6uR/tcp/1234", // dms3fs not last
	"/ip4/1.2.3.4/tcp/dms3fs/5dru6bJPUM1B7N69528u49DJiWZnok",          // bad tcp part
	"/ip4/tcp/1234/dms3fs/kTRX47RthhwNzWdi6ggwqjuX",                   // bad ip part
	"/ip4/1.2.3.4/tcp/1234/dms3fs",                                    // no id
	"/ip4/1.2.3.4/tcp/1234/dms3fs/",                                   // no id
}

func newMultiaddr(t *testing.T, s string) ma.Multiaddr {
	maddr, err := ma.NewMultiaddr(s)
	if err != nil {
		t.Fatal(err)
	}
	return maddr
}

func TestParseStringGood(t *testing.T) {
	for _, g := range good {
		if _, err := ParseString(g); err != nil {
			t.Error("failed to parse", g, err)
		}
	}
}

func TestParseStringBad(t *testing.T) {
	for _, b := range bad {
		if _, err := ParseString(b); err == nil {
			t.Error("succeeded in parsing", b)
		}
	}
}

func TestParseMultiaddrGood(t *testing.T) {
	for _, g := range good {
		if _, err := ParseMultiaddr(newMultiaddr(t, g)); err != nil {
			t.Error("failed to parse", g, err)
		}
	}
}

func TestParseMultiaddrBad(t *testing.T) {
	for _, b := range bad {
		m, err := ma.NewMultiaddr(b)
		if err != nil {
			continue // skip these.
		}

		if _, err := ParseMultiaddr(m); err == nil {
			t.Error("succeeded in parsing", m)
		}
	}
}

func TestIDMatches(t *testing.T) {
	for _, g := range good {
		a, err := ParseString(g)
		if err != nil {
			t.Error("failed to parse", g, err)
			continue
		}

		sp := strings.Split(g, "/")
		sid := sp[len(sp)-1]
		id, err := peer.IDB58Decode(sid)
		if err != nil {
			t.Error("failed to parse", sid, err)
			continue
		}

		if a.ID() != id {
			t.Error("not equal", a.ID(), id)
		}
	}
}

func TestMultiaddrMatches(t *testing.T) {
	for _, g := range good {
		a, err := ParseString(g)
		if err != nil {
			t.Error("failed to parse", g, err)
			continue
		}

		m := newMultiaddr(t, g)
		if !a.Multiaddr().Equal(m) {
			t.Error("not equal", a.Multiaddr(), m)
		}
	}
}

func TestTransport(t *testing.T) {
	for _, g := range good {
		a, err := ParseString(g)
		if err != nil {
			t.Error("failed to parse", g, err)
			continue
		}

		m := newMultiaddr(t, g)
		split := ma.Split(m)
		if len(split) <= 1 {
			if Transport(a) != nil {
				t.Error("should not have a transport")
			}
		} else {
			m = ma.Join(split[:len(split)-1]...)
			if a.Multiaddr().Equal(m) {
				t.Error("should not be equal", a.Multiaddr(), m)
			}
			if !Transport(a).Equal(m) {
				t.Error("should be equal", Transport(a), m)
			}
		}
	}
}
