package main

import (
	"bufio"
	"container/list"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	ClientList list.List
	mu sync.Mutex
)

func main() {
	print("bind listen port:")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	service := ":"+input.Text()
	mylistener, err := net.Listen("tcp", service)
	if err != nil{
		return
	}
	print("service working ...\n")
	for {
		conn, aErr := mylistener.Accept()
		if aErr != nil {
			continue
		}

		go msgWorking(conn)
	}
}
func msgWorking(conn net.Conn) {

	var rbuf [512]byte
	for {
		msgLen, err := conn.Read(rbuf[0:])
		if err != nil {
			return
		}
		inStr := strings.TrimSpace(string(rbuf[0:msgLen]))
		inputs := strings.Split(inStr, " ")

		switch inputs[0] {

		case "GetAddrs":
			print("GetAddrs\n")
			wStr := string("onlineClient")
			mu.Lock()
			for item := ClientList.Front();nil != item ;item = item.Next() {
				interfaceValue := fmt.Sprintf("%s",item.Value)
				if interfaceValue == strings.TrimSpace(fmt.Sprintf("%s",conn.RemoteAddr())){
					continue
				}
				wStr = strings.Join([]string{wStr,interfaceValue},"##")
			}
			mu.Unlock()
			data := []byte(wStr)
			conn.Write(data)
			fmt.Println("push addr to client:", conn.RemoteAddr(),"\n",wStr)

			break

		case "ClientOn":
			print("ClientOn\n")
			mu.Lock()
			if ClientList.Len() == 0{
				ClientList.PushBack(conn.RemoteAddr())
				fmt.Println("client:", conn.RemoteAddr()," on (new)")
			}else{
				i := 0
				localIpAndPort  := strings.Split(strings.TrimSpace(fmt.Sprintf("%s",conn.RemoteAddr())),":")
				for item := ClientList.Front();nil != item ;item = item.Next() {
					remoteIpAndPort := strings.Split(strings.TrimSpace(fmt.Sprintf("%s",item.Value)),":")
					if remoteIpAndPort[0] ==  localIpAndPort[0]{
						if remoteIpAndPort[1] ==  localIpAndPort[1]{
							fmt.Println("client:", conn.RemoteAddr()," on (no change)")
							i=1
						} else {
							ClientList.PushBack(conn.RemoteAddr())
							fmt.Println("client:", conn.RemoteAddr()," on (ip exist)")
							i=1
						}
						//} else {
						//	ClientList.Remove(item)
						//	ClientList.PushBack(conn.RemoteAddr())
						//	fmt.Println("client:", conn.RemoteAddr()," on (port changed)")
						//	i=1
						//}
						break
					}
				}
				if i != 1 {
					ClientList.PushBack(conn.RemoteAddr())
					fmt.Println("client:", conn.RemoteAddr(), " on (new)")
				}
			}
			mu.Unlock()
			break

		case "ClientDrop":
			mu.Lock()
			for item := ClientList.Front();nil != item ;item = item.Next() {
				if item.Value ==  conn.RemoteAddr(){
					ClientList.Remove(item)
					break
				}
			}
			fmt.Println("client:", conn.RemoteAddr()," drop")
			conn.Close()
			mu.Unlock()
			break

		case "KeepAlive":
			fmt.Println("remoteaddr :", conn.RemoteAddr()," keep alive")
			break
		}
		time.Sleep(1)
	}
}