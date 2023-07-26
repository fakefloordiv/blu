package sdp

import (
	"strconv"
	"strings"
)

type NetType string

const (
	IN NetType = "IN"
)

type AddrType string

const (
	IP4 AddrType = "IP4"
	IP6 AddrType = "IP6"
)

type Origin struct {
	Username       string
	SessId         string // TODO: these may be ordinary integers
	SessVersion    string // TODO: these may be ordinary integers
	NetType        NetType
	AddrType       AddrType
	UnicastAddress string
}

func (o Origin) Parse(value string) (Origin, error) {
	sp := strings.IndexByte(value, ' ')
	if sp == -1 {
		return o, ErrBadSyntax
	}

	o.Username, value = value[:sp], value[sp+1:]

	sp = strings.IndexByte(value, ' ')
	if sp == -1 {
		return o, ErrBadSyntax
	}

	o.SessId, value = value[:sp], value[sp+1:]

	sp = strings.IndexByte(value, ' ')
	if sp == -1 {
		return o, ErrBadSyntax
	}

	o.SessVersion, value = value[:sp], value[sp+1:]

	sp = strings.IndexByte(value, ' ')
	if sp == -1 {
		return o, ErrBadSyntax
	}

	switch nettype := value[:sp]; nettype {
	case "IN":
		o.NetType = IN
	default:
		return o, ErrUnknownNetType
	}

	value = value[sp+1:]

	sp = strings.IndexByte(value, ' ')
	if sp == -1 {
		return o, ErrBadSyntax
	}

	switch addrtype := value[:sp]; addrtype {
	case "IP4":
		o.AddrType = IP4
	case "IP6":
		o.AddrType = IP6
	default:
		return o, ErrUnknownAddrType
	}

	value = value[sp+1:]

	if len(value) == 0 {
		return o, ErrBadSyntax
	}

	o.UnicastAddress = value

	return o, nil
}

type ConnectionInfo struct {
	NetType   NetType
	AddrType  AddrType
	Address   string
	TTL       int
	AddrRange int
}

func (c ConnectionInfo) Parse(value string) (info ConnectionInfo, err error) {
	sp := strings.IndexByte(value, ' ')
	if sp == -1 {
		return c, ErrBadSyntax
	}

	switch nettype := value[:sp]; nettype {
	case "IN":
		c.NetType = IN
	default:
		return c, ErrUnknownNetType
	}

	value = value[sp+1:]

	sp = strings.IndexByte(value, ' ')
	if sp == -1 {
		return c, ErrBadSyntax
	}

	switch addrtype := value[:sp]; addrtype {
	case "IP4":
		c.AddrType = IP4
	case "IP6":
		c.AddrType = IP6
	default:
		return c, ErrUnknownAddrType
	}

	value = value[sp+1:]

	if sp = strings.IndexByte(value, ' '); sp != -1 {
		value = value[:sp]
	}

	c.Address, c.TTL, c.AddrRange, err = parseAddress(value)

	return c, err
}

// parseAddress parses the address provided in connection info field. It follows
// the following grammar:
//
//	(ip4addr | ip6addr) [ "/" ttl [ "/" addr-range ] ]
//
// If no TTL is provided, -1 will be returned.
// If no addr-range is provided, 0 will be returned.
//
// Actually this grammar doesn't follow the official one from the RFC. Those one is pretty confusing,
// as TTL there is optional, while address range isn't. But in case there is just one element
// after the address separated with a slash, then it isn't an addr-range, but TTL. In case you want
// to specify the addr-range, then you MUST also specify TTL.
func parseAddress(addr string) (outAddr string, ttl, addrrange int, err error) {
	slash := strings.IndexByte(addr, '/')
	if slash == -1 {
		return addr, -1, 0, nil
	}

	var rawTTL string
	addr, rawTTL = addr[:slash], addr[slash+1:]
	if slash = strings.IndexByte(rawTTL, '/'); slash != -1 {
		addrrange, err = strconv.Atoi(rawTTL[slash+1:])
		if err != nil {
			return "", 0, 0, err
		}

		rawTTL = rawTTL[:slash]
	}

	ttl, err = strconv.Atoi(rawTTL)
	if err != nil {
		return "", 0, 0, err
	}

	return addr, ttl, addrrange, nil
}
