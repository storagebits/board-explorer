bluetooth.onBluetoothConnected(function () {
    basic.showLeds(`
        . # # # .
        # . . . .
        # . . . .
        # . . . .
        . # # # .
        `)
})
bluetooth.onBluetoothDisconnected(function () {
    basic.showLeds(`
        # # # . .
        # . . # .
        # . . # .
        # . . # .
        # # # . .
        `)
})
basic.showLeds(`
    # . . # #
    # . . # #
    # # # . .
    # . # . .
    # # # . .
    `)
bluetooth.startAccelerometerService()
bluetooth.startButtonService()
bluetooth.startLEDService()
bluetooth.startTemperatureService()
