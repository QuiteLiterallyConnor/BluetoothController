import dbus

def get_media_player_object():
    bus = dbus.SystemBus()
    obj = bus.get_object('org.bluez', '/')
    mgr = dbus.Interface(obj, 'org.freedesktop.DBus.ObjectManager')
    for path, interfaces in mgr.GetManagedObjects().items():
        if 'org.bluez.MediaPlayer1' in interfaces:
            return dbus.Interface(bus.get_object('org.bluez', path), 'org.bluez.MediaPlayer1')
    return None

def play_pause():
    player = get_media_player_object()
    if player is None:
        print("Media player not found. Make sure your device is connected.")
        return
    player.PlayPause()

def next_track():
    player = get_media_player_object()
    if player is None:
        print("Media player not found. Make sure your device is connected.")
        return
    player.Next()

def previous_track():
    player = get_media_player_object()
    if player is None:
        print("Media player not found. Make sure your device is connected.")
        return
    player.Previous()

# Add interactive user input to call functions
if __name__ == "__main__":
    while True:
        print("AVRCP Controller:")
        print("1. Play/Pause")
        print("2. Next Track")
        print("3. Previous Track")
        print("4. Exit")
        choice = input("Select an action (1-4): ")

        if choice == '1':
            play_pause()
        elif choice == '2':
            next_track()
        elif choice == '3':
            previous_track()
        elif choice == '4':
            print("Exiting...")
            break
        else:
            print("Invalid choice. Please select a number between 1-4.")
