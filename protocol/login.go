package protocol

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
)

var (
	publicKeyBytes []byte
	privateKey     *rsa.PrivateKey
)

func init() {
	var err error
	privateKey, err = rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	privateKey.Precompute()

	publicKeyBytes, err = x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		panic(err)
	}
}

type Authenticator interface {
	Authenticate(username string, serverID string, sharedSecret, publicKey []byte) (uuid string, err error)
}

//Login function authenticates a player who tries to join server
func (conn *Conn) Login(handshake *Handshake, authenticator Authenticator, protocolVersion int) (name string, uuid string, err error) {
	if handshake.ProtocolVersion > VarInt(protocolVersion) {
		return "", "", errors.New("Server out of date!")
	} else if handshake.ProtocolVersion < VarInt(protocolVersion) {
		return "", "", errors.New("Client out of date!")
	}

	conn.State = Login

	packet, err := conn.ReadPacket()
	if err != nil {
		return
	}
	lStart, ok := packet.(*LoginStart)
	if !ok {
		err = fmt.Errorf("Unexpected packet")
		return
	}
	name = lStart.Username

	verifyToken := make([]byte, 16) //Used by the server to check encryption is working correctly
	rand.Read(verifyToken)

	var serverID = "-"
	if authenticator != nil {
		serverBytes := make([]byte, 10)
		rand.Read(serverBytes)
		serverID = hex.EncodeToString(serverBytes)
	} else {
		if len(name) > 16 {
			name = name[:16]
		}
	}

	conn.WritePacket(&EncryptionKeyRequest{
		ServerID:    serverID,
		PublicKey:   publicKeyBytes,
		VerifyToken: verifyToken,
	})

	packet, err = conn.ReadPacket()
	if err != nil {
		return
	}
	encryptionResponse, ok := packet.(*EncryptionKeyResponse)
	if !ok {
		err = fmt.Errorf("Unexpected packet")
		return
	}

	sharedSecret, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptionResponse.SharedSecret)
	if err != nil {
		return
	}

	verifyTokenResponse, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptionResponse.VerifyToken)
	if err != nil {
		return
	}
	if !bytes.Equal(verifyToken, verifyTokenResponse) {
		return
	}

	if authenticator != nil {
		if uuid, err = authenticator.Authenticate(name, serverID, sharedSecret, publicKeyBytes); err != nil {
			return
		}
	} else {
		idBytes := make([]byte, 16)
		rand.Read(idBytes)
		uuid = hex.EncodeToString(idBytes)
	}

	aesCipher, err := aes.NewCipher(sharedSecret)
	if err != nil {
		return
	}

	conn.In = cipher.StreamReader{
		R: conn.In,
		S: newCFB8Decrypt(aesCipher, sharedSecret),
	}
	conn.Out = cipher.StreamWriter{
		W: conn.Out,
		S: newCFB8Encrypt(aesCipher, sharedSecret),
	}
	conn.WritePacket(&LoginSuccess{
		uuid,
		name,
	})
	conn.State = Play

	return
}
