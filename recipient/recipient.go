package recipient

import (
	//"github.com/actionpay/postmanq/logger"
	"bitbucket.org/asolomonov/postmanq/logger"
	"fmt"
	"net"
	"net/textproto"
)

type Recipient struct {
	id    int
	state State
	conn  net.Conn
}

func newRecipient(id int, events chan *Event) {
	//quit := new(QuitState)
	//
	//commonPossibles := []State{
	//	quit,
	//}
	//
	//input := new(InputState)
	//
	//data := new(DataState)
	//data.SetNext(input)
	//
	rcpt := new(RcptState)
	//rcpt.SetNext(data)
	//
	mail := new(MailState)
	mail.SetNext(rcpt)
	mail.SetPossibles([]State{})
	//input.SetNext(mail)

	ehlo := new(EhloState)
	ehlo.SetNext(mail)
	ehlo.SetPossibles([]State{})

	conn := new(ConnectState)
	conn.SetNext(ehlo)
	//conn.SetPossibles([]State{})

	recipient := &Recipient{
		id:    id,
		state: conn,
	}
	for event := range events {
		recipient.handle(event)
	}
}

func (r *Recipient) handle(event *Event) {
	var id uint
	var buf []byte
	txt := textproto.NewConn(event.conn)
	status := ReadStatus

	for {
		//goto handleStatus

		//handleStatus:
		if r.state == nil {
			continue
		}

		switch status {
		case ReadStatus:
			r.state.SetEvent(event)
			id = txt.Next()
			txt.StartRequest(id)
			buf = r.state.Read(txt)
			txt.EndRequest(id)
			status = r.state.Process(buf)
			logger.By(event.serverHostname).Debug("-> %s", string(buf))
			fmt.Println("->", string(buf))
			fmt.Println("status: ", status)

		case WriteStatus:
			txt.StartResponse(id)
			r.state.Write(txt)
			txt.EndResponse(id)

			r.state = r.state.GetNext()
			status = ReadStatus

		case PossibleStatus:
			fmt.Println("PossibleStatus")
		}

		//goto handleStatus

		//handleStatus:
		//	switch status {
		//	case SuccessStatus:
		//		r.state.Write(txt)
		//		state := r.state.GetNext()
		//		state.SetId(r.state.GetId())
		//		r.state = state
		//
		//	case QuitStatus:
		//		r.state.Write(txt)
		//		event.conn.Close()
		//		return
		//
		//	case FailureStatus:
		//		txt.Cmd("500 Syntax error, command unrecognized")
		//		return
		//
		//	case PossibleStatus:
		//		var possibleStatus StateStatus
		//		var state State
		//		for _, possible := range r.state.GetPossibles() {
		//			possible.SetEvent(event)
		//			possibleStatus = possible.Read(r.txt)
		//			if possibleStatus != FailureStatus {
		//				possible.SetId(r.state.GetId())
		//				state = possible
		//				status = possibleStatus
		//				break
		//			}
		//		}
		//		if state == nil {
		//			txt.Cmd("511 Syntax error, command unrecognized")
		//		} else {
		//			r.state = state
		//			goto handleStatus
		//		}
		//	}
	}
}
