package protocol

import (
	"crypto/cipher"
)

type cfb8 struct {
	c                cipher.Block
	blockSize        int
	iv, iv_real, tmp []byte
	de               bool
}

func newCFB8(c cipher.Block, iv []byte, decrypt bool) *cfb8 {
	if len(iv) != 16 {
		panic("bad iv length!")
	}
	cp := make([]byte, 256)
	copy(cp, iv)
	return &cfb8{
		c:         c,
		blockSize: c.BlockSize(),
		iv:        cp[:16],
		iv_real:   cp,
		tmp:       make([]byte, 16),
		de:        decrypt,
	}
}

func newCFB8Decrypt(c cipher.Block, iv []byte) *cfb8 {
	return newCFB8(c, iv, true)
}

func newCFB8Encrypt(c cipher.Block, iv []byte) *cfb8 {
	return newCFB8(c, iv, false)
}

func (cf *cfb8) XORKeyStream(dst, src []byte) {
	for i := 0; i < len(src); i++ {
		val := src[i]
		cf.c.Encrypt(cf.tmp, cf.iv)
		val = val ^ cf.tmp[0]

		if cap(cf.iv) >= 17 {
			cf.iv = cf.iv[1:17]
		} else {
			copy(cf.iv_real, cf.iv[1:])
			cf.iv = cf.iv_real[:16]
		}

		if cf.de {
			cf.iv[15] = src[i]
		} else {
			cf.iv[15] = val
		}
		dst[i] = val
	}
}
