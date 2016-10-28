package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/oikomi/FishChatServer2/codec"
	"github.com/oikomi/FishChatServer2/libnet"
	"github.com/oikomi/FishChatServer2/protocol/external"
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "false")
}

func checkErr(err error) {
	if err != nil {
		glog.Error(err)
	}
}

func clientLoop(session *libnet.Session, protobuf *codec.ProtobufProtocol) {
	var err error
	var clientMsg *libnet.Session
	err = session.Send(&external.ReqAccessServer{
		Cmd: external.ReqAccessServerCMD,
	})
	checkErr(err)
	rsp, err := session.Receive()
	checkErr(err)
	glog.Info(string(rsp))
	if rsp != nil {
		baseCMD := &external.Base{}
		if err = proto.Unmarshal(rsp, baseCMD); err != nil {
			glog.Error(err)
		}
		switch baseCMD.Cmd {
		case external.ReqAccessServerCMD:
			resSelectMsgServerForClientPB := &external.ResSelectAccessServerForClient{}
			if err = proto.Unmarshal(rsp, resSelectMsgServerForClientPB); err != nil {
				glog.Error(err)
			}
			glog.Info(resSelectMsgServerForClientPB)
			glog.Info(resSelectMsgServerForClientPB.Addr)
			clientMsg, err = libnet.Connect("tcp", resSelectMsgServerForClientPB.Addr, protobuf, 0)
			checkErr(err)
		}
	}
	fmt.Print("input my id :")
	var myID int64
	if _, err := fmt.Scanf("%d\n", &myID); err != nil {
		glog.Error(err.Error())
	}
	err = clientMsg.Send(&external.ReqLogin{
		Cmd: external.LoginCMD,
		UID: myID,
	})
	checkErr(err)
	rsp, err = clientMsg.Receive()
	checkErr(err)
	glog.Info(string(rsp))

	go func() {
		for {
			rsp, err := clientMsg.Receive()
			if err != nil {
				glog.Error(err.Error())
			}
			// fmt.Printf("%s\n", rmsg)
			glog.Info(string(rsp))
		}
	}()
	for {
		glog.Info("send p2p msg")
		var targetID int64
		fmt.Print("send the id you want to talk :")
		if _, err = fmt.Scanf("%d\n", &targetID); err != nil {
			glog.Error(err.Error())
		}
		var msg string
		fmt.Print("input msg :")
		if _, err = fmt.Scanf("%s\n", &msg); err != nil {
			glog.Error(err.Error())
		}
		err = clientMsg.Send(&external.ReqSendP2PMsg{
			Cmd:       external.SendP2PMsgCMD,
			SourceUID: myID,
			TargetUID: targetID,
			Msg:       msg,
		})
	}
}

func main() {
	var addr string
	flag.StringVar(&addr, "addr", "127.0.0.1:10000", "echo server address")
	flag.Parse()
	protobuf := codec.Protobuf()
	// session, err := libnet.Connect("tcp", addr, libnet.Packet(2, 1024*1024, 1024, binary.BigEndian, TestCodec{}))
	client, err := libnet.Connect("tcp", addr, protobuf, 0)
	checkErr(err)
	clientLoop(client, protobuf)
}
