package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gwbeacon/gwbeacon/lib"
	"github.com/gwbeacon/sdk/v1"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type client struct {
	sync.RWMutex
	v1.MessageServiceClient
	v1.QueryServiceClient
	v1.UserServiceClient
	v1.RosterServiceClient
	v1.MUCServiceClient
	msgCli  v1.MessageService_OnChatMessageClient
	idMaker lib.IDMaker
	waitAck map[uint64]*v1.ChatMessage
	errCh   chan error
	user    string
	domain  string
	device  string
	notice  string
}

func (c *client) Login() error {
	account := &v1.UserAccount{
		Domain: c.domain,
		Device: c.device,
		Name:   c.user,
		Passwd: "test",
	}
	info, err := c.SignIn(context.Background(), account)
	if err != nil {
		return err
	}
	log.Println(info)
	return nil
}

func (c *client) Run() error {
	ackCli, err := c.OnAckMessage(context.Background())
	if err != nil {
		return err
	}
	errCh := make(chan error, 3)
	go func() {
		for {
			ack, err := ackCli.Recv()
			if err != nil {
				errCh <- err
				return
			}
			id := uint64(ack.Id)
			c.Lock()
			if _, ok := c.waitAck[id]; ok {
				delete(c.waitAck, id)
			}
			c.Unlock()
		}
	}()
	hbCli, err := c.OnHeartbeat(context.Background())
	if err != nil {
		return err
	}
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		var err error
		for {
			<-ticker.C
			hb := &v1.Heartbeat{
				Domain:       c.domain,
				From:         c.user,
				Index:        uint64(c.idMaker.MakeID().GetIndex()),
				ClientTime:   time.Now().Unix(),
				NextDuration: 10,
			}
			if err = hbCli.Send(hb); err == nil {
				_, err = hbCli.Recv()
			}
			if err != nil {
				c.errCh <- err
				return
			}
		}
	}()
	msgCli, err := c.OnChatMessage(context.Background())
	if err != nil {
		return err
	}
	go func() {
		for {
			msg, err := msgCli.Recv()
			if err != nil {
				c.errCh <- err
				return
			}
			tm := lib.ID(msg.Id).GetTimestamp()
			fmt.Printf("\r%s@%s at %s:\n%s\n", msg.From, msg.Domain, time.Unix(int64(tm), 0).String(), msg.Msg)
			c.Lock()
			fmt.Printf(c.notice)
			c.Unlock()
			ack := &v1.AckMessage{
				Domain: c.domain,
				From:   c.user,
				To:     msg.From,
				Id:     msg.Id,
			}
			err = ackCli.Send(ack)

			if err != nil {
				c.errCh <- err
				return
			}
		}
	}()

	go func() {
		inputReader := bufio.NewReader(os.Stdin)
		for {
			var line string
			var to string
			var err error
			c.Lock()
			c.notice = "\rto >"
			c.Unlock()
			for {
				fmt.Printf(c.notice)
				to, _ = inputReader.ReadString('\n')
				to = strings.Replace(to, "\r", "", -1)
				to = strings.Replace(to, "\n", "", -1)
				if to != "" {
					break
				}
			}
			c.Lock()
			c.notice = "\rto " + to + " msg: "
			c.Unlock()
			for {
				fmt.Printf(c.notice)
				line, _ = inputReader.ReadString('\n')
				line = strings.Replace(line, "\r", "", -1)
				line = strings.Replace(line, "\n", "", -1)
				if line != "" {
					break
				}
			}
			msg := &v1.ChatMessage{
				Id:     uint64(c.idMaker.MakeID()),
				Domain: c.domain,
				From:   c.user,
				To:     to,
				Msg:    line,
			}
			c.Lock()
			c.notice = "\rto >"
			c.Unlock()
			err = msgCli.Send(msg)
			if err != nil {
				c.errCh <- err
				return
			}
		}
	}()
	select {
	case err := <-c.errCh:
		log.Println(err)
		ackCli.CloseSend()
		hbCli.CloseSend()
		msgCli.CloseSend()
		return err
	}

	return nil
}

func main() {
	var addr string
	var user string
	flag.StringVar(&addr, "server", "localhost:8888", "-server localhost:8888")
	flag.StringVar(&user, "user", "Alice", "-user Alice")
	flag.Parse()
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Println(err)
		return
	}
	cli := &client{
		domain:               "mi.com",
		device:               "goClient",
		user:                 user,
		idMaker:              lib.NewMessageIDMaker(0, 0),
		MessageServiceClient: v1.NewMessageServiceClient(conn),
		QueryServiceClient:   v1.NewQueryServiceClient(conn),
		UserServiceClient:    v1.NewUserServiceClient(conn),
		RosterServiceClient:  v1.NewRosterServiceClient(conn),
		MUCServiceClient:     v1.NewMUCServiceClient(conn),
	}

	err = cli.Login()
	if err != nil {
		log.Println(1, err)
		return
	}
	err = cli.Run()
	if err != nil {
		log.Println(2, err)
		return
	}
}
