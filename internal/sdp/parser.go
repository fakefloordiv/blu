package sdp

import (
	"bytes"
)

type Parser struct{}

func NewParser() Parser {
	return Parser{}
}

func (p *Parser) Parse(data []byte) (desc Description, err error) {
	desc.Session, data, err = p.parseSession(data)
	if err != nil {
		return desc, err
	}

	for len(data) > 0 {
		var media Media
		media, data, err = p.parseMedia(data)
		if err != nil {
			return desc, err
		}

		desc.Media = append(desc.Media, media)
	}
	
	return desc, nil
}

func (p *Parser) parseSession(data []byte) (session Session, rest []byte, err error) {
	for len(data) > 0 {
		if len(data) < 2 {
			return session, nil, ErrIncompleteData
		}

		if data[1] != '=' {
			return session, nil, ErrBadSyntax
		}

		if data[0] == 'm' {
			// media starts here
			return session, data, nil
		}

		key := data[0]
		var value string
		value, data = parseValue(data[2:])

		switch key {
		case 'v':
			session.Protocol = value
		case 'o':
			session.Originator, err = session.Originator.Parse(value)
			if err != nil {
				return session, nil, err
			}
		case 's':
			session.Name = value
		case 'i':
			session.Info = value
		case 'u':
			session.URI = value
		case 'e':
			session.Email = value
		case 'p':
			session.Phone = value
		case 'c':
			connInfo, err := ConnectionInfo{}.Parse(value)
			if err != nil {
				return session, nil, err
			}

			session.ConnectionInfo = append(session.ConnectionInfo, connInfo)
		case 'b':
			bwInfo, err := Bandwidth{}.Parse(value)
			if err != nil {
				return session, nil, err
			}

			session.BandwidthInfo = append(session.BandwidthInfo, bwInfo)
		case 'z':
			session.TimeZoneAdjustments = append(session.TimeZoneAdjustments, value)
		case 'k':
			session.EncryptionKey, err = session.EncryptionKey.Parse(value)
			if err != nil {
				return session, nil, err
			}
		case 'a':
			session.Attributes = append(session.Attributes, Attribute{}.Parse(value))
		default:
			return session, nil, ErrUnrecognizedKey
		}
	}

	return session, data, nil
}

func (p *Parser) parseMedia(data []byte) (media Media, rest []byte, err error) {
	for len(data) > 0 {
		if len(data) < 2 {
			return media, nil, ErrIncompleteData
		}

		if data[1] != '=' {
			return media, nil, ErrBadSyntax
		}

		if data[0] == 'm' && media.Name != "" {
			// the next media block description has begun
			return media, data, nil
		}

		key := data[0]
		var value string
		value, data = parseValue(data[2:])

		switch key {
		case 'm':
			media.Name = value
		case 'i':
			media.Title = value
		case 'c':
			connInfo, err := ConnectionInfo{}.Parse(value)
			if err != nil {
				return media, nil, err
			}

			media.ConnectionInfo = append(media.ConnectionInfo, connInfo)
		case 'b':
			media.BandwidthInfo = append(media.BandwidthInfo, value)
		case 'k':
			media.EncryptionKey = value
		case 'a':
			media.Attributes = append(media.Attributes, value)
		default:
			return media, nil, ErrUnrecognizedKey
		}
	}

	return media, nil, nil
}

func parseValue(data []byte) (value string, rest []byte) {
	rest = data
	lf := bytes.IndexByte(data, '\n')
	if lf >= 0 {
		rest, data = data[lf+1:], data[:lf]

		if data[len(data)-1] == '\r' {
			data = data[:len(data)-1]
		}
	}

	return string(data), rest
}
