// +build !js

package webrtc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateDataChannelID(t *testing.T) {
	sctpTransportWithChannels := func(ids []uint16) *SCTPTransport {
		ret := &SCTPTransport{dataChannels: []*DataChannel{}}

		for i := range ids {
			id := ids[i]
			ret.dataChannels = append(ret.dataChannels, &DataChannel{id: &id})
		}

		return ret
	}

	t.Run("OK", func(t *testing.T) {
		testCases := []struct {
			role   DTLSRole
			s      *SCTPTransport
			result uint16
		}{
			{DTLSRoleClient, sctpTransportWithChannels([]uint16{}), 0},
			{DTLSRoleClient, sctpTransportWithChannels([]uint16{1}), 0},
			{DTLSRoleClient, sctpTransportWithChannels([]uint16{0}), 2},
			{DTLSRoleClient, sctpTransportWithChannels([]uint16{0, 2}), 4},
			{DTLSRoleClient, sctpTransportWithChannels([]uint16{0, 4}), 2},
			{DTLSRoleServer, sctpTransportWithChannels([]uint16{}), 1},
			{DTLSRoleServer, sctpTransportWithChannels([]uint16{0}), 1},
			{DTLSRoleServer, sctpTransportWithChannels([]uint16{1}), 3},
			{DTLSRoleServer, sctpTransportWithChannels([]uint16{1, 3}), 5},
			{DTLSRoleServer, sctpTransportWithChannels([]uint16{1, 5}), 3},
		}
		for i, testCase := range testCases {
			msg := fmt.Sprintf("case %v", i)
			idPtr := new(uint16)
			err := testCase.s.generateAndSetDataChannelID(testCase.role, &idPtr)
			require.NoError(t, err, msg)
			require.NotNil(t, idPtr, msg)
			assert.Equal(t, testCase.result, *idPtr, msg)
		}
	})

	t.Run("IgnoresClosed", func(t *testing.T) {
		s := sctpTransportWithChannels([]uint16{0, 1})
		for _, dc := range s.dataChannels {
			dc.setReadyState(DataChannelStateClosed)
		}

		// Server
		idPtr := new(uint16)
		err := s.generateAndSetDataChannelID(DTLSRoleServer, &idPtr)
		require.NoError(t, err)
		require.NotNil(t, idPtr)
		assert.EqualValues(t, 1, *idPtr)

		// Client
		err = s.generateAndSetDataChannelID(DTLSRoleClient, &idPtr)
		require.NoError(t, err)
		require.NotNil(t, idPtr)
		assert.EqualValues(t, 0, *idPtr)
	})
}
