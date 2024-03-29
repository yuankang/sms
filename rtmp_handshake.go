package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
)

var (
	FPKey = []byte{
		// Genuine Adobe Flash Player 001
		'G', 'e', 'n', 'u', 'i', 'n', 'e', ' ', 'A', 'd',
		'o', 'b', 'e', ' ', 'F', 'l', 'a', 's', 'h', ' ',
		'P', 'l', 'a', 'y', 'e', 'r', ' ', '0', '0', '1',
		0xF0, 0xEE, 0xC2, 0x4A, 0x80, 0x68, 0xBE, 0xE8, 0x2E, 0x00,
		0xD0, 0xD1, 0x02, 0x9E, 0x7E, 0x57, 0x6E, 0xEC, 0x5D, 0x2D,
		0x29, 0x80, 0x6F, 0xAB, 0x93, 0xB8, 0xE6, 0x36, 0xCF, 0xEB,
		0x31, 0xAE,
	}
	FMSKey = []byte{
		// Genuine Adobe Flash Media Server 001
		'G', 'e', 'n', 'u', 'i', 'n', 'e', ' ', 'A', 'd',
		'o', 'b', 'e', ' ', 'F', 'l', 'a', 's', 'h', ' ',
		'M', 'e', 'd', 'i', 'a', ' ', 'S', 'e', 'r', 'v',
		'e', 'r', ' ', '0', '0', '1',
		0xF0, 0xEE, 0xC2, 0x4A, 0x80, 0x68, 0xBE, 0xE8, 0x2E, 0x00,
		0xD0, 0xD1, 0x02, 0x9E, 0x7E, 0x57, 0x6E, 0xEC, 0x5D, 0x2D,
		0x29, 0x80, 0x6F, 0xAB, 0x93, 0xB8, 0xE6, 0x36, 0xCF, 0xEB,
		0x31, 0xAE,
	}
	FPKeyP  = FPKey[:30]
	FMSKeyP = FMSKey[:36]
)

/*************************************************/
/* RtmpHandshakeServer
/*************************************************/
// C1(1536): time(4) + zero(4) + randomByte(1528)
// randomByte(1528): key(764) + digest(764) or digest(764) + key(764)
// key(764): randomData(n) + keyData(128) + randomData(764-n-128-4) + offset(4)
// n = offset的值
// digest(764): offset(4) + randomData(n) + digestData(32) + randomData(764-4-n-32)
// n = offset的值
// 简单握手时c2(1536): time(4) + time(4) + randomEcho(1528)
// 复杂握手时c2(1536): randomData(1504) + digestData(32)
func CreateComplexS1S2(s *Stream, C1, S1, S2 []byte) error {
	// 发起rtmp连接的一方key用FPKeyP，被连接的一方key用FMSKeyP
	// 1 重新计算C1的digest和C1中的digest比较,一样才行
	// 2 计算S1的digest
	// 3 计算S2的key和digest
	var c1Digest []byte
	if c1Digest = DigestFind(s, C1, 8); c1Digest == nil {
		c1Digest = DigestFind(s, C1, 772)
	}
	if c1Digest == nil {
		err := fmt.Errorf("can't find digest in C1")
		s.log.Println(err)
		return err
	}

	cTime := ByteToUint32(C1[0:4], BE)
	cZero := ByteToUint32(C1[4:8], BE)
	sTime := cTime
	sZero := cZero
	//sZero := uint32(0x0d0e0a0d)

	Uint32ToByte(sTime, S1[0:4], BE)
	Uint32ToByte(sZero, S1[4:8], BE)
	rand.Read(S1[8:])
	pos := DigestFindPos(S1, 8)
	s1Digest := DigestCreate(S1, pos, FMSKeyP)
	//s.log.Println("s1Digest create:", s1Digest)
	copy(S1[pos:], s1Digest)

	//???
	s2DigestKey := DigestCreate(c1Digest, -1, FMSKey)
	//s.log.Println("s2DigestKey create:", s2DigestKey)

	rand.Read(S2)
	pos = len(S2) - 32
	s2Digest := DigestCreate(S2, pos, s2DigestKey)
	//s.log.Println("s2Digest create:", s2Digest)
	copy(S2[pos:], s2Digest)
	return nil
}

