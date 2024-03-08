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
    args = parser.parse_args()
    file_path = args.file

    if os.path.exists(file_path):
        with open(file_path, 'r') as file:
            devices = json.load(file)
    else:
        print(f"File not found: {file_path}")
        exit(1)

    if len(devices) == 0:
        print("No devices found in the file.")
        exit(1)
        
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