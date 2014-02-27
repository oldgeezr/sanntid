package main

import (
	//. "./algorithm"
	. "./lift"
	. "./lift/log"
	. "./messages"
	. "./network"
	. "./network/tcp"
	. "./network/udp"
	// . "fmt"
	. "net"
	"time"
)

func main() {

	saddr, _ := ResolveUDPAddr("udp", UDP_PORT)
	ln, _ := ListenUDP("udp", saddr)
	ln.SetReadDeadline(time.Now().Add(250 * time.Millisecond))

	b := make([]byte, 16)

	_, _, err := ln.ReadFromUDP(b)
	ln.Close()

	array_update := make(chan int)
	get_array := make(chan []int)
	new_master := make(chan bool)
	flush := make(chan bool)
	master := make(chan bool)
	int_order := make(chan Dict)
	ext_order := make(chan Dict)
	last_order := make(chan Dict)
	new_job_queue := make(chan string)
	master_request := make(chan string)
	master_order := make(chan Dict)
	master_pop := make(chan string)
	algo_out := make(chan Order)
	last_floor := make(chan Dict)
	get_last_queue := make(chan []Dict)
	get_last_queue_request := make(chan bool)
	job_queue := make(chan []Jobs)
	last_queue := make(chan []Dict)

	go IP_array(array_update, get_array, flush)
	// Println("Starter IP_array...")
	go Timer(flush)
	// Println("Starter Timer...")
	go Last_queue(last_floor, get_last_queue, get_last_queue_request, new_job_queue)
	go Job_queues(new_job_queue, master_request, master_pop, master_order, algo_out)

	if err != nil { // MASTER
		go IMA(master, get_array)
		// Println("Starter IMA...")
		master <- true
		go UDP_listen(array_update)
		// Println("Starter UDP_listen...")
	} else { // SLAVE
		go Internal(int_order, ext_order, last_order)
		// Println("slave")
		go IMA(master, get_array)
		// Println("Starter IMA...")
		master <- false
		go UDP_listen(array_update)
		// Println("Starter UDP_listen...")
		go IMA_master(get_array, master, new_master)
		// Println("Starter IMA_master...")
		go Connect_to_MASTER(get_array, new_master, int_order, ext_order, last_order)
		new_master <- true
	}

	neverQuit := make(chan string)
	<-neverQuit
}
