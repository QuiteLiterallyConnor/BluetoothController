import dbus
from dbus.mainloop.glib import DBusGMainLoop
from gi.repository import GLib
import datetime

DBusGMainLoop(set_as_default=True)
bus = dbus.SystemBus()
bluez_hci0_path = "/org/bluez/hci0/"

class BluetoothManager:
    def __init__(self, device):
        self.bus = dbus.SystemBus()
        self.device_name = device["name"]
        self.device_address = device["mac_address"]
        self.device_path = f"{bluez_hci0_path}dev_{self.device_address.replace(':', '_')}"

    def on_properties_changed(self, interface, changed_properties, invalidated_properties, path=None):
        timestamp = datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')

        for property_name, value in changed_properties.items():
            if property_name == "Volume":
                print(f"[{timestamp}] {self.device_name} Volume Changed: {value} on path: {path}")
            elif property_name == "Status":
                print(f"[{timestamp}] {self.device_name} Playback Status Changed: {value}")
            elif property_name == "Track":
                title = value.get('Title', 'Unknown Title')
                artist = value.get('Artist', 'Unknown Artist')
                print(f"[{timestamp}] {self.device_name} Now Playing: {title} by {artist}")

        if "Connected" in changed_properties:
            connected = changed_properties["Connected"]
            status = "connected" if connected else "disconnected"
            print(f"[{timestamp}] {self.device_name} {status}: {path}")

    def on_device_discovered(self, object_path, interfaces_added):
        if self.device_address in object_path:
            print(f"[{datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')}] {self.device_name} Device detected in range: {self.device_address}")
            GLib.idle_add(self.attempt_reconnect)

    def attempt_reconnect(self):
        if not self.is_device_connected():
            print(f"[{datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')}] {self.device_name} Attempting to reconnect to {self.device_address}")
            try:
                device_proxy = self.bus.get_object("org.bluez", self.device_path)
                device_interface = dbus.Interface(device_proxy, "org.bluez.Device1")
                device_interface.Connect()
                print("Reconnected successfully.")
            except dbus.exceptions.DBusException as e:
                print(f"Reconnection failed: {e}")
        else:
                print(f"[{datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')}] {self.device_name} ({self.device_address}) is already connected.")
        return False

    def is_device_connected(self):
        try:
            device_proxy = self.bus.get_object("org.bluez", self.device_path)
            device_properties = dbus.Interface(device_proxy, dbus.PROPERTIES_IFACE)
            return device_properties.Get("org.bluez.Device1", "Connected")
        except dbus.DBusException as e:
            print(f"Error checking device connection status: {e}")
            return False

    def main(self):
        self.bus.add_signal_receiver(self.on_properties_changed,
                                     dbus_interface="org.freedesktop.DBus.Properties",
                                     signal_name="PropertiesChanged",
                                     path_keyword="path")

        self.bus.add_signal_receiver(self.on_device_discovered,
                                     dbus_interface="org.freedesktop.DBus.ObjectManager",
                                     signal_name="InterfacesAdded")

        print(f"Started listening for {self.device_name} ({self.device_address}) events")

        loop = GLib.MainLoop()
        try:
            loop.run()
        except KeyboardInterrupt:
            loop.quit()