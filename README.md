# OpenHeimer
[![Discord](https://img.shields.io/discord/1100991922018979850?color=%235865F2&label=discord&logo=discord&style=for-the-badge)](https://discord.com/invite/55FxeEd6dx)

**:warning: openheimer is undergoing a major rewrite in Rust (`rust` branch)**

OpenHeimer is the (unofficial) open source version of the Copenheimer project. It scans all the IPv4 addresses in order to find Minecraft servers exposed to the internet. The original version is currently closed-source and not available to anyone else, that's why this version exists.

:rocket: Capable of scanning up to <b>16,000 IP addresses per second</b>, finishing everything in just 3 days! :rocket:

## Installation

In order to use the current version of openheimer, you must compile from source.

### Requirements
- git (or download the tarball/zip)
- [go](https://go.dev/dl/)

In your terminal,

```shell
$ # Not necessary with tarball/zip
$ git clone https://github.com/ErrorNoInternet/openheimer.git

$ cd openheimer
$ go build
```

You should now have an executable in the same directory.

## Usage
Run `openheimer -help` to see a list of commands
```
Usage of ./openheimer:
  -database string
        The database to store the results in (default "openheimer.db")
  -ipFile string
        The file to extract IP addresses from
  -logFile string
        The file to store the logs in (default "openheimer.log")
  -maxPingWorkers int
        The maximum amount of workers to ping IPs (default 4000)
  -maxScanWorkers int
        The maximum amount of workers to scan IPs (default 1000)
  -startingIP string
        The IP address to start scanning from (default "1.0.0.0")
  -timeout int
        The amount of seconds to wait before timing out (default 5)
  -verbose
        Display everything that's happening
  -version
        Display the current version of OpenHeimer
```
Note: `maxPingWorkers` are the amount of workers to use to check for open ports. `maxScanWorkers` are the amount of workers to use to check if an IP has a valid Minecraft server. Normally you would increase `maxPingWorkers` to ping IP addresses faster, but when you're scanning IP addresses from a file (produced by masscan for example), you should increase `maxScanWorkers`.

<sub>If you would like to modify or use this repository (including its code) in your own project, please be sure to credit!</sub>
