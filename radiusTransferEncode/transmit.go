package radiusTransferEncode

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"time"
)

//整形转换成字节
func Int32ToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt32(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int(x)
}

func EncryptMessage(key []byte, data []byte) ([]byte, error) {
	aesEnc := AesEncrypt{}
	return aesEnc.Encrypt(key, data)
}
func DecryptMessage(key []byte, data []byte) ([]byte, error) {
	aesEnc := AesEncrypt{}
	return aesEnc.Decrypt(key, data)
}

func PackMessage(userName string, userPwd string, src []byte) ([]byte, error) {
	key := []byte(userPwd)
	addhash := GetHash("sha1", src)
	addhash = append(addhash, src...)

	encrypt, err := EncryptMessage(key, addhash)
	if err != nil {
		return nil, err
	}
	buf1 := GetHash("md5", []byte(userName))
	buf1 = append(buf1, encrypt...)
	buf := Int32ToBytes(len(buf1))
	return append(buf, buf1...), nil
}
func UnpackMessage(userName string, userPwd string, data []byte) ([]byte, error) {
	//len := BytesToInt32(data[:4])
	key := []byte(userPwd)
	buf1 := data[4:20]
	buf2 := GetHash("md5", []byte(userName))
	if !bytes.Equal(buf1, buf2) {
		return nil, errors.New("this is not my message")
	}
	decrypt, err := DecryptMessage(key, data[20:])
	if err != nil {
		return nil, err
	}
	hs := GetHash("sha1", decrypt[20:])
	if bytes.Equal(hs, decrypt[:20]) {
		return decrypt[20:], nil
	}
	return nil, errors.New("this content is error, not passing hash check")
}
func ServerUnpackMessage(data []byte, userMap *map[string]*UserInfo) ([]byte, *UserInfo, error) {
	//len := BytesToInt32(data[:4])
	buf1 := data[4:20]
	user, ok := (*userMap)[string(buf1)]

	if !ok {
		return nil, nil, errors.New("not find this user")
	}

	decrypt, err := DecryptMessage([]byte(user.UserPwd), data[20:])
	if err != nil {
		return nil, user, err
	}
	hs := GetHash("sha1", decrypt[20:])
	if bytes.Equal(hs, decrypt[:20]) {
		return decrypt[20:], user, nil
	}
	return nil, user, errors.New("this content is error, not passing hash check")
}
func RecvPacket(conn *net.Conn) ([]byte, error) {
	var buf []byte
	timeoutSec := time.Duration(int64(60) * int64(time.Second))
	(*conn).SetReadDeadline(time.Now().Add(timeoutSec))
	packetLen := 0
	for {
		var tmp []byte
		tmp = make([]byte, 4096)
		ret, err := (*conn).Read(tmp)
		if err != nil {
			return nil, err
		}
		if ret == 0 {
			return buf, errors.New("remote close this socket!")
		}
		if len(buf) < 4 {
			buf = append(buf, tmp[:ret]...)
			if len(buf) >= 4 {
				packetLen = BytesToInt32(buf[:4])
				if len(buf) >= packetLen+4 {
					break
				}
			}
		} else {
			buf = append(buf, tmp[:ret]...)
			if len(buf) >= packetLen+4 {
				break
			}
		}
	}
	return buf, nil
}
