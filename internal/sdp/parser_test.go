package sdp

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParser(t *testing.T) {
	t.Run("rfc sample", func(t *testing.T) {
		rfcSample := []byte("v=0\r\n" +
			"o=jdoe 2890844526 2890842807 IN IP4 10.47.16.5\r\n" +
			"s=SDP Seminar\r\n" +
			"i=A Seminar on the session description protocol\r\n" +
			"u=http://www.example.com/seminars/sdp.pdf\r\n" +
			"e=j.doe@example.com (Jane Doe)\r\n" +
			"c=IN IP4 224.2.17.12/127\r\n" +
			//"t=2873397496 2873404696\r\n" +
			"a=recvonly\r\n" +
			"m=audio 49170 RTP/AVP 0\r\n" +
			"m=video 51372 RTP/AVP 99\r\n" +
			"a=rtpmap:99 h263-1998/90000\r\n")
		parser := NewParser()
		desc, err := parser.Parse(rfcSample)
		require.NoError(t, err)
		require.Equal(t, "0", desc.Session.Protocol)
		require.Equal(t, "jdoe 2890844526 2890842807 IN IP4 10.47.16.5", desc.Session.Originator)
		require.Equal(t, "SDP Seminar", desc.Session.Name)
		require.Equal(t, "A Seminar on the session description protocol", desc.Session.Info)
		require.Equal(t, "http://www.example.com/seminars/sdp.pdf", desc.Session.URI)
		require.Equal(t, "j.doe@example.com (Jane Doe)", desc.Session.Email)
		require.Equal(t, "IN IP4 224.2.17.12/127", desc.Session.ConnectionData)
		require.Equal(t, []string{"recvonly"}, desc.Session.Attributes)
		require.Equal(t, 2, len(desc.Media), "must be exactly 2 media blocks")
		media := desc.Media[0]
		require.Equal(t, "audio 49170 RTP/AVP 0", media.Name)
		media = desc.Media[1]
		require.Equal(t, "video 51372 RTP/AVP 99", media.Name)
		require.Equal(t, []string{"rtpmap:99 h263-1998/90000"}, media.Attributes)
	})
}
