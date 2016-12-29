package o3

import (
	"golang.org/x/crypto/nacl/box"
	"crypto/rand"
)


// SessionContext is a passable structure containing all
// established keys and nonces required for communication with
// the server
type SessionContext struct {
	ID ThreemaID
	//TODO it might make more sense in a lot of places to use pointers here
	clientSPK   [32]byte //client short-term public key
	clientSSK   [32]byte //client short-term secret key
	serverSPK   [32]byte //server short-term public key
	serverLPK   [32]byte //server long-term public key
	clientNonce nonce
	serverNonce nonce
}

// NewSessionContext returns a new SessionContext
func NewSessionContext(ID ThreemaID) SessionContext {
	sc := SessionContext{
		serverLPK: [32]byte{69, 11, 151, 87, 53, 39, 159, 222, 203, 51, 19, 100, 143, 95, 198, 238, 159, 244, 54, 14, 169, 42, 140, 23, 81, 198, 97, 228, 192, 216, 201, 9},
		ID:        ID}
		
	// New Session means new ephemeral keys and nonce
	sc.clientNonce = newNonce()

	pk, sk, err := box.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	copy(sc.clientSPK[:], (*pk)[:])
	copy(sc.clientSSK[:], (*sk)[:])

	return sc
}