func DigestFindPos(C1 []byte, base int) (pos int) {
	pos, offset := 0, 0
	for i := 0; i < 4; i++ {
		offset += int(C1[base+i])
	}
	//764 - 4 - 32 = 728, offset 最大值不能超过728
	pos = base + 4 + (offset % 728)
	return
}

func DigestCreate(b []byte, pos int, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	if pos <= 0 {
		h.Write(b)
	} else {
		h.Write(b[:pos])
		h.Write(b[pos+32:])
	}
	return h.Sum(nil)
}

func DigestFind(s *Stream, C1 []byte, base int) []byte {
	pos := DigestFindPos(C1, base)
	c1Digest := DigestCreate(C1, pos, FPKeyP)
	//s.log.Println("c1Digest in C1:", C1[pos:pos+32])
	//s.log.Println("c1Digest create:", c1Digest)
	if hmac.Equal(C1[pos:pos+32], c1Digest) {
		return c1Digest
	}
	return nil
}

func RtmpHandshakeServer(s *Stream) error {
	var err error
	var C0C1C2S0S1S2 [(1 + 1536*2) * 2]byte

	C0C1C2 := C0C1C2S0S1S2[:1536*2+1]
	C0 := C0C1C2[:1]
	C1 := C0C1C2[1 : 1536+1]
	C0C1 := C0C1C2[:1536+1]
	C2 := C0C1C2[1536+1:]

	S0S1S2 := C0C1C2S0S1S2[1536*2+1:]
	S0 := S0S1S2[:1]
	S1 := S0S1S2[1 : 1536+1]
	S2 := S0S1S2[1536+1:]

	if _, err = io.ReadFull(s.Conn, C0C1); err != nil {
		s.log.Println(err)
		return err
	}

	//rtmp协议版本号, 0x03明文, 0x06密文
	if C0[0] != 3 {
		err = fmt.Errorf("invalid rtmp client version %d", C0[0])
		s.log.Println(err)
		return err
	}
	S0[0] = 3

	cZero := ByteToUint32(C1[4:8], BE)
	//s.log.Printf("cZero: 0x%x", cZero)
	if cZero == 0 {
		s.log.Println("rtmp simple handshake")
		copy(S1, C2)
		copy(S2, C1)
	} else {
		s.log.Println("rtmp complex handshake")
		err = CreateComplexS1S2(s, C1, S1, S2)
		if err != nil {
			s.log.Println(err)
			return err
		}
	}

	if _, err = s.Conn.Write(S0S1S2); err != nil {
		s.log.Println(err)
		return err
	}
	s.Conn.Flush()

	if _, err = io.ReadFull(s.Conn, C2); err != nil {
		s.log.Println(err)
		return err
	}
	return nil
}

/*************************************************/
/* RtmpHandshakeClient
/*************************************************/
func RtmpHandshakeClient(s *Stream) error {
	var C0C1C2S0S1S2 [(1 + 1536*2) * 2]byte

	C0C1C2 := C0C1C2S0S1S2[:1536*2+1]
	C0 := C0C1C2[:1]
	C0C1 := C0C1C2[:1536+1]
	C2 := C0C1C2[1536+1:]

	S0S1S2 := C0C1C2S0S1S2[1536*2+1:]

	C0[0] = 3
	_, err := s.Conn.Write(C0C1)
	if err != nil {
		s.log.Println(err)
		return err
	}
	s.Conn.Flush()

	_, err = io.ReadFull(s.Conn, S0S1S2)
	if err != nil {
		s.log.Println(err)
		return err
	}

	S1 := S0S1S2[1 : 1536+1]
	ver := ByteToUint32(S1[4:8], BE)
	if ver != 0 {
		C2 = S1
	} else {
		C2 = S1
	}

	_, err = s.Conn.Write(C2)
	if err != nil {
		s.log.Println(err)
		return err
	}
	s.Conn.Flush()
	return nil
}
