#!/usr/bin/python3

from bluezero import microbit
ubit = microbit.Microbit(adapter_addr='DC:A6:32:5D:64:8D',
                         #device_addr='DC:4A:04:CA:C3:6F',
                         device_addr='E2:31:3A:95:93:94',
                         accelerometer_service=True,
                         button_service=True,
                         led_service=True,
                         magnetometer_service=False,
                         pin_service=False,
                         temperature_service=True)
my_text = 'Hello, world'
ubit.connect()

while my_text is not '':
    ubit.text = my_text
    my_text = input('Enter message: ')

ubit.disconnect()
