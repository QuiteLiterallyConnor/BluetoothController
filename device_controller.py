import argparse
from src.bluetooth_manager import BluetoothManager

def get_device():
    parser = argparse.ArgumentParser()
    parser.add_argument("--name", help="Name of the device")
    parser.add_argument("--mac_address", help="MAC address of the device")
    args = parser.parse_args()

    if args.help:
        parser.print_help()
        exit()

    device = {
        "name": args.name,
        "mac_address": args.mac_address
    }

    return device

def main():
    device = get_device()
    bluetooth_manager = BluetoothManager(device)
    bluetooth_manager.main()

if __name__ == "__main__":
    main()