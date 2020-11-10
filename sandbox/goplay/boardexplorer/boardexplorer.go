package main

import (
	"fmt"
	"log"

	"github.com/simulatedsimian/joystick"
)

func readJoystick(js joystick.Joystick) {
	for {

		jinfo, err := js.Read()

		if err != nil {
			log.Printf("Error: " + err.Error())
			return
		}

		// BUTTONS
		/* printAt(1, 5, "Buttons:")
		for button := 0; button < js.ButtonCount(); button++ {
			if jinfo.Buttons&(1<<uint32(button)) != 0 {
				printAt(10+button, 5, "X")
			} else {
				printAt(10+button, 5, ".")
			}
		}*/

		// AXE
		for axis := 0; axis < js.AxisCount(); axis++ {
			//printAt(1, axis+7, fmt.Sprintf("Axis %2d Value: %7d", axis, jinfo.AxisData[axis]))
			//log.Printf("Axis %2d Value: %7d", axis, jinfo.AxisData[axis])
			if axis == 0 && jinfo.AxisData[axis] == 32767 {
				log.Printf("UP")
			}
		}
	}
}

func main() {
	log.Printf("Welcome to boardexplorer ! have fun :)")

	// Init joystick
	jsid := 0
	/*if len(os.Args) > 1 {
		i, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		jsid = i
	}*/

	js, jserr := joystick.Open(jsid)

	if jserr != nil {
		fmt.Println(jserr)
		return
	}

	ch := make(chan byte, 1)
	go readJoystick(js)
	<-ch
}
