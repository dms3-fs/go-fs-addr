package dms3fsaddr

import (
	"errors"
	"strings"

	logging "github.com/dms3-fs/go-log"
	circuit "github.com/dms3-p2p/go-p2p-circuit"
	peer "github.com/dms3-p2p/go-p2p-peer"
	ma "github.com/dms3-mft/go-multiaddr"
)

var log = logging.Logger("dms3fsaddr")

// ErrInvalidAddr signals an address is not a valid DMS3FS address.
var ErrInvalidAddr = errors.New("invalid DMS3FS address")

type DMS3FSAddr interface {
	ID() peer.ID
	Multiaddr() ma.Multiaddr
	Transport() ma.Multiaddr
	String() string
	Equal(b interface{}) bool
}

type dms3fsAddr struct {
	ma ma.Multiaddr
	id peer.ID
}

func (a dms3fsAddr) ID() peer.ID {
	return a.id
}

func (a dms3fsAddr) Multiaddr() ma.Multiaddr {
	return a.ma
}

func (a dms3fsAddr) Transport() ma.Multiaddr {
	return Transport(a)
}

func (a dms3fsAddr) String() string {
	return a.ma.String()
}

func (a dms3fsAddr) Equal(b interface{}) bool {
	if ib, ok := b.(DMS3FSAddr); ok {
		return a.Multiaddr().Equal(ib.Multiaddr())
	}
	if mb, ok := b.(ma.Multiaddr); ok {
		return a.Multiaddr().Equal(mb)
	}
	return false
}

// ParseString parses a string representation of an address into an DMS3FSAddr
func ParseString(str string) (a DMS3FSAddr, err error) {
	if str == "" {
		return nil, ErrInvalidAddr
	}

	m, err := ma.NewMultiaddr(str)
	if err != nil {
		return nil, err
	}

	return ParseMultiaddr(m)
}

// ParseMultiaddr parses a multiaddr into an DMS3FSAddr
func ParseMultiaddr(m ma.Multiaddr) (a DMS3FSAddr, err error) {
	// never panic.
	defer func() {
		if r := recover(); r != nil {
			log.Debug("recovered from panic: ", r)
			a = nil
			err = ErrInvalidAddr
		}
	}()

	if m == nil {
		return nil, ErrInvalidAddr
	}

	// make sure it's an DMS3FS addr
	parts := ma.Split(m)
	if len(parts) < 1 {
		return nil, ErrInvalidAddr
	}
	dms3fspart := parts[len(parts)-1] // last part
	if dms3fspart.Protocols()[0].Code != ma.P_DMS3FS {
		return nil, ErrInvalidAddr
	}

	// make sure 'dms3fs id' parses as a peer.ID
	peerIdParts := strings.Split(dms3fspart.String(), "/")
	peerIdStr := peerIdParts[len(peerIdParts)-1]
	id, err := peer.IDB58Decode(peerIdStr)
	if err != nil {
		return nil, err
	}

	return dms3fsAddr{ma: m, id: id}, nil
}

func Transport(iaddr DMS3FSAddr) ma.Multiaddr {
	maddr := iaddr.Multiaddr()

	// /dms3fs/QmId is part of the transport address for p2p-circuit
	// TODO clean up the special case
	// we need a consistent way of composing and consumig multiaddrs
	// so that we don't have to do this
	_, err := maddr.ValueForProtocol(circuit.P_CIRCUIT)
	if err == nil {
		return maddr
	}

	split := ma.Split(maddr)
	if len(split) == 1 {
		return nil
	}
	return ma.Join(split[:len(split)-1]...)
}
