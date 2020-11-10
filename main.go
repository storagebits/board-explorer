package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
	"github.com/simulatedsimian/joystick"
)

var (
	microbitName = flag.String("m", "", "microbit name")
)

var (
	UART_SERVICE_UUID = ble.MustParse(`6E400001-B5A3-F393-E0A9-E50E24DCCA9E`)
	TX_CHAR_UUID      = ble.MustParse(`6E400002-B5A3-F393-E0A9-E50E24DCCA9E`)
	RX_CHAR_UUID      = ble.MustParse(`6E400003-B5A3-F393-E0A9-E50E24DCCA9E`)
)

var c *ble.Characteristic

func readJoystick(js joystick.Joystick, messages chan byte) {
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
	if jinfo.Buttons&(1<<uint32(0)) != 0 {
		log.Printf("BUTTON PRESSED")
		messages <- 0x6f
		log.Printf("BUTTON PRESSED")
	} else {
		log.Printf("BUTTON RELEASED")
		messages <- 0x70
		log.Printf("BUTTON RELEASED")
	}

	// AXE
	for axis := 0; axis < js.AxisCount(); axis++ {
		//printAt(1, axis+7, fmt.Sprintf("Axis %2d Value: %7d", axis, jinfo.AxisData[axis]))
		//log.Printf("Axis %2d Value: %7d", axis, jinfo.AxisData[axis])
		if axis == 0 && jinfo.AxisData[axis] == 32767 {
			if jinfo.AxisData[1] == -32767 {
				// UPLEFT
				log.Printf("UPLEFT")
				messages <- 0x61
				log.Printf("UPLEFT sent")

			} else {
				// UP
				log.Printf("UP")
				messages <- 0x7a
				log.Printf("UP sent")
			}
		}

		if axis == 0 && jinfo.AxisData[axis] == -32767 {
			log.Printf("REVERSE")
			messages <- 0x73
			log.Printf("REVERSE sent")
		}

		if axis == 0 && jinfo.AxisData[axis] == 0 {
			log.Printf("BREAK")
			messages <- 0x77
			log.Printf("BREAK sent")
		}

		if axis == 1 && jinfo.AxisData[axis] == -32767 {
			log.Printf("LEFT")
			messages <- 0x71
			log.Printf("LEFT sent")
		}

		if axis == 1 && jinfo.AxisData[axis] == 32767 {
			log.Printf("RIGHT")
			messages <- 0x64
			log.Printf("LEFT sent")
		}

	}

	return
}

func main() {

	log.Printf("Welcome to board-explorer ! have fun :)")

	// Init BLE
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-sig
		cancel()
	}()

	device, err := dev.DefaultDevice()
	if err != nil {
		log.Fatal(err)
	}

	ble.SetDefaultDevice(device)

	log.Println("connecting...")

	client, err := ble.Connect(ctx, func(a ble.Advertisement) bool {
		if a.Connectable() && strings.HasPrefix(a.LocalName(), "BBC micro:bit [tavez]") && strings.Contains(a.LocalName(), *microbitName) {
			log.Printf("connect to %s", a.LocalName())
			return true
		}
		return false
	})
	if err != nil {
		log.Fatalf("failed to connect: %s", err)
	}
	go func() {
		<-client.Disconnected()
		cancel()
	}()

	p, err := client.DiscoverProfile(true)
	if err != nil {
		log.Fatalf("failed to discover profile: %s", err)
	}

	c = p.FindCharacteristic(ble.NewCharacteristic(RX_CHAR_UUID))

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

	ticker := time.NewTicker(time.Millisecond * 40)

	messages := make(chan byte)
	for {
		select {
		case ev := <-messages:
			log.Printf("Message received")
			if err := client.WriteCharacteristic(c, []byte{ev, 0x0a}, true); err != nil {
				log.Printf("send data: %s", err)
			}
		case <-ticker.C:
			go readJoystick(js, messages)
		default:
			//fmt.Println("no message received")
		}
	}
}
