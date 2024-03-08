# Car Audio System

This is a Python program that manages a car audio system using Bluetooth. It can detect when a permitted device is in range, attempt to reconnect to it if it's not currently connected, and handle changes in the device's properties such as volume, status, and track.

## Requirements

This program requires Python 3 and the following Python packages:

- dbus-python
- PyGObject

You can install these packages using pip:

```bash
pip install -r requirements.txt
```

## Usage 

 You can run this program using the following command:

```python
python main.py --name [DEVICE NAME] --mac_address [MAC ADDRESS]
```


Replace [DEVICE NAME] with the name of the device and [MAC ADDRESS] with the MAC address of the device.

## Configuration

You can configure the permitted devices by modifying the `devices.json` file. Each device in this file should have a `mac_address` and `name`.

Here's an example of what this file might look like:

```json
[
    {
        "mac_address": "00:11:22:33:44:55",
        "name": "My Phone"
    },
    {
        "mac_address": "66:77:88:99:AA:BB",
        "name": "My Tablet"
    }
]
```