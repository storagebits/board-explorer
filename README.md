# board-explorer

board-explorer is an experimental interactive board project. The aim is to learn mixing electronic components with gaming in mind. Code is provived as is mostly for educational purpose. 

Base components are Raspberry Pi, micro:bit, arcade controllers and an IR frame (optional). 

**WARNING:** this project is at early stage and is subject to change often.

![board-explorer](https://github.com/storagebits/board-explorer/blob/master/images/board-explorer.jpg?raw=true)

## Getting started

### Setting up your micro:bit

board-explorer communicate with micro:bit(s) component(s) via BLE (Bluetooth Low Energy). It sends commands on the BLE UART which are then interpreted localy on the micro:bit. So we need first to flash our micro:bit(s) with the following makecode project :

#### Direct link to makecode project

This is the easier way to get started, you can open the makecode project with following link and flash your micro:bit :
https://makecode.microbit.org/_5up1Y89u5J5X

#### Repository of the makecode project

If you need more infos you can check the corresponding repository here :
https://github.com/storagebits/board-explorer-makecode

### Setting up your USB joystick

In the best way it should just work plug'n'play ! just plug it before launching board-explorer.

## Build board-explorer

In the bin folder you'll find a binary compiled on a Raspberry Pi with GOARM=5 build option.

If you're running on a Rapsberry Pi you can run this one (see Run paragraph below) or build your own with following command :

```console
foo@bar:~/board-explorer $ go build board-explorer.go -o bin/board-explorer
```

## Run board-explorer

Now it's time to powerup your flashed micro:bit(s) and launch board-explorer :

```console
foo@bar:~/board-explorer/bin $ sudo ./board-explorer -player1 microbit1name [-player2 microbit2name]
```

_note:_ You can retrieve your micro:bit name in your Android or IOS micro:bit app

_note:_ to avoid sudo you could set following setcap rights to board-explorer binary :
```console
foo@bar:~$ sudo setcap 'cap_net_raw,cap_net_admin=eip' ./board-explorer
```
