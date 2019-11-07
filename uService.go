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
	uClientList list.List
	umu sync.Mutex
)

func main() {
	print("bind listen port:")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	service := ":"+input.Text()

	print("service working ...\n")
	umsgWorking(service)

}
func umsgWorking(service string) {
	udpaddr, _ := net.ResolveUDPAddr("udp", service)
	udpconn, _ := net.ListenUDP("udp", udpaddr)

	var rbuf [512]byte
	for {
		msgLen, remoteAddr, err := udpconn.ReadFromUDP(rbuf[0:])
		if err != nil {
			return
		}
		inStr := strings.TrimSpace(string(rbuf[0:msgLen]))
		inputs := strings.Split(inStr, " ")

		switch inputs[0] {

		case "GetAddrs":
			print("GetAddrs\n")
			wStr := string("onlineClient")
			umu.Lock()
			for item := uClientList.Front();nil != item ;item = item.Next() {
				interfaceValue := fmt.Sprintf("%s",item.Value)
				if interfaceValue == strings.TrimSpace(fmt.Sprintf("%s",remoteAddr)){
					continue
				}
				wStr = strings.Join([]string{wStr,interfaceValue},"##")
			}
			umu.Unlock()
			data := []byte(wStr)
			_,_ = udpconn.WriteToUDP(data, remoteAddr)
			fmt.Println("push addr to client:", remoteAddr,"\n",wStr)

			break

		case "ClientOn":
			print("ClientOn\n")
			umu.Lock()
			if uClientList.Len() == 0{
				uClientList.PushBack(remoteAddr)
				fmt.Println("client:", remoteAddr," on (new)")
			}else{
				i := 0
				localIpAndPort  := strings.Split(strings.TrimSpace(fmt.Sprintf("%s",remoteAddr)),":")
				for item := uClientList.Front();nil != item ;item = item.Next() {
					remoteIpAndPort := strings.Split(strings.TrimSpace(fmt.Sprintf("%s",item.Value)),":")
					if remoteIpAndPort[0] ==  localIpAndPort[0]{
						if remoteIpAndPort[1] ==  localIpAndPort[1]{
							fmt.Println("client:",remoteAddr," on (no change)")
							i=1
						} else {
							uClientList.PushBack(remoteAddr)
							fmt.Println("client:", remoteAddr," on (ip exist)")
							i=1
						}
						//} else {
						//	ClientList.Remove(item)
						//	ClientList.PushBack(remoteAddr)
						//	fmt.Println("client:", remoteAddr," on (port changed)")
						//	i=1
						//}
						break
					}
				}
				if i != 1 {
					uClientList.PushBack(remoteAddr)
					fmt.Println("client:", remoteAddr, " on (new)")
				}
			}
			umu.Unlock()
			break

		case "ClientDrop":
			umu.Lock()
			for item := uClientList.Front();nil != item ;item = item.Next() {
				if item.Value ==  remoteAddr{
					uClientList.Remove(item)
					break
				}
			}
			fmt.Println("client:", remoteAddr," drop")
			umu.Unlock()
			break

		case "KeepAlive":
			fmt.Println("remoteaddr :", remoteAddr," keep alive")
			break
		}
		time.Sleep(1)
	}
}