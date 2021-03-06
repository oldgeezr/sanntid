package log

import (
	. "../.././algorithm"
	. "../.././formatting"
	. "../.././functions"
	"sort"
)

const (
	M int = 4 // Number of floors
)

func Job_queues(log_order chan Dict, slave_queues, queues_to_tcp, do_first chan Queues) {

	Fo.WriteString("Entered Job_queues\n")

	last_queue := []Dict{}
	job_queue := []Jobs{}
	ext_queue := []Dict{}

	the_queue := Queues{job_queue, ext_queue, last_queue}

	for {
		select {
		case msg := <-log_order:
			switch {
			case msg.Dir == "int":
				job_queue = Append_to_correct_queue(job_queue, msg)

			case msg.Ip_order == "ext":
				ext_queue = Append_if_missing_ext_queue(ext_queue, msg.Floor, msg.Dir)

			case msg.Floor >= M:
				last_queue = Update_last_queue(last_queue, msg, true)

			case msg.Dir == "standby":
				if len(last_queue) != 0 {
					for _, last := range last_queue {
						if last.Ip_order != msg.Ip_order {
							job_queue = Create_job_queue_if_missing(job_queue, msg.Ip_order)
						}
					}
				} else {
					job_queue = Create_job_queue_if_missing(job_queue, msg.Ip_order)
				}
				last_queue = Update_last_queue(last_queue, msg, false)
			}

			the_queue = Queues{}
			the_queue = Queues{job_queue, ext_queue, last_queue}
			the_queue = Algo(the_queue, msg)

			job_queue = the_queue.Job_queue
			ext_queue = the_queue.Ext_queue
			last_queue = the_queue.Last_queue

			Format_queues_term(the_queue, "MASTER")

		case msg := <-slave_queues:
			the_queue.Job_queue = msg.Job_queue
			the_queue.Ext_queue = msg.Ext_queue
			the_queue.Last_queue = msg.Last_queue
			Format_queues_term(the_queue, "SLAVE")

		case queues_to_tcp <- the_queue:
		case do_first <- the_queue:
		}
	}
}

func IP_array(ip_array_update chan int, get_ip_array chan []int, flush chan bool) {

	Fo.WriteString("Entered IP_array\n")

	IPaddresses := []int{}

	for {
		select {
		case ip := <-ip_array_update:
			IPaddresses = Append_if_missing_ip(IPaddresses, ip)
			sort.Ints(IPaddresses)

		case get_ip_array <- IPaddresses:
		case msg := <-flush:
			_ = msg
			IPaddresses = IPaddresses[:0]
		}
	}
}
