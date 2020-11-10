#!/usr/bin/python3

import time
import sys
import os

from bluezero import microbit

DEVICE = "DC:4A:04:CA:C3:6F"   # device #1
#DEVICE = "DC:A6:32:5D:64:8D"   # device #2

if len(sys.argv) == 2:
  DEVICE = str(sys.argv[1])

ubit = microbit.Microbit(adapter_addr='DC:A6:32:5D:64:8D',
                         device_addr=DEVICE,
                         accelerometer_service=True,
                         button_service=True,
                         led_service=True,
                         magnetometer_service=False,
                         pin_service=False,
                         temperature_service=True)

my_text = 'Hello, world'
ubit.connect()
ubit.text = "TEST"
#while my_text is not '':
    #ubit.text = my_text
    #my_text = input('Enter message: ')
ubit.disconnect()
