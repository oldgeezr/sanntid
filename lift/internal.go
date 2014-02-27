package lift

import (
	. ".././driver"
	. ".././network"
	. "./log"
	. "fmt"
	. "strconv"
	"time"
)

//Sends elevator to specified floor
func Send_to_floor(floor int, button string) {
	current_floor := Get_floor_sensor()
	Elev_set_door_open_lamp(0)
	Set_stop_lamp(0)

	if current_floor < floor {
		Println("Going up")
		for {
			Speed(150)
			if Get_floor_sensor() == floor {
				Println("I am now at floor: " + Itoa(Get_floor_sensor()))
				Set_stop_lamp(1)
				Elev_set_door_open_lamp(1)
				Speed(-150)
				time.Sleep(25 * time.Millisecond)
				Speed(0)
				if button == "int" {
					Set_button_lamp(BUTTON_COMMAND, floor, 0)
				} else {
					if button == "up" {
						Set_button_lamp(BUTTON_CALL_UP, floor, 0)
					} else {
						Set_button_lamp(BUTTON_CALL_DOWN, floor, 0)
					}
				}
				return
			}
			time.Sleep(25 * time.Millisecond)
		}
	} else {
		Println("Going down")
		for {
			Speed(-150)
			if Get_floor_sensor() == floor {
				Println("I am now at floor: " + Itoa(Get_floor_sensor()))
				Set_stop_lamp(1)
				Elev_set_door_open_lamp(1)
				Speed(150)
				time.Sleep(25 * time.Millisecond)
				Speed(0)
				if button == "int" {
					Set_button_lamp(BUTTON_COMMAND, floor, 0)
				} else {
					if button == "up" {
						Set_button_lamp(BUTTON_CALL_UP, floor, 0)
					} else {
						Set_button_lamp(BUTTON_CALL_DOWN, floor, 0)
					}
				}
				return
			}
			time.Sleep(25 * time.Millisecond)
		}
	}
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
func Ext_order(ext_order chan Dict) {

	i := 0

	for {

		if i < 3 {
			if Get_button_signal(BUTTON_CALL_UP, i) == 1 {
				Println("External call up button nr: " + Itoa(i) + " has been pressed!")
				Set_button_lamp(BUTTON_CALL_UP, i, 1)
				ext_order <- Dict{"up", i}
				time.Sleep(300 * time.Millisecond)
			}
		}
		if i > 0 {
			if Get_button_signal(BUTTON_CALL_DOWN, i) == 1 {
				Println("External call down button nr: " + Itoa(i) + " has been pressed!")
				Set_button_lamp(BUTTON_CALL_DOWN, i, 1)
				ext_order <- Dict{"down", i}
				time.Sleep(300 * time.Millisecond)
			}
		}

		i++
		i = i % 4
		time.Sleep(25 * time.Millisecond)

	}
}

//Handles internal button presses
func Int_order(int_order chan Dict) {

	i := 0
	for {
		if Get_button_signal(BUTTON_COMMAND, i) == 1 {
			Println("Internal button nr: " + Itoa(i) + " has been pressed!")
			Set_button_lamp(BUTTON_COMMAND, i, 1)
			int_order <- Dict{GetMyIP(), i}
			time.Sleep(300 * time.Millisecond)
		}

		i++
		i = i % 4
		time.Sleep(25 * time.Millisecond)

	}
}

//Checks which floor the elevator is on and sets the floor-light
func Floor_indicator(last_order chan Dict) {
	Println("executing floor indicator!")
	//_ = last_order
	var floor int
	for {
		floor = Get_floor_sensor()
		if floor != -1 {
			Set_floor_indicator(floor)
			last_order <- Dict{GetMyIP(), i}
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func To_nearest_floor() {
	for {
		Speed(150)
		if Get_floor_sensor() != -1 {
			time.Sleep(25 * time.Millisecond)
			Speed(0)
		}
	}
}

func Internal(int_order, ext_order, last_order chan Dict) {

	// Initialize
	Init()
	Speed(150)
	last_floor := -1

	go func() {
		for {

			last_floor = Get_floor_sensor()

			if last_floor != -1 {

				Speed(-150)
				time.Sleep(25 * time.Millisecond)
				Speed(0)
				return
			}
		}
	}()

	go Floor_indicator(last_order)
	go Int_order(int_order)
	go Ext_order(ext_order)

	neverQuit := make(chan string)
	<-neverQuit

}
