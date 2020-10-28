#!/usr/bin/python3
# Python script to get data from multiple devices
# Usage:
# python GetDataAll.py

import subprocess

# list of devices
devices = [
  "DC:4A:04:CA:C3:6F",   # device #1
  "DC:A6:32:5D:64:8D"   # device #2
   ]

for x in range(0,len(devices)):
  cmd = "python pybit.py " + devices[x]
  subprocess.call(cmd, shell=True)

print("finished all devices!")