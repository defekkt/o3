package o3

import "testing"

func TestHandleMessagePacketUnknownType(t *testing.T) {
	sc := &SessionContext{}
	for _, mt := range []MsgType{CONTACT_SET_PHOTO, CONTACT_DELETE_PHOTO, CONTACT_REQ_PHOTO, POLL_CREATE_MESSAGE, 0x82} {
		payload := append([]byte{byte(mt)}, []byte("payload-with-padding-0123456789")...)
		msg, err := sc.handleMessagePacket(messagePacket{Plaintext: payload})
		if err != nil {
			t.Fatalf("type 0x%02x: unexpected error: %v", mt, err)
		}
		um, ok := msg.(UnknownMessage)
		if !ok {
			t.Fatalf("type 0x%02x: got %T, want UnknownMessage", mt, msg)
		}
		if um.Type != mt || string(um.Payload) != "payload-with-padding-0123456789" {
			t.Fatalf("type 0x%02x: got type 0x%02x payload %q", mt, um.Type, um.Payload)
		}
	}
}

func TestHandleMessagePacketMalformedDoesNotPanic(t *testing.T) {
	sc := &SessionContext{}
	// a group text message too short to contain its group header
	msg, err := sc.handleMessagePacket(messagePacket{Plaintext: []byte{byte(GROUPTEXTMESSAGE), 1, 2, 3}})
	if err == nil || msg != nil {
		t.Fatalf("expected an error for a malformed message, got msg=%v err=%v", msg, err)
	}
}

func TestSendErrorDoesNotBlock(t *testing.T) {
	sc := &SessionContext{ErrorChan: make(chan error, 1)}
	for i := 0; i < 5; i++ {
		sc.sendError(nil)
	}
}
