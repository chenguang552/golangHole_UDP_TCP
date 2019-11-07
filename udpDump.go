package main

import (
	"bufio"
	"container/list"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

var uipList list.List

func main() {

	print("bind port:")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	service := ":"+input.Text()
	remote := readServiceAddr("service.txt")
	print("service addr:"+remote+"\n")
	if remote == ""{
		fmt.Printf("Fail to get service addr .\n")
		return
	}
	uService, _ := net.ResolveUDPAddr("udp", service)
	uRemote, _ := net.ResolveUDPAddr("udp", remote)
	conn, err := net.DialUDP("udp", uService, uRemote)
	if err != nil {
		fmt.Printf("Fail to connect, %s\n", err)
		return
	}
	uhandleRequest(conn)
	uhandleHoleWork(service, uipList.Front().Value.(string))
}

func uhandleHoleWork(service string, remote string){
	//udpaddr, _ := net.ResolveUDPAddr("udp", service)
	//udpconn, _ := net.ListenUDP("udp", udpaddr)
	uService, _ := net.ResolveUDPAddr("udp", service)
	uRemote, _ := net.ResolveUDPAddr("udp", remote)
	conn2, _ := net.DialUDP("udp", uService, uRemote)
	defer conn2.Close()
	go getPeerMsg(conn2)
	//for i := 0; i < 10; i++ {
	for{
		print("send to peer 'hello'\n")
		_, err2 := conn2.Write([]byte("hello\r\n"))
		if err2 != nil {
			fmt.Print(err2)
		}

		time.Sleep(time.Duration(2)*time.Second)
	}
}

func getPeerMsg(conn *net.UDPConn) {
	var resp2 [512]byte
	for {
		respLen, rAddr, _ := conn.ReadFromUDP(resp2[0:])
		print("peerMsg from[",rAddr.IP.String(),":",rAddr.Port,"] ",string(resp2[0:respLen]), "\n")
	}
}

func uhandleRequest(conn *net.UDPConn) {
	defer conn.Close()
	//上线
	_, err := conn.Write([]byte("ClientOn"))
	if err != nil {
		return
	}
	time.Sleep(1)
	// 获取ip列表
	var resp [512]byte
	var respLen int
	var respErr error
	for {
		_, err = conn.Write([]byte("GetAddrs"))
		respLen, _, respErr = conn.ReadFromUDP(resp[0:])
		//respLen, respErr = conn.Read(resp[0:])
		if respErr != nil {
			return
		}
		if err != nil {
			return
		}
		tmp := strings.TrimSpace(string(resp[0:respLen]))
		if len(tmp) > 14 {
			break
		}
		print("no ip list , try again .\n")
		time.Sleep(time.Duration(2)*time.Second)
	}
	inStr := strings.TrimSpace(string(resp[0:respLen]))
	inputs := strings.Split(inStr, "##")
	for i := 1; i < len(inputs); i++{
		uipList.PushBack(inputs[i])
		print(inputs[i]+"\n")
	}

}

func readServiceAddr(fN string) string {
	server,fErr := os.Open(fN)
	defer server.Close()
	if fErr != nil{
		fmt.Printf("%s",fErr)
	}
	fileScanner := bufio.NewScanner(server)
	for fileScanner.Scan() {
		//print("fileScanner.Scan()\n")
		// 以#开头视为注释 空行和注释不读取
		if strings.HasPrefix(fileScanner.Text(),"#") || fileScanner.Text() == "" {
			continue
		} else {
			return fileScanner.Text()
		}
	}
	return ""
}