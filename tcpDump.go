package main

import (
	"bufio"
	"container/list"
	"fmt"
	reuse "github.com/jbenet/go-reuseport"
	"net"
	"os"
	"strings"
	"time"
)

var ipList list.List

func main() {

	print("bind port:")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	service := ":"+input.Text()
	remote := "207.246.121.35:2019"
	//remote := "127.0.0.1:12306"
	//tcpAddr, _ := net.ResolveTCPAddr("tcp", service)
	//remoteAddr,_ := net.ResolveTCPAddr("tcp", remote)
	//conn, err := net.DialTCP("tcp", tcpAddr, remoteAddr)
	conn, err := reuse.Dial("tcp", service, remote)
	if err != nil {
		fmt.Printf("Fail to connect, %s\n", err)
		return
	}
	handleRequest(conn)
	//conn.Close()
	remoteIp := ipList.Front().Value.(string)
	//remoteIp := "207.246.121.35:80"
	//remoteAddr2,_ := net.ResolveTCPAddr("tcp", remoteIp)
	//conn2, _ := net.DialTCP("tcp", tcpAddr, remoteAddr2)
	go handleReListen(service)
	conn2, _ := reuse.Dial("tcp", service, remoteIp)
	var resp2 [512]byte
	for i := 0; i<10; i++ {
		_, err2 := conn2.Write([]byte("hello\r\n"))
		fmt.Print(err2)
		respLen, _ := conn2.Read(resp2[0:])
		print(string(resp2[0:respLen]))
		time.Sleep(time.Duration(2)*time.Second)
	}
	//conn.Close()
	conn2.Close()
}

func handleReListen(lService string){
	reListen, err := reuse.Listen("tcp",lService)
	if err != nil{
		fmt.Printf("Err: %s\n",err)
	}
	var tmpBuf [512]byte
	for{
		conn, lErr := reListen.Accept()
		if lErr != nil{
			continue
		}

		dLen,dErr := conn.Read(tmpBuf[0:])

		if dErr != nil{
			continue
		}

		fmt.Printf("%s\n", string(tmpBuf[0:dLen]))
	}
}
func handleRequest(conn net.Conn) {
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
		respLen, respErr = conn.Read(resp[0:])
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
		ipList.PushBack(inputs[i])
		print(inputs[i]+"\n")
	}

}