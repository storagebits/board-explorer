package main

import (
	"context"
	"flag"
	"fmt"
	evdev "golang-evdev"
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

func readInputEvents(inputDev *evdev.InputDevice, messages chan byte) {
	// read events comming from the overlay touch device

	var events []evdev.InputEvent
	var err error

	for {
		events, err = inputDev.Read()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		for i := range events {
			str := formatEvent(&events[i])
			/*if strings.Contains(str, "ABS_MT_TRACKING_ID") || strings.Contains(str, "ABS_MT_POSITION_X") || strings.Contains(str, "ABS_MT_POSITION_Y") || strings.Contains(str, "SYN") {
			if !strings.Contains(laststr, "SYN") {*/
			fmt.Println(str)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func formatEvent(ev *evdev.InputEvent) string {
	// format events comming from the overlay touch device

	var res, f, codeName string

	code := int(ev.Code)
	etype := int(ev.Type)

	switch ev.Type {
	case evdev.EV_SYN:
		if ev.Code == evdev.SYN_MT_REPORT {
			f = "time %d.%-8d +++++++++ %s ++++++++"
		} else {
			f = "time %d.%-8d --------- %s --------"
		}
		return fmt.Sprintf(f, ev.Time.Sec, ev.Time.Usec, evdev.SYN[code])
	case evdev.EV_KEY:
		val, haskey := evdev.KEY[code]
		if haskey {
			codeName = val
		} else {
			val, haskey := evdev.BTN[code]
			if haskey {
				codeName = val
			} else {
				codeName = "?"
			}
		}
	default:
		m, haskey := evdev.ByEventType[etype]
		if haskey {
			codeName = m[code]
		} else {
			codeName = "?"
		}
	}

	evfmt := "time %d.%-8d type %d (%s), code %-3d (%s), value %d"
	res = fmt.Sprintf(evfmt, ev.Time.Sec, ev.Time.Usec, etype,
		evdev.EV[int(ev.Type)], ev.Code, codeName, ev.Value)

	return res
}

func readJoystick(js joystick.Joystick, messages chan byte) {
	jinfo, err := js.Read()
	// read joytick event and send corresponding byte to microbit

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
			// UP
			log.Printf("UP")
			messages <- 0x7a
			log.Printf("UP sent")
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

	// parse flags
	microbitName1 := flag.String("microbitName1", "", "Name of microbit 1")
	microbitName2 := flag.String("microbitName2", "", "Name of microbit 2 (optional)")
	flag.Parse()

	if *microbitName1 == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// init inputDevice (Multi touch overlay device)
	/*var inputDev *evdev.InputDevice
	var err error
	multiTouchChannel := make(chan byte)

	inputDev, err = evdev.Open("/dev/input/event0")
	log.Printf("Evdev protocol version: %d\n", inputDev.EvdevVersion)
	log.Printf("Device name: %s\n", inputDev.Name)
	go readInputEvents(inputDev, multiTouchChannel)*/

	// Init BLE controller
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

	log.Println("connecting BLE devices...")

	// Init Player 1 BLE and joytick
	player1, err := ble.Connect(ctx, func(a ble.Advertisement) bool {
		if a.Connectable() && strings.HasPrefix(a.LocalName(), "BBC micro:bit ["+*microbitName1+"]") && strings.Contains(a.LocalName(), *microbitName) {
			log.Printf("connect to %s", a.LocalName())
			return true
		}
		return false
	})
	if err != nil {
		log.Fatalf("failed to connect: %s", err)
	}
	go func() {
		<-player1.Disconnected()
		cancel()
	}()

	p, err := player1.DiscoverProfile(true)
	if err != nil {
		log.Fatalf("failed to discover profile: %s", err)
	}

	// Select the UART for writing
	c := p.FindCharacteristic(ble.NewCharacteristic(ble.MustParse(`6E400003-B5A3-F393-E0A9-E50E24DCCA9E`)))

	// Init joystick
	jsid := 0
	js, jserr := joystick.Open(jsid)
	if jserr != nil {
		fmt.Println(jserr)
		return
	}

	// Init player 2 BLE and joystick
	var player2 ble.Client
	var c2 *ble.Characteristic
	var js2 joystick.Joystick

	if *microbitName2 != "" {
		player2, err := ble.Connect(ctx, func(a ble.Advertisement) bool {
			if a.Connectable() && strings.HasPrefix(a.LocalName(), "BBC micro:bit ["+*microbitName2+"]") && strings.Contains(a.LocalName(), *microbitName) {
				log.Printf("connect to %s", a.LocalName())
				return true
			}
			return false
		})
		if err != nil {
			log.Fatalf("failed to connect: %s", err)
		}
		go func() {
			<-player2.Disconnected()
			cancel()
		}()

		p2, err := player2.DiscoverProfile(true)
		if err != nil {
			log.Fatalf("failed to discover profile: %s", err)
		}

		// Select the UART for writing
		c2 = p2.FindCharacteristic(ble.NewCharacteristic(ble.MustParse(`6E400003-B5A3-F393-E0A9-E50E24DCCA9E`)))

		jsid = 1
		js2, jserr = joystick.Open(jsid)
		if jserr != nil {
			fmt.Println(jserr)
			return
		}

	}

	ticker := time.NewTicker(time.Millisecond * 40)

	channelBlePlayer1 := make(chan byte)
	channelBlePlayer2 := make(chan byte)
	for {
		select {
		case ev := <-channelBlePlayer1:
			log.Printf("Message received: %b", ev)
			if err := player1.WriteCharacteristic(c, []byte{ev, 0x0a}, true); err != nil {
				log.Printf("send data: %s", err)
			}
		case ev2 := <-channelBlePlayer2:
			if err := player2.WriteCharacteristic(c2, []byte{ev2, 0x0a}, true); err != nil {
				log.Printf("send data: %s", err)
			}
		/* case ev3 := <-multiTouchChannel:
		if err := player2.WriteCharacteristic(c2, []byte{ev3, 0x0a}, true); err != nil {
			log.Printf("send data: %s", err)
		}*/
		case <-ticker.C:
			go readJoystick(js, channelBlePlayer1)
			go readJoystick(js2, channelBlePlayer2)
		default:
			//fmt.Println("no message received")
		}
	}
}
