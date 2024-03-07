import dbus
from dbus.mainloop.glib import DBusGMainLoop
from gi.repository import GLib
import datetime

# Initialize the D-Bus main loop
DBusGMainLoop(set_as_default=True)

# Create the system bus connection
bus = dbus.SystemBus()

DEVICE_MAC = "0C:C4:13:12:67:62"  # Your Pixel 6 MAC address
device_path = f"/org/bluez/hci0/dev_{DEVICE_MAC.replace(':', '_')}"

def on_properties_changed(interface, changed_properties, invalidated_properties, path=None):
    timestamp = datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')
    
    # Ensure we're handling properties changes for the device of interest
    if path and path.endswith(device_path.split('/')[-1]):
        # Process and print relevant property changes
        for property_name, value in changed_properties.items():
            if property_name in ["Volume", "Status", "Track"]:
                print(f"[{timestamp}] {property_name} Changed: {value} on path: {path}")

def on_device_discovered(object_path, interfaces_added):
    global device_path
    if device_path in object_path:
        print(f"[{datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')}] Device detected in range: {DEVICE_MAC}")
        GLib.idle_add(attempt_reconnect)

def is_device_connected():
    try:
        device_proxy = bus.get_object("org.bluez", device_path)
        device_properties = dbus.Interface(device_proxy, dbus.PROPERTIES_IFACE)
        return device_properties.Get("org.bluez.Device1", "Connected")
    except dbus.DBusException as e:
        print(f"Error checking device connection status: {e}")
        return False

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
    return False  # Ensure this function does not repeat

if __name__ == "__main__":
    bus.add_signal_receiver(on_properties_changed,
                            dbus_interface="org.freedesktop.DBus.Properties",
                            signal_name="PropertiesChanged",
                            path_keyword="path")
                            
    bus.add_signal_receiver(on_device_discovered,
                            dbus_interface="org.freedesktop.DBus.ObjectManager",
                            signal_name="InterfacesAdded")

    print("Listening for device connection changes, property changes, and device discovery...")

    loop = GLib.MainLoop()
    try:
        loop.run()
    except KeyboardInterrupt:
        print("Program interrupted by user")
        loop.quit()
