import dbus
from dbus.mainloop.glib import DBusGMainLoop
from gi.repository import GLib
import datetime
import json
import os
import threading
import threading

DBusGMainLoop(set_as_default=True)
bluez_hci0_path = "/org/bluez/hci0/"

def get_permitted_devices():
    permitted_devices_file = "permitted_devices.json"
    if not os.path.exists(permitted_devices_file):
        raise FileNotFoundError("permitted_devices.json does not exist. Please create it from the template file.")
    with open(permitted_devices_file, "r") as file:
        devices = json.load(file)

    for device in devices:
        if "mac_address" not in device or "name" not in device:
            raise ValueError("Each device in permitted_devices.json must have a mac_address and name.")
        device["device_path"] = f"{bluez_hci0_path}dev_{device['mac_address'].replace(':', '_')}"

    return devices

class BluetoothManager:
    def __init__(self, device):
        self.bus = dbus.SystemBus()
        self.device_name = device["name"]
        self.device_address = device["mac_address"]
        self.device_path = device["device_path"]

    def on_properties_changed(self, interface, changed_properties, invalidated_properties, path=None):
        timestamp = datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')
        for property_name, value in changed_properties.items():
            if property_name == "Volume":
                print(f"[{timestamp}] Volume Changed: {value}")
            elif property_name == "Status":
                print(f"[{timestamp}] Playback Status Changed: {value}")
            elif property_name == "Track":
                title = value.get('Title', 'Unknown Title')
                artist = value.get('Artist', 'Unknown Artist')
                print(f"[{timestamp}] Now Playing: {title} by {artist}")

        if "Connected" in changed_properties:
            connected = changed_properties["Connected"]
            status = "connected" if connected else "disconnected"
            print(f"[{timestamp}] {self.device_name} {status}: {self.device_address}")


    def on_device_discovered(self, object_path, interfaces_added):
        if self.device_address in object_path:
            print(f"[{datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')}] {self.device_name} Device detected in range: {self.device_address}")
            GLib.idle_add(self.attempt_connect)

    def attempt_connect(self):
        if not self.is_device_connected():
            print(f"[{datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')}] {self.device_name} Attempting to connect to {self.device_address}")
            try:
                device_proxy = self.bus.get_object("org.bluez", self.device_path)
                device_interface = dbus.Interface(device_proxy, "org.bluez.Device1")
                device_interface.Connect()
                print("Connected successfully.")
            except dbus.exceptions.DBusException as e:
                print(f"Connection failed: {e}")
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
        if not self.is_device_connected():
            self.attempt_connect()

        self.bus.add_signal_receiver(self.on_device_discovered,
                                     dbus_interface="org.freedesktop.DBus.ObjectManager",
                                     signal_name="InterfacesAdded")

        print(f"Started listening for {self.device_name} ({self.device_address}) events")

        loop = GLib.MainLoop()
        try:
            loop.run()
        except KeyboardInterrupt:
            loop.quit()

if __name__ == "__main__":
    devices = get_permitted_devices()

    print("Permitted Devices:")
    for device in devices:
        print(f"Device: {device['name']} ({device['mac_address']}) device_path: {device['device_path']}")

    self.bus.add_signal_receiver(self.on_properties_changed,
                                dbus_interface="org.freedesktop.DBus.Properties",
                                signal_name="PropertiesChanged",
                                path_keyword="path")

    threads = []
    for device in devices:
        bluetooth_manager = BluetoothManager(device)
        thread = threading.Thread(target=bluetooth_manager.main)
        threads.append(thread)
        thread.start()

    for thread in threads:
        thread.join()
