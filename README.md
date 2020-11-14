## board-explorer

board-explorer is an experimental interactive board project. The aim is to learn mixing electronic components with gaming in mind. 

Base components are Raspberry Pi, micro:bit, arcade controllers and an IR frame (optional). 

**WARNING:** this project is at early stage and is subject to change often. I put this here mainly for educational purpose.

![board-explorer](https://github.com/storagebits/board-explorer/blob/master/images/board-explorer.jpg?raw=true)

## Getting started

```console
foo@bar:~$ sudo ./board-explorer -player1 microbit1name [-player2 microbit2name]
```

_note:_ to avoid sudo you could set following setcap rights to board-explorer binary :
```console
foo@bar:~$ sudo setcap 'cap_net_raw,cap_net_admin=eip' ./board-explorer
```
