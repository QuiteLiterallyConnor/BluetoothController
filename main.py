import dbus
from dbus.mainloop.glib import DBusGMainLoop
from gi.repository import GLib
import datetime

DBusGMainLoop(set_as_default=True)
bus = dbus.SystemBus()

DEVICE_MAC = "0C:C4:13:12:67:62"
device_path = f"/org/bluez/hci0/dev_{DEVICE_MAC.replace(':', '_')}"

def on_properties_changed(interface, changed_properties, invalidated_properties, path=None):
    timestamp = datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')
    # Existing property change handling logic

def is_device_connected(device_path):
    try:
        device_proxy = bus.get_object("org.bluez", device_path)
        device_properties = dbus.Interface(device_proxy, dbus.PROPERTIES_IFACE)
        return device_properties.Get("org.bluez.Device1", "Connected")
    except dbus.DBusException:
        return False

def on_device_discovered(object_path, interfaces_added):
    if device_path in object_path:
        print(f"[{datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')}] Device detected in range: {DEVICE_MAC}")
        if not is_device_connected(device_path):
            attempt_reconnect()

def attempt_reconnect():
    print(f"[{datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')}] Attempting to reconnect to {DEVICE_MAC}")
    try:
        device_proxy = bus.get_object("org.bluez", device_path)
        device_interface = dbus.Interface(device_proxy, "org.bluez.Device1")
        device_interface.Connect()
        print("Reconnected successfully.")
    except dbus.exceptions.DBusException as e:
        print(f"Reconnection failed: {e}")

def main():
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
        loop.quit()

if __name__ == "__main__":
    main()