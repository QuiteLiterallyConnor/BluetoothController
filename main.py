import dbus
from dbus.mainloop.glib import DBusGMainLoop
from gi.repository import GLib
import datetime

# Initialize the D-Bus main loop
DBusGMainLoop(set_as_default=True)

# Create the system bus connection
bus = dbus.SystemBus()

DEVICE_MAC = "0C:C4:13:12:67:62"  # Pixel 6 MAC address
device_path = f"/org/bluez/hci0/dev_{DEVICE_MAC.replace(':', '_')}"

def on_properties_changed(interface, changed_properties, invalidated_properties, path=None):
    timestamp = datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')

    # Print any property changes
    for property_name, value in changed_properties.items():
        if property_name == "Volume":
            print(f"[{timestamp}] Volume Changed: {value} on path: {path}")
        elif property_name == "Status":
            print(f"[{timestamp}] Playback Status Changed: {value}")
        elif property_name == "Track":
            title = value.get('Title', 'Unknown Title')
            artist = value.get('Artist', 'Unknown Artist')
            print(f"[{timestamp}] Now Playing: {title} by {artist}")

    # Handle connection status changes
    if "Connected" in changed_properties:
        connected = changed_properties["Connected"]
        status = "connected" if connected else "disconnected"
        print(f"[{timestamp}] Device {status}: {path}")

def is_device_connected(device_path):
    """Check if the device is currently connected."""
    try:
        device_proxy = bus.get_object("org.bluez", device_path)
        device_properties = dbus.Interface(device_proxy, dbus.PROPERTIES_IFACE)
        return device_properties.Get("org.bluez.Device1", "Connected")
    except dbus.DBusException as e:
        print(f"Error checking connection status: {e}")
        return False

def attempt_reconnect():
    if not is_device_connected(device_path):
        print(f"[{datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')}] Attempting to reconnect to {DEVICE_MAC}")
        try:
            device_proxy = bus.get_object("org.bluez", device_path)
            device_interface = dbus.Interface(device_proxy, "org.bluez.Device1")
            device_interface.Connect()
            print("Reconnected successfully.")
        except dbus.exceptions.DBusException as e:
            print(f"Reconnection failed: {e}")
            # Schedule another attempt if the first one fails
            return True
    else:
        print(f"[{datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')}] Device {DEVICE_MAC} is already connected.")
    return False  # Stop trying to reconnect

if __name__ == "__main__":
    bus.add_signal_receiver(on_properties_changed,
                            dbus_interface="org.freedesktop.DBus.Properties",
                            signal_name="PropertiesChanged",
                            path_keyword="path")

    print("Listening for device connection changes and property changes...")

    # Schedule the first reconnection check
    GLib.timeout_add_seconds(10, attempt_reconnect)

    loop = GLib.MainLoop()
    try:
        loop.run()
    except KeyboardInterrupt:
        print("Program interrupted by user")
        loop.quit()
