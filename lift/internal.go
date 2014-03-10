package lift

import (
	. ".././driver"
	. ".././network"
	// . ".././formating"
	// . "./log"
	. ".././functions"
	. "fmt"
	. "strconv"
	"time"
)

func Do_first(do_first chan Queues, order chan Dict, kill_send_to_floor chan bool) {

	var last_floor int
	// var doing Dict
	var running bool = false

	Fo.WriteString("Entered Do_first\n")

	myIP := GetMyIP()

	for {
		time.Sleep(50 * time.Millisecond)

		do := <-do_first

		Println("TODOLIST:", do)

		job_queue := do.Int_queue
		last_queue := do.Last_queue
		ext_queue := do.Ext_queue

		for _, last := range last_queue {
			if last.Ip_order == myIP {
				last_floor = last.Floor
				// last_dir = last.Dir
				break
			}
		}

		Println(running)

		if len(job_queue) != 0 {
			for _, yours := range job_queue {
				if yours.Ip == GetMyIP() {
					if len(yours.Dest) != 0 {
						// if yours.Dest[0] != doing {
						// Println("YOU ARE DOING:", doing, yours.Dest[0])
						// Println("Running:", running)
						// doing = yours.Dest[0]
						go func() {
							if running {
								kill_send_to_floor <- true
							}
						}()
						go func() {
							if !running {
								running = true
								running = Send_to_floor(yours.Dest[0].Floor, last_floor, "int", kill_send_to_floor)
							}
						}()
						// }
					} else {
						if len(ext_queue) != 0 {
							Println("ext", ext_queue, len(ext_queue))
							if ext_queue[0] != doing {
								doing = ext_queue[0]
								go func() {
									if running {
										kill_send_to_floor <- true
									}
								}()
								go func() {
									running = true
									running = Send_to_floor(yours.Dest[0].Floor, last_floor, "int", kill_send_to_floor)
								}()
							}
						} else {
							order <- Dict{myIP, M + 1, "standby"}
						}
					}
				}
			}
		}
	}
}

//Sends elevator to specified floor
func Send_to_floor(floor, current_floor int, button string, kill_send_to_floor chan bool) bool {

	Elev_set_door_open_lamp(0)
	Set_stop_lamp(0)

	var stop bool = false

	if current_floor < floor {
		for {
			Speed(150)
			go func() {
				stop = <-kill_send_to_floor
				Println("recieved kill")
			}()
			if Get_floor_sensor() == floor || stop {
				Set_stop_lamp(1)
				Elev_set_door_open_lamp(1)
				Speed(-150)
				time.Sleep(25 * time.Millisecond)
				Speed(0)
				time.Sleep(1500 * time.Millisecond)
				if button == "int" {
					Set_button_lamp(BUTTON_COMMAND, floor, 0)
				} else {
					if button == "up" {
						Set_button_lamp(BUTTON_CALL_UP, floor, 0)
					} else {
						Set_button_lamp(BUTTON_CALL_DOWN, floor, 0)
					}
				}
				return false
			}
			time.Sleep(25 * time.Millisecond)
		}
	} else if current_floor > floor {
		for {
			Speed(-150)
			go func() {
				stop = <-kill_send_to_floor
				Println("recieved kill")
			}()
			if Get_floor_sensor() == floor || stop {
				Set_stop_lamp(1)
				Elev_set_door_open_lamp(1)
				Speed(150)
				time.Sleep(25 * time.Millisecond)
				Speed(0)
				time.Sleep(1500 * time.Millisecond)
				if button == "int" {
					Set_button_lamp(BUTTON_COMMAND, floor, 0)
				} else {
					if button == "up" {
						Set_button_lamp(BUTTON_CALL_UP, floor, 0)
					} else {
						Set_button_lamp(BUTTON_CALL_DOWN, floor, 0)
					}
				}
				return false
			}
			time.Sleep(25 * time.Millisecond)
		}
	}
	return false
}

//Keyboard terminal input (For testing)
func KeyboardInput(ch chan int) {
	var a int

	for {
		Scan(&a)
		ch <- a
	}
}

//Handles external button presses
func Ext_order(order chan Dict) {

	Fo.WriteString("Entered Ext_order\n")

	i := 0

	for {

		if i < 3 {
			if Get_button_signal(BUTTON_CALL_UP, i) == 1 {
				Println("External call up button nr: " + Itoa(i) + " has been pressed!")
				Set_button_lamp(BUTTON_CALL_UP, i, 1)
				order <- Dict{"ext", i, "up"}
				time.Sleep(300 * time.Millisecond)
			}
		}
		if i > 0 {
			if Get_button_signal(BUTTON_CALL_DOWN, i) == 1 {
				Println("External call down button nr: " + Itoa(i) + " has been pressed!")
				Set_button_lamp(BUTTON_CALL_DOWN, i, 1)
				order <- Dict{"ext", i, "down"}
				time.Sleep(300 * time.Millisecond)
			}
		}

		i++
		i = i % 4
		time.Sleep(25 * time.Millisecond)

	}
}

//Handles internal button presses
func Int_order(order chan Dict) {

	Fo.WriteString("Entered Int_order\n")

	i := 0
	for {
		if Get_button_signal(BUTTON_COMMAND, i) == 1 {
			Println("Internal button nr: " + Itoa(i) + " has been pressed!")
			Set_button_lamp(BUTTON_COMMAND, i, 1)
			order <- Dict{GetMyIP(), i, "int"}
			Fprintln(Fo, "INTERNAL: btn -> order -> tcp")
			time.Sleep(300 * time.Millisecond)
		}

		i++
		i = i % 4
		time.Sleep(25 * time.Millisecond)

	}
}

//Checks which floor the elevator is on and sets the floor-light
func Floor_indicator(order chan Dict) {

	Fo.WriteString("Entered Floor_indicator\n")

	Println("executing floor indicator!")
	var floor int
	for {
		floor = Get_floor_sensor()
		if floor != -1 {
			Set_floor_indicator(floor)
			order <- Dict{GetMyIP(), floor, "standby"}
			// Fprintln(Fo, "222: @floor -> order -> tcp")
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func To_nearest_floor() {

	Fo.WriteString("Entered To_nearest_floor\n")

	for {
		Speed(150)
		if Get_floor_sensor() != -1 {
			time.Sleep(25 * time.Millisecond)
			Speed(0)
		}
	}
}

func Internal(order chan Dict) {

	Fo.WriteString("Entered Internal\n")

	// Initialize
	Init()
	Speed(150)
	floor := -1
	Println("UP")

	go func() {
		for {

			floor = Get_floor_sensor()

			if floor != -1 {

				Speed(-150)
				Println("DOWN")
				time.Sleep(10 * time.Millisecond)
				Println("STOP")
				Speed(0)
				return
			}
		}
	}()

	go Floor_indicator(order)
	go Int_order(order)
	go Ext_order(order)
}
