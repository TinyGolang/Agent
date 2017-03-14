package agent

import (
	"fmt"
	"net"

	"io"

	"github.com/TinyGolang/protobuf"
	//	"github.com/name5566/leaf/log"
)

type Agent struct {
	Conn          net.Conn
	Read          func(conn io.Reader) ([]byte, error)       //设置读数据回调
	Write         func(conn io.Writer, args ...[]byte) error //设置写数据回调
	Processor     *protobuf.Processor
	OnUserConnect func(args []interface{})
	OnUserClose   func(args []interface{})
}

/*
func NewAgent(
	OnUserConnect func(args []interface{}),
	OnUserClose func(args []interface{}),
	Processor *protobuf.Processor) *Agent {

	a := &Agent{
		OnUserConnect: OnUserConnect,
		OnUserClose:   OnUserClose,
		Processor:     Processor,
	}
	return a
}
*/
func (a *Agent) Run() {

	if a.Read == nil || a.Write == nil {
		fmt.Println("请设置 Read 和 Write")
		return
	}
	for {
		data, err := a.Read(a.Conn)
		fmt.Println("接收数据")
		if err != nil {
			//			log.Debug("read message: %v", err)
			break
		}

		if a.Processor != nil {
			msg, err := a.Processor.Unmarshal(data)
			if err != nil {
				//				log.Debug("unmarshal message error: %v", err)
				break
			}
			err = a.Processor.Route(msg, a)
			if err != nil {
				//				log.Debug("route message error: %v", err)
				break
			}
		}
	}
}
func (a *Agent) WriteMsg(msg interface{}) {
	if a.Processor != nil {
		data, err := a.Processor.Marshal(msg)
		if err != nil {
			//			log.Error("marshal message %v error: %v", reflect.TypeOf(msg), err)
			return
		}
		err = a.Write(a.Conn, data...)
		if err != nil {
			//			log.Error("write message %v error: %v", reflect.TypeOf(msg), err)
		}
	}
}
func (a *Agent) OnConnect() {
	fmt.Println("用户连接")

	if a.OnUserConnect != nil {

		a.OnUserConnect([]interface{}{a})
	}
}

func (a *Agent) OnClose() {
	fmt.Println("用户失连")

	if a.OnUserClose != nil {

		a.OnUserClose([]interface{}{a})
	}
}
