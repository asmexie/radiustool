package radiusTransferEncode

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"time"
)

func SendAndRecv(conn *net.Conn, send []byte) (recv []byte, err error) {
	if len(send) == 0 {
		err = errors.New("要发送的数据包为空")
		return
	}
	ret, err := (*conn).Write(send)
	if err != nil {
		return
	}
	if ret != len(send) {
		err = errors.New("send is not complete")
		return
	}
	recv, err = RecvPacket(conn)

	return
}

func AddRadiusUserToServer(info *AddRadiusUser, serverAddr string, user *UserInfo) (int, error) {

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	send, err := CreateAddRadiusNamesRequest(user.UserName, user.UserPwd, info)
	if err != nil {
		return 0, err
	}
	recv, err := SendAndRecv(&conn, send[:])
	if err != nil {
		return 0, err
	}
	unpacket, err := UnpackMessage(user.UserName, user.UserPwd, recv)
	if err != nil {
		return 0, err
	}
	r := bytes.NewReader(unpacket)

	header, err := ReadMessageHeader(r)
	if err != nil {
		return 0, err
	}
	if header.RequestType != TypeAddRadiusNames || header.RequestStyle != StyleAnswer {
		return 0, fmt.Errorf("回复数据包类型不符:%v:%v", header.RequestType, header.RequestStyle)
	}
	body, err := ReadAnswerBody(r)
	if err != nil {
		return 0, err
	}
	return int(body.Result), nil
}

func AddRadiusUserToServerAndWait(info *AddRadiusUser, serverAddr string, user *UserInfo, waitSecond int) (int, error) {
	ch := make(chan bool)
	ret := 0
	var err error
	go func() {
		ret, err = AddRadiusUserToServer(info, serverAddr, user)
		ch <- true
	}()

	select {
	case <-ch:
		return ret, err
	case <-time.After(time.Duration(waitSecond) * time.Second):
		return 0, errors.New("timeout")
	}
}
func DelRadiusUserToServer(info *DelRadiusUser, serverAddr string, user *UserInfo) (int, error) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	send, err := CreateDelRadiusNamesRequest(user.UserName, user.UserPwd, info)
	if err != nil {
		return 0, err
	}
	recv, err := SendAndRecv(&conn, send[:])
	if err != nil {
		return 0, err
	}
	unpacket, err := UnpackMessage(user.UserName, user.UserPwd, recv)
	if err != nil {
		return 0, err
	}
	r := bytes.NewReader(unpacket)
	if err != nil {
		return 0, err
	}
	header, err := ReadMessageHeader(r)
	if err != nil {
		return 0, err
	}
	if header.RequestType != TypeDelRadiusNames || header.RequestStyle != StyleAnswer {
		return 0, errors.New("回复数据包类型不符")
	}
	body, err := ReadAnswerBody(r)
	if err != nil {
		return 0, err
	}
	return int(body.Result), nil
}
func DelRadiusUserToServerAndWait(info *DelRadiusUser, serverAddr string, user *UserInfo, waitSecond int) (int, error) {
	ch := make(chan bool)
	ret := 0
	var err error
	go func() {
		ret, err = DelRadiusUserToServer(info, serverAddr, user)
		ch <- true
	}()

	select {
	case <-ch:
		return ret, err
	case <-time.After(time.Duration(waitSecond) * time.Second):
		return 0, errors.New("timeout")
	}
}
