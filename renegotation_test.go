package webrtc

import (
	"context"
	"testing"
	"time"

	"github.com/pion/transport/test"
	"github.com/pion/webrtc/v3/pkg/media"
	"github.com/stretchr/testify/assert"
)

func Test_Renegotation_E2E(t *testing.T) {
	lim := test.TimeOut(time.Second * 30)
	defer lim.Stop()

	report := test.CheckRoutines(t)
	defer report()

	pcOffer, pcAnswer, err := newPair()
	if err != nil {
		t.Fatal(err)
	}

	haveRenegotiated := &atomicBool{}
	onTrackFired, onTrackFiredFunc := context.WithCancel(context.Background())
	pcAnswer.OnTrack(func(track *TrackRemote, r *RTPReceiver) {
		if !haveRenegotiated.get() {
			t.Fatal("OnTrack was called before renegotiation")
		}
		onTrackFiredFunc()
	})

	assert.NoError(t, signalPair(pcOffer, pcAnswer))

	_, err = pcAnswer.AddTransceiverFromKind(RTPCodecTypeVideo, RTPTransceiverInit{Direction: RTPTransceiverDirectionRecvonly})
	assert.NoError(t, err)

	vp8Track, err := NewTrackLocalStaticSample(RTPCodecCapability{MimeType: MimeTypeVP8}, "foo", "bar")
	assert.NoError(t, err)

	sender, err := pcOffer.AddTrack(vp8Track)
	assert.NoError(t, err)

	// Send 10 packets, OnTrack MUST not be fired
	for i := 0; i <= 10; i++ {
		assert.NoError(t, vp8Track.WriteSample(media.Sample{Data: []byte{0x00}, Duration: time.Second}))
		time.Sleep(20 * time.Millisecond)
	}

	haveRenegotiated.set(true)
	assert.False(t, sender.isNegotiated())
	offer, err := pcOffer.CreateOffer(nil)
	assert.True(t, sender.isNegotiated())
	assert.NoError(t, err)
	assert.NoError(t, pcOffer.SetLocalDescription(offer))
	assert.NoError(t, pcAnswer.SetRemoteDescription(offer))
	answer, err := pcAnswer.CreateAnswer(nil)
	assert.NoError(t, err)

	assert.NoError(t, pcAnswer.SetLocalDescription(answer))

	pcOffer.ops.Done()
	assert.Equal(t, 0, len(vp8Track.rtpTrack.bindings))

	assert.NoError(t, pcOffer.SetRemoteDescription(answer))

	pcOffer.ops.Done()
	assert.Equal(t, 1, len(vp8Track.rtpTrack.bindings))

	sendVideoUntilDone(onTrackFired.Done(), t, []*TrackLocalStaticSample{vp8Track})

	closePairNow(t, pcOffer, pcAnswer)
}
