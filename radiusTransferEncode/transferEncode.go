package radiusTransferEncode

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"

	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"hash"
)

func GetHash(style string, inputSrc []byte) []byte {
	var hs hash.Hash
	switch style {
	case "md5":
		hs = md5.New()
	case "sha1":
		hs = sha1.New()
	case "sha512":
		hs = sha512.New()
	default:
		hs = sha256.New()
	}
	hs.Write(inputSrc)
	return hs.Sum(nil)
}
func test() {
	aesEnc := new(AesEncrypt)
	arrEncrypt, err := aesEnc.Encrypt([]byte("123456"), []byte("abcde"))
	if err != nil {
		fmt.Println(arrEncrypt)
		return
	}
	strMsg, err := aesEnc.Decrypt([]byte("123456"), arrEncrypt)
	if err != nil {
		fmt.Println(arrEncrypt)
		return
	}
	fmt.Println(strMsg)
}

type AesEncrypt struct {
}

var (
	ErrAESTextSize = errors.New("ciphertext is not a multiple of the block size")
	ErrAESPadding  = errors.New("cipher padding size error")
)

//加密字符串
func (this *AesEncrypt) Encrypt(strKey []byte, src []byte) ([]byte, error) {
	key := GetHash("sha256", strKey)
	iv := GetHash("md5", strKey)
	padLen := aes.BlockSize - (len(src) % aes.BlockSize)
	for i := 0; i < padLen; i++ {
		src = append(src, byte(padLen))
	}
	encrypted := make([]byte, len(src))
	aesBlockEncrypter, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesEncrypter := cipher.NewCBCEncrypter(aesBlockEncrypter, iv)
	aesEncrypter.CryptBlocks(encrypted, src)
	return encrypted, nil
}

//解密字符串
func (this *AesEncrypt) Decrypt(strKey []byte, src []byte) ([]byte, error) {
	if len(src) < aes.BlockSize || len(src)%aes.BlockSize != 0 {
		return nil, ErrAESTextSize
	}
	key := GetHash("sha256", strKey)
	iv := GetHash("md5", strKey)
	decrypted := make([]byte, len(src))
	aesBlockDecrypter, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesDecrypter := cipher.NewCBCDecrypter(aesBlockDecrypter, iv)
	aesDecrypter.CryptBlocks(decrypted, src)
	paddingLen := int(decrypted[len(decrypted)-1])
	if paddingLen > 16 {
		fmt.Printf("paddingLen is %d, key is %s\n", paddingLen, string(strKey))
		//BinaryPrint(src)
		//BinaryPrint(decrypted)
		return nil, ErrAESPadding
	}
	for i := 2; i <= paddingLen; i++ {
		if paddingLen != int(decrypted[len(decrypted)-i]) {
			//BinaryPrint(decrypted)
			return nil, ErrAESPadding
		}
	}
	return decrypted[:len(decrypted)-paddingLen], nil
}

// UserInfo ...
type UserInfo struct {
	UserName string `json:"UserName"`
	UserPwd  string `json:"UserPwd"`
}

// RadiusUserInfo ...
type RadiusUserInfo struct {
	UserName   string `json:"UserName"`
	UserPwd    string `json:"UserPwd"`
	ExpireTime string `json:"ExpireTime"`
}
type AddRadiusUser struct {
	ChildUsers []RadiusUserInfo `json:"ChildUsers"`
}
type DelRadiusUser struct {
	ChildUsers []string `json:"ChildUsers"`
}

const (
	base64Table = "123QRSTUabcdVWXYZHijKLAWDCABDstEFGuvwxyzGHIJklmnopqr234560178912"
)

func GzipCompress(input []byte) ([]byte, error) {
	var res bytes.Buffer
	w := gzip.NewWriter(&res)

	_, err := w.Write(input)
	if err != nil {
		fmt.Println("gzip compress failed, ", err)
		return []byte{}, err
	}
	w.Close()
	return res.Bytes(), nil

}
func GzipUnCompress(input []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(input))
	if err != nil {
		fmt.Println("create gzip reader failed, ", err)
		return []byte{}, err
	}
	defer r.Close()
	undatas, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Println("gzip decompress failed, ", err)
		return []byte{}, err
	}

	return undatas, nil
}
func CreateAddRadiusNamesRequest(userName string, userPwd string, names *AddRadiusUser) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := WriteMessageHeader(buf, &RequestAddRadiusNames)
	if err != nil {
		return []byte{}, err
	}
	err = WriteAddRadiusNamesRequest(buf, names)
	if err != nil {
		return []byte{}, err
	}
	return PackMessage(userName, userPwd, buf.Bytes()[:])
}

func CreateDelRadiusNamesRequest(userName string, userPwd string, names *DelRadiusUser) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := WriteMessageHeader(buf, &RequestDelRadiusNames)
	if err != nil {
		return []byte{}, err
	}
	err = WriteDelRadiusNamesRequest(buf, names)
	if err != nil {
		return []byte{}, err
	}
	return PackMessage(userName, userPwd, buf.Bytes()[:])
}

func CreateAnswer(typeMajor int32, typeMinor int32, result int32, msg string) ([]byte, error) {
	var header TransmitHeader = TransmitHeader{typeMajor, typeMinor}
	var body RadiusResult = RadiusResult{result, msg}
	buf := new(bytes.Buffer)
	err := WriteMessageHeader(buf, &header)
	if err != nil {
		return []byte{}, err
	}
	err = WriteAnswerBody(buf, &body)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes()[:], nil
}
