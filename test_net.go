package main 

import(
	."fmt"
	."net"
	"time"
	"strings"
)

func send_message(conn Conn) {
        
        str := GetMyIP()
        _, err := conn.Write([]byte(str))
        _ = err
}

func udp_listen(ch chan bool) {

        saddr, _ := ResolveUDPAddr("udp", ":10020")        
        ln, _ := ListenUDP("udp", saddr)
        
        for {
                b := make([]byte,16)
                _, _, err := ln.ReadFromUDP(b)
		remoteIP := string(b[0:15]) 
                if err == nil {
                        time.Sleep(100*time.Millisecond)
                        ch<- true
                }
               
                if remoteIP != GetMyIP() {
                	Println(remoteIP)
                }
        }
}

func udp_send(ch chan bool) {

        saddr, _ := ResolveUDPAddr("udp","129.241.187.255:10020")
        conn, _ := DialUDP("udp", nil, saddr)
        
        for {
                if <-ch == true {
                    	send_message(conn)
                } else {
                    	time.Sleep(100*time.Millisecond)
                    	send_message(conn)
                }
        }        
}

func GetMyIP() string {

        allIPs, _ := InterfaceAddrs()
        
        IPString := make([]string, len(allIPs))
        for i := range allIPs {
                temp := allIPs[i].String()
                ip := strings.Split(temp, "/")
                IPString[i] = ip[0]
        }
        var myIP string
        for i:=range IPString {
                if IPString[i][0:3] == "129" {
                        myIP = IPString[i]
                }
        }
        return myIP
}

func network_modul() {	
	
	ch := make(chan bool)

	go udp_listen(ch)
	go udp_send(ch)

	ch<- true
}

func main() {

	go network_modul()
	
	neverQuit := make(chan string)
	<-neverQuit
}
