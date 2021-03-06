package Driver  // where "driver" is the folder that contains io.go, io.c, io.h, channels.go, channels.c and driver.go
/*
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
*/


//NOTE TO SELF: Denne er lik elev.go



import "C"

import (
	."fmt"
)


const N_BUTTONS = 3
const N_FLOORS = 4


type elev_motor_direction_t int
const ( 
    DIRN_DOWN elev_motor_direction_t = -1
    DIRN_STOP elev_motor_direction_t = 0
    DIRN_UP elev_motor_direction_t = 1
)

type elev_button_type_t int
const (
    BUTTON_CALL_UP elev_button_type_t = 0
    BUTTON_CALL_DOWN elev_button_type_t = 1
    BUTTON_COMMAND elev_button_type_t = 2
)




var lamp_channel_matrix [][]int
var button_channel_matrix [][]int




func fuckGO(){
	//Make 4x3 lightmatrix:
	lamp_channel_matrix = make ([][]int, N_FLOORS)
	for i := 0; i<N_FLOORS; i++ {
		lamp_channel_matrix [i] = make([]int, N_BUTTONS)
	}
	
	lamp_channel_matrix[0][0] = LIGHT_UP1
	lamp_channel_matrix[1][0] = LIGHT_UP2
	lamp_channel_matrix[2][0] = LIGHT_UP3
	lamp_channel_matrix[3][0] = LIGHT_UP4
	lamp_channel_matrix[0][1] = LIGHT_DOWN1
	lamp_channel_matrix[1][1] = LIGHT_DOWN2
	lamp_channel_matrix[2][1] = LIGHT_DOWN3
	lamp_channel_matrix[3][1] = LIGHT_DOWN4
	lamp_channel_matrix[0][2] = LIGHT_COMMAND1
	lamp_channel_matrix[1][2] = LIGHT_COMMAND2
	lamp_channel_matrix[2][2] = LIGHT_COMMAND3
	lamp_channel_matrix[3][2] = LIGHT_COMMAND4
	
	
	//Make 4x3 buttonmatrix:
	button_channel_matrix = make ([][] int, N_FLOORS)
	for i := 0; i<N_FLOORS; i++{
		button_channel_matrix [i] = make([]int, N_BUTTONS)
	}

	button_channel_matrix[0][0] = BUTTON_UP1
	button_channel_matrix[1][0] = BUTTON_UP2
	button_channel_matrix[2][0] = BUTTON_UP3
	button_channel_matrix[3][0] = BUTTON_UP4
	button_channel_matrix[0][1] = BUTTON_DOWN1
	button_channel_matrix[1][1] = BUTTON_DOWN2
	button_channel_matrix[2][1] = BUTTON_DOWN3
	button_channel_matrix[3][1] = BUTTON_DOWN4
	button_channel_matrix[0][2] = BUTTON_COMMAND1
	button_channel_matrix[1][2] = BUTTON_COMMAND2
	button_channel_matrix[2][2] = BUTTON_COMMAND3
	button_channel_matrix[3][2] = BUTTON_COMMAND4
	
}


	
	
	



func Elev_init() int {
	fuckGO() 
	
	//Init hardware
	if (io_init() == 0)	{return 0}
	
	// Zero all floor button lamps
	for i := 0; i<N_FLOORS; i++ {
		if i != 0 {
			Elev_set_button_lamp(BUTTON_CALL_DOWN, i, 0)
		}
		if i != N_FLOORS - 1 {
			Elev_set_button_lamp(BUTTON_CALL_UP, i, 0)
		}
		
		Elev_set_button_lamp(BUTTON_COMMAND, i, 0)
	}
	
	// Clear stop lamp, door open lamp, and set floor indicator to ground floor.
	elev_set_stop_lamp(0);
	Elev_set_door_open_lamp(0);
	Elev_set_floor_indicator(0);


	// Return success
	return 1;
}



