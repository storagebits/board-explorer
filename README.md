## board-explorer

board-explorer is an experimental interactive board project. The aim is to learn mixing electronic components with gaming in mind. 

Base components are Raspberry Pi, micro:bit, arcade controllers and an IR frame (optional). 

**WARNING:** this project is at early stage and is subject to change often. I put this here mainly for educational purpose.

![board-explorer](https://github.com/storagebits/board-explorer/blob/master/images/board-explorer.jpg?raw=true)

## Getting started

## Setting up your micro:bit

board-explorer communicate with micro:bit(s) component(s) via BLE (Bluetooth Low Energy). It sends commands on the BLE UART which are then interpreted localy on the micro:bit. So we need first to flash our micro:bit(s) with the following makecode project :

### Direct link to makecode project 
https://makecode.microbit.org/_5up1Y89u5J5X

### Repository of the makecode project
https://github.com/storagebits/board-explorer-makecode

## Build board-explorer

## Run board-explorer
```console
foo@bar:~$ sudo ./board-explorer -player1 microbit1name [-player2 microbit2name]
```

_note:_ to avoid sudo you could set following setcap rights to board-explorer binary :
```console
foo@bar:~$ sudo setcap 'cap_net_raw,cap_net_admin=eip' ./board-explorer
```
