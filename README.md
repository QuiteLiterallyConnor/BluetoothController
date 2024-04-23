# Bluetooth Manager

This package provides functionality to VERY SIMPLY control and scan Bluetooth devices. It includes the following main components:

- `BluetoothController`: This component is responsible for controlling Bluetooth devices. It provides functions to connect, disconnect, and manage Bluetooth devices. It can also send AVRCP commands to the device, such as play, pause, next, previous, etc

- `BluetoothScanner`: This component is responsible for scanning for nearby Bluetooth devices. It provides functions to start and stop scanning, and to retrieve the list of discovered devices.

## Usage

To use this package, you need to import it in your Go code:

```go
import "github.com/QuiteLiterallyConnor/BluetoothManager"
```

Then, you can create a new BluetoothController or BluetoothScanner:

```go
controller, err := bluetoothmanager.NewBluetoothController()
if err != nil {
    // handle error
}

scanner, err := bluetoothmanager.NewBluetoothScanner()
if err != nil {
    // handle error
}
```

You can then use the methods provided by these objects to control and scan Bluetooth devices.

## Examples
Examples of how to use this package can be found in the examples directory. These examples demonstrate how to use the BluetoothController and BluetoothScanner to perform common tasks.

## Testing
This package includes unit tests. To run the tests, navigate to the package directory and run "go test"

## Contributing
Contributions to this package are welcome. Please submit a pull request or open an issue if you have any improvements or suggestions.

## License
This package is licensed under the MIT License. See the LICENSE file for more details.