func Elev_set_motor_direction(dirn int) {
	if dirn == 0 {
		io_write_analog(MOTOR, 0)
	} else if dirn > 0 {
		io_clear_bit(MOTORDIR)
		io_write_analog(MOTOR, 2800)
	} else if dirn < 0 {
		io_set_bit(MOTORDIR)
		io_write_analog(MOTOR, 2800)
	}
}


func Elev_set_door_open_lamp(value int) {
	if value == 1 {
		io_set_bit(LIGHT_DOOR_OPEN)
	} else {
		io_clear_bit(LIGHT_DOOR_OPEN)
	}
}


func elev_get_obstruction_signal() int {
	return io_read_bit(OBSTRUCTION)
}


func Elev_get_stop_signal() int {
	return io_read_bit(STOP)
}


func elev_set_stop_lamp(value int) {
	if value == 1 {
		io_set_bit(LIGHT_STOP)
	} else {
		io_clear_bit(LIGHT_STOP)
	}
}


func Elev_get_floor_sensor_signal() int {

	if (io_read_bit(SENSOR_FLOOR1) != 0) {
		return 0 
	} else if (io_read_bit(SENSOR_FLOOR2) != 0) {
		return 1 
	} else if (io_read_bit(SENSOR_FLOOR3) != 0) {
		return 2 
	} else if (io_read_bit(SENSOR_FLOOR4) != 0) {
		return 3 
	} else {
		return -1 
	}
}


func Elev_set_floor_indicator(floor int) {
	
	if floor < 0 {
		Println ("Floor is lower than 0\n")
		return
	}
	if floor > N_FLOORS {
		Println ("Floor variable is too high\n")
		return
	}
	
	// Binary encoding. One light must always be on.
	if ((floor & 0x02) != 0){
		io_set_bit(LIGHT_FLOOR_IND1)
	} else {
		io_clear_bit(LIGHT_FLOOR_IND1)
	}
		
	if ((floor & 0x01) != 0){
		io_set_bit(LIGHT_FLOOR_IND2)
	} else {
		io_clear_bit(LIGHT_FLOOR_IND2)
	}
}


func Elev_get_button_signal(button elev_button_type_t, floor int) bool {
	
	if floor < 0 {
		Println ("Floor is lower than 0\n")
		return false }
	if floor > N_FLOORS {
		Println ("Floor variable is too high\n")
		return false }
	if (button == BUTTON_CALL_UP && floor == N_FLOORS-1) {
		Println("Unvalid button call or floor\n")
		println("button 1")
		return false }
	if (button == BUTTON_CALL_DOWN && floor == 0) {
		Println("Unvalid button call or floor\n")
		println("button 2")
		return false }
	if !(button == BUTTON_CALL_UP || button == BUTTON_CALL_DOWN || button == BUTTON_COMMAND) {
		Println("Unvalid button\n")
		return false }
	

	if (io_read_bit(button_channel_matrix[floor][button]) != 0) {
		return true
	} else {
		return false
	}
}



func Elev_set_button_lamp(button elev_button_type_t, floor int, value int) {

	if floor < 0 {
		Println("Floor is negative\n")
		return }
	if floor >= N_FLOORS {
		Println("Floor variable is too high\n")
		return }
	if (button == BUTTON_CALL_UP && floor == N_FLOORS-1) {
		println("lampe 1")	
		Println("Unvalid button call or floor")
		return }
	if (button == BUTTON_CALL_DOWN && floor == 0) {
		println("lampe 2")
		Println("Unvalid button call or floor\n")
		return }
	if !(button == BUTTON_CALL_UP || button == BUTTON_CALL_DOWN || button == BUTTON_COMMAND) {
		Println("Unvalid button\n")
		return }

	
	if (value != 0) {
		io_set_bit(lamp_channel_matrix[floor][button])
	} else {
		io_clear_bit(lamp_channel_matrix[floor][button])
	}
}





































