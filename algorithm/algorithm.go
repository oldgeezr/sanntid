package algorithm

import (
	. "../functions"
	. "fmt"
)

func Algo(get_at_floor chan Dict, get_queues chan Queues) {

	Fo.WriteString("Entered Algo\n")

	var dir string

	for {
		select {
		case at_floor := <-get_at_floor:
			queues := <-get_queues
			Println(queues)
			int_queue := queues.Int_queue
			ext_queue := queues.Ext_queue
			last_queue := queues.Last_queue
			for i, order := range int_queue {
				Println("ranging")
				if order.Ip == at_floor.Ip_order { // Finn riktig kø
					Println("found correct queue:", order, at_floor.Floor)
					if !Missing_int_job(order, at_floor.Floor) { // Noen skal av
						// Stopp heis
						Println("queue before remove:", order)
						int_queue[i] = Remove_order_int_queue(int_queue[i], at_floor.Floor)
						Println("queue after remove:", order) // Slett alle interne
						ext_queue = Remove_order_ext_queue(ext_queue, at_floor.Floor, at_floor.Dir)
						Println("ext_queue after remove:", ext_queue) // Slett alle eksterne i riktig retning
					} else { // Ingen skal av
						Fprintln(Fo, "EXT: ",ext_queue)
						Fprintln(Fo, "EXT: ",at_floor.Floor, at_floor.Dir)

						if (ext_queue[0].Floor - at_floor.Floor) > 0 {
							dir = "up"
						} else if (ext_queue[0].Floor - at_floor.Floor) < 0 {
							dir = "down"
						} else {
							ext_queue = Remove_order_ext_queue(ext_queue, at_floor.Floor, dir)
						}
						if !Missing_ext_job(ext_queue, at_floor.Floor, dir) { // Noen skal på
							// Stopp heis
							Println("I was here?")
							ext_queue = Remove_order_ext_queue(ext_queue, at_floor.Floor, dir) // Slett alle eksterne i riktig retning
							Fprintln(Fo, "Removed from ext_queue")
						}
						Println("I was here??")
					}
					break // Avslutt å gå gjennom køen fordi det er unødvendig da det kun finnes en instans av hver heis
				}
			}
			queues = Queues{int_queue, ext_queue, last_queue}
			Println(queues)
			get_queues <- queues
			Fprintln(Fo, "\t \t \t \t ALGO: queue -> get_queues -> log")
		}
	}
}
