# OpenHeimer
OpenHeimer is the (unofficial) open source version of the Copenheimer project. It scans the entire internet in order to find public Minecraft servers. The original version is currently closed-source and not available to anyone else, that's why this version exists.

:warning: OpenHeimer is under heavy development. The current version of OpenHeimer (v0.0.1) is not stable and might crash (or break) anytime. Feel free to submit an issue or a pull reuqest!

## Usage
View `openheimer --help` to see a list of commands
```
OpenHeimer v0.0.1

	-h, --help		Display a list of available arguments
	-s, --start		Run OpenHeimer and start scanning IPs
		--goroutines		Maximum number of goroutines (10000)
		--timeout		Maximum number of seconds to wait (5 seconds)
		--ip-address		Specify which IP address to start scanning from
	-v, --verbose		Display everthing that's currently happening
	-q, --query		Query something from the database
		--server		Query a Minecraft server from the database
		--player		Query a Minecraft player from the database
```

