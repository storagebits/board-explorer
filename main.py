def on_bluetooth_connected():
    pass
bluetooth.on_bluetooth_connected(on_bluetooth_connected)

basic.show_leds("""
    # . . # #
    # . . # #
    # # # . .
    # . # . .
    # # # . .
    """)
bluetooth.start_accelerometer_service()
bluetooth.start_button_service()
bluetooth.start_led_service()
bluetooth.start_temperature_service()