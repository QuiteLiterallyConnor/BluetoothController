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

def on_device_discovered(object_path, interfaces_added):
    if DEVICE_MAC in object_path:
        print(f"[{datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')}] Device detected in range: {DEVICE_MAC}")
        # Defer reconnection to avoid direct call from signal
        GLib.idle_add(attempt_reconnect)

def attempt_reconnect():
    if not is_device_connected():
        print(f"[{datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')}] Attempting to reconnect to {DEVICE_MAC}")
        try:
            device_proxy = bus.get_object("org.bluez", device_path)
            device_interface = dbus.Interface(device_proxy, "org.bluez.Device1")
            device_interface.Connect()
            print("Reconnected successfully.")
        except dbus.exceptions.DBusException as e:
            print(f"Reconnection failed: {e}")
    else:
        print(f"[{datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')}] Device {DEVICE_MAC} is already connected.")
    return False  # Stop further idle calls

def is_device_connected():
    try:
        device_proxy = bus.get_object("org.bluez", device_path)
        device_properties = dbus.Interface(device_proxy, dbus.PROPERTIES_IFACE)
        return device_properties.Get("org.bluez.Device1", "Connected")
    except dbus.DBusException as e:
        print(f"Error checking device connection status: {e}")
        return False

if __name__ == "__main__":
    # Listen for property changes globally but apply logic specifically
    bus.add_signal_receiver(on_properties_changed,
                            dbus_interface="org.freedesktop.DBus.Properties",
                            signal_name="PropertiesChanged",
                            path_keyword="path")

    # Listen for new devices being discovered
    bus.add_signal_receiver(on_device_discovered,
                            dbus_interface="org.freedesktop.DBus.ObjectManager",
                            signal_name="InterfacesAdded")

    print("Listening for device connection changes, property changes, and device discovery...")

    loop = GLib.MainLoop()
    try:
        loop.run()
    except KeyboardInterrupt:
        loop.quit()
