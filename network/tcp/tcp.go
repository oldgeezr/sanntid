package tcp

import (
	. "../.././functions"
	. "../.././network"
	"encoding/json"
	. "fmt"
	. "net"
	. "strconv"
	"time"
)

func TCP_master_connect(order, master_order chan Dict, queues chan Queues) {

	ln, _ := Listen("tcp", TCP_PORT)
	for {
		conn, _ := ln.Accept()
		go TCP_master_com(conn, order, master_order, queues)
	}
}

func TCP_master_com(conn Conn, order, master_order chan Dict, queues chan Queues) {

	for {
		select {
		case msg := <-queues:
			b, _ := json.Marshal(msg)
			conn.Write(b)
		case msg := <-order:
			master_order <- msg
		default:
			b := make([]byte, BUF_LEN)
			conn.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
			length, err := conn.Read(b)
			Println("master_err:", err)
			/*if err != nil {
				Println("closed connection")
				return
			}
			var c Dict
			json.Unmarshal(b[0:length], &c)
			master_order <- c*/
		}
	}
}

func Connect_to_MASTER(get_ip_array chan []int, new_master chan bool, order chan Dict, queues chan Queues) {

	for {
		select {
		case <-new_master:
			ip := <-get_ip_array
			if len(ip) != 0 {
				if ip[len(ip)-1] > 255 {
					master_ip := ip[len(ip)-1] - 255
					go TCP_slave_com(Itoa(master_ip), order, queues)
				}
			}
		default:
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func TCP_slave_com(master_ip string, order chan Dict, queues chan Queues) {

	conn, _ := Dial("tcp", IP_BASE+master_ip+TCP_PORT)

	for {
		select {
		case msg := <-order:
			Println(msg)
			b, _ := json.Marshal(msg)
			conn.Write(b)
		default:
			b := make([]byte, BUF_LEN)
			conn.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
			length, err := conn.Read(b)
			Println("slave_err:", err)
			/*if err != nil {
				Println("closed connection")
				return
			}
			var c Queues
			json.Unmarshal(b[0:length], &c)
			queues <- c*/
		}
	}
}

/*func TCP_slave_recieve(conn Conn, queues chan Queues) {

	for {
		b := make([]byte, BUF_LEN)
		length, _ := conn.Read(b)
		var c Queues
		json.Unmarshal(b[0:length], &c)
		queues <- c
	}
}*/
