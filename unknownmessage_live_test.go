package o3

import (
	"bytes"
	"testing"
	"time"
)

// TestUnknownTypesLive sends message types the library cannot decode from Alice
// to Bob over the real Threema servers and checks that Bob's receive loop
// survives and still delivers a normal text message afterwards.
func TestUnknownTypesLive(t *testing.T) {
	aliceCtx := initSession(t, "test/idAlice", "test/idAlice.ab", "ThisIsAlice1")
	bobCtx := initSession(t, "test/idBob", "test/idBob.ab", "ThisIsBob234")

	aliceSend, _, err := aliceCtx.Run()
	if err != nil {
		t.Fatal(err)
	}
	_, bobRecv, err := bobCtx.Run()
	if err != nil {
		t.Fatal(err)
	}

	newRaw := func(mt MsgType, payload []byte) UnknownMessage {
		return UnknownMessage{
			messageHeader: messageHeader{
				sender:    aliceCtx.ID.ID,
				recipient: bobCtx.ID.ID,
				id:        NewMsgID(),
				time:      time.Now(),
				pubNick:   aliceCtx.ID.Nick,
			},
			Type:    mt,
			Payload: payload,
		}
	}
	padding := bytes.Repeat([]byte{32}, 32) // valid PKCS#7-style padding, 32 bytes

	aliceSend <- newRaw(CONTACT_SET_PHOTO, padding)
	aliceSend <- newRaw(GROUPTEXTMESSAGE, []byte{1, 2, 3}) // too short for its group header
	if _, err := aliceCtx.SendTextMessage(bobCtx.ID.String(), "still alive", aliceSend); err != nil {
		t.Fatal(err)
	}

	var gotPhoto, gotMalformed, gotText bool
	timeout := time.After(20 * time.Second)
	for !(gotPhoto && gotMalformed && gotText) {
		select {
		case r := <-bobRecv:
			switch m := r.Msg.(type) {
			case UnknownMessage:
				t.Logf("Bob got UnknownMessage type 0x%02x (%d payload bytes)", m.Type, len(m.Payload))
				if m.Type == CONTACT_SET_PHOTO && bytes.Equal(m.Payload, padding) {
					gotPhoto = true
				}
			case TextMessage:
				t.Logf("Bob got text %q", m.Text())
				gotText = m.Text() == "still alive"
			default:
				if r.Err != nil {
					t.Logf("Bob got error (expected for the malformed one): %v", r.Err)
					gotMalformed = true
				}
			}
		case <-timeout:
			t.Fatalf("timeout: photo=%v malformed=%v text=%v", gotPhoto, gotMalformed, gotText)
		}
	}
}
