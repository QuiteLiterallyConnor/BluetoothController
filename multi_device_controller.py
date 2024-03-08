import json
import os
import argparse
import asyncio
from src.bluetooth_manager import BluetoothManager

def get_devices():
    devices = []
    default_file = 'devices.json'
    parser = argparse.ArgumentParser()
    parser.add_argument('--file', help='Specify the JSON file to open', default=default_file)
    parser.add_argument('-h', '--help', action='help', help='Show this help message and exit')
    args = parser.parse_args()
    file_path = args.file

    if os.path.exists(file_path):
        with open(file_path, 'r') as file:
            devices = json.load(file)
    return devices

async def main():
    devices = get_devices()
    tasks = []
    for device in devices:
        bluetooth_manager = BluetoothManager(device)
        task = asyncio.create_task(bluetooth_manager.main())
        tasks.append(task)
    await asyncio.gather(*tasks)

if __name__ == "__main__":
    asyncio.run(main())