package main

import (
	"context"
	"flag"
	"fmt"
	evdev "github.com/gvalkov/golang-evdev"
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
	joystick1, joystick2 bool
)

// readInputEvents reads events from the touch input device and sends formatted messages
func readInputEvents(inputDev *evdev.InputDevice) {
	for {
		events, err := inputDev.Read()
		if err != nil {
			log.Printf("Error reading input events: %v", err)
			os.Exit(1)
		}

		for _, event := range events {
			fmt.Println(formatEvent(&event))
		}

		time.Sleep(100 * time.Millisecond)
	}
}

// formatEvent formats the input events for logging
func formatEvent(ev *evdev.InputEvent) string {
	switch ev.Type {
	case evdev.EV_SYN:
		if ev.Code == evdev.SYN_MT_REPORT {
			return fmt.Sprintf("time %d.%-8d +++++++++ %s ++++++++", ev.Time.Sec, ev.Time.Usec, evdev.SYN[ev.Code])
		}
		return fmt.Sprintf("time %d.%-8d --------- %s --------", ev.Time.Sec, ev.Time.Usec, evdev.SYN[ev.Code])

	case evdev.EV_KEY:
		codeName := evdev.KEY[ev.Code]
		return fmt.Sprintf("time %d.%-8d type %d (EV_KEY), code %-3d (%s), value %d", ev.Time.Sec, ev.Time.Usec, ev.Type, ev.Code, codeName, ev.Value)

	default:
		return fmt.Sprintf("time %d.%-8d type %d, code %-3d, value %d", ev.Time.Sec, ev.Time.Usec, ev.Type, ev.Code, ev.Value)
	}
}

// readJoystick reads joystick events and sends commands based on joystick input
func readJoystick(js joystick.Joystick, messages chan<- byte) {
	jinfo, err := js.Read()
	if err != nil {
		log.Printf("Error reading joystick: %v", err)
		return
	}

	if jinfo.Buttons&(1<<0) != 0 {
		messages <- 0x6f // Button pressed
	} else {
		messages <- 0x70 // Button released
	}

	// Handle axis input for directions
	switch {
	case jinfo.AxisData[0] == 32767:
		messages <- 0x7a // UP
	case jinfo.AxisData[0] == -32767:
		messages <- 0x73 // REVERSE
	case jinfo.AxisData[1] == -32767:
		messages <- 0x71 // LEFT
	case jinfo.AxisData[1] == 32767:
		messages <- 0x64 // RIGHT
	}
}

func main() {
	log.Printf("Welcome to board-explorer!")

	// Parse command line flags
	microbitName1 := flag.String("microbitName1", "", "Name of microbit 1")
	microbitName2 := flag.String("microbitName2", "", "Name of microbit 2 (optional)")
	flag.Parse()

	if *microbitName1 == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Setup signal handling to gracefully terminate
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-sig
		cancel()
	}()

	// Initialize BLE device
	device, err := dev.DefaultDevice()
	if err != nil {
		log.Fatalf("Failed to initialize BLE device: %v", err)
	}
	ble.SetDefaultDevice(device)

	// Connect to microbit 1 via BLE
	player1, err := ble.Connect(ctx, func(a ble.Advertisement) bool {
		return a.Connectable() && strings.HasPrefix(a.LocalName(), "BBC micro:bit ["+*microbitName1+"]")
	})
	if err != nil {
		log.Fatalf("Failed to connect to microbit 1: %v", err)
	}
	defer player1.CancelConnection()

	// Discover services and characteristics
	p1, err := player1.DiscoverProfile(true)
	if err != nil {
		log.Fatalf("Failed to discover profile: %v", err)
	}
	characteristic1 := p1.FindCharacteristic(ble.NewCharacteristic(ble.MustParse("6E400003-B5A3-F393-E0A9-E50E24DCCA9E")))

	// Initialize joystick 1
	js1, err := joystick.Open(0)
	if err != nil {
		log.Printf("Failed to open joystick 1: %v", err)
	} else {
		joystick1 = true
	}

	var (
		player2       ble.Client
		characteristic2 *ble.Characteristic
		js2           joystick.Joystick
	)

	// If microbit 2 is specified, connect to it as well
	if *microbitName2 != "" {
		player2, err = ble.Connect(ctx, func(a ble.Advertisement) bool {
			return a.Connectable() && strings.HasPrefix(a.LocalName(), "BBC micro:bit ["+*microbitName2+"]")
		})
		if err != nil {
			log.Fatalf("Failed to connect to microbit 2: %v", err)
		}
		defer player2.CancelConnection()

		p2, err := player2.DiscoverProfile(true)
		if err != nil {
			log.Fatalf("Failed to discover profile for player 2: %v", err)
		}
		characteristic2 = p2.FindCharacteristic(ble.NewCharacteristic(ble.MustParse("6E400003-B5A3-F393-E0A9-E50E24DCCA9E")))

		// Initialize joystick 2
		js2, err = joystick.Open(1)
		if err != nil {
			log.Printf("Failed to open joystick 2: %v", err)
		} else {
			joystick2 = true
		}
	}

	ticker := time.NewTicker(40 * time.Millisecond)
	defer ticker.Stop()

	// Create communication channels
	channelBlePlayer1 := make(chan byte)
	channelBlePlayer2 := make(chan byte)

	for {
		select {
		case msg := <-channelBlePlayer1:
			if err := player1.WriteCharacteristic(characteristic1, []byte{msg, 0x0a}, true); err != nil {
				log.Printf("Failed to send data to player 1: %v", err)
			}
		case msg := <-channelBlePlayer2:
			if err := player2.WriteCharacteristic(characteristic2, []byte{msg, 0x0a}, true); err != nil {
				log.Printf("Failed to send data to player 2: %v", err)
			}
		case <-ticker.C:
			if joystick1 {
				go readJoystick(js1, channelBlePlayer1)
			}
			if joystick2 {
				go readJoystick(js2, channelBlePlayer2)
			}
		}
	}
}
