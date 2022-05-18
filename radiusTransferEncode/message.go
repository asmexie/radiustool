package radiusTransferEncode

import (
	"bytes"
	"encoding/binary"
)

const (
	TypeUnknow = iota
	TypeAddRadiusNames
	TypeDelRadiusNames
)
const (
	StyleRequest = iota
	StyleAnswer
	StyleNotFindUser
	StyleDecryptFailed
	StylePacketError
	StyleUnknowError
	StyleNetworkRecvError
	StyleTypeNotSurport
	StyleNeedRequest
)

type TransmitHeader struct {
	RequestType  int32
	RequestStyle int32
}

type RadiusResult struct {
	Result int32
	Msg    string
}

var (
	RequestAddRadiusNames TransmitHeader = TransmitHeader{
		TypeAddRadiusNames, StyleRequest}

	RequestDelRadiusNames TransmitHeader = TransmitHeader{
		TypeDelRadiusNames, StyleRequest}

	AnswerAddRadiusNames TransmitHeader = TransmitHeader{
		TypeAddRadiusNames, StyleAnswer}

	AnswerDelRadiusNames TransmitHeader = TransmitHeader{
		TypeDelRadiusNames, StyleAnswer}
)

func WriteMessageHeader(buf *bytes.Buffer, header *TransmitHeader) error {
	err := binary.Write(buf, binary.BigEndian, header.RequestType)
	if err != nil {
		return err
	}
	err = binary.Write(buf, binary.BigEndian, header.RequestStyle)
	return err
}
func ReadMessageHeader(buf *bytes.Reader) (*TransmitHeader, error) {
	var header TransmitHeader
	err := binary.Read(buf, binary.BigEndian, &header.RequestType)
	if err != nil {
		return nil, err
	}

	err = binary.Read(buf, binary.BigEndian, &header.RequestStyle)
	if err != nil {
		return nil, err
	}
	return &header, nil
}
func MyWriteString(b *bytes.Buffer, s string) error {
	data := []byte(s)
	err := binary.Write(b, binary.BigEndian, int32(len(data)))
	if err != nil {
		return err
	}
	if len(data) > 0 {
		_, err = b.Write(data)
	}
	return err
}
func MyReadString(b *bytes.Reader, s *string) error {
	var sLen int32
	err := binary.Read(b, binary.BigEndian, &sLen)
	if err != nil {
		return err
	}
	if sLen <= 0 {
		*s = ""
		return nil
	}
	tmp := make([]byte, sLen)
	n, err := b.Read(tmp)
	if err != nil {
		return err
	}
	*s = string(tmp[:n])
	return nil
}
func WriteAddRadiusNamesRequest(buf *bytes.Buffer, names *AddRadiusUser) error {
	err := binary.Write(buf, binary.BigEndian, int32(len(names.ChildUsers)))
	if err != nil {
		return err
	}
	for _, r := range names.ChildUsers {
		err = MyWriteString(buf, r.UserName)
		if err != nil {
			return err
		}
		err = MyWriteString(buf, r.UserPwd)
		if err != nil {
			return err
		}
		err = MyWriteString(buf, r.ExpireTime)
		if err != nil {
			return err
		}
	}

	return nil

}
func ReadAddRadiusNamesRequest(r *bytes.Reader) (*AddRadiusUser, error) {
	var sLen int32
	err := binary.Read(r, binary.BigEndian, &sLen)
	if err != nil {
		return nil, err
	}
	tmp := &AddRadiusUser{make([]RadiusUserInfo, sLen)}
	for i, _ := range tmp.ChildUsers {
		err = MyReadString(r, &tmp.ChildUsers[i].UserName)
		if err != nil {
			return nil, err
		}
		err = MyReadString(r, &tmp.ChildUsers[i].UserPwd)
		if err != nil {
			return nil, err
		}
		err = MyReadString(r, &tmp.ChildUsers[i].ExpireTime)
		if err != nil {
			return nil, err
		}
	}
	return tmp, nil
}
func WriteDelRadiusNamesRequest(buf *bytes.Buffer, names *DelRadiusUser) error {

	//buf := new(bytes.Buffer)
	//err := WriteMessageHeader(buf, &RequestPutTableData)
	//if err != nil {
	//	return nil, err
	//}
	err := binary.Write(buf, binary.BigEndian, int32(len(names.ChildUsers)))
	if err != nil {
		return err
	}
	for _, u := range names.ChildUsers {
		err = MyWriteString(buf, u)
		if err != nil {
			return err
		}
	}

	return nil

}
func ReadDelRadiusNamesRequest(r *bytes.Reader) (*DelRadiusUser, error) {
	var i, sLen int32
	err := binary.Read(r, binary.BigEndian, &sLen)
	if err != nil {
		return nil, err
	}
	tmp := &DelRadiusUser{make([]string, sLen)}
	for i = 0; i < sLen; i++ {
		err = MyReadString(r, &tmp.ChildUsers[i])
		if err != nil {
			return nil, err
		}

	}
	return tmp, nil
}
func WriteAnswerBody(buf *bytes.Buffer, result *RadiusResult) error {
	err := binary.Write(buf, binary.BigEndian, int32(result.Result))
	if err != nil {
		return err
	}
	err = MyWriteString(buf, result.Msg)
	return err
}
func ReadAnswerBody(r *bytes.Reader) (*RadiusResult, error) {
	body := RadiusResult{}
	err := binary.Read(r, binary.BigEndian, &body.Result)
	if err != nil {
		return nil, err
	}
	err = MyReadString(r, &body.Msg)
	if err != nil {
		return nil, err
	}
	return &body, nil
}
