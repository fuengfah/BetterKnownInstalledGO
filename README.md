# Better Known Installed GO

A basic package list spoofer

<h1>WARNING</h1>

<h2>This program can break your package manager. Use it carefully.</h2>

## Credits

- [Tesla](https://github.com/0x11DFE) for [BetterKnownInstalled idea](https://github.com/Pixel-Props/BetterKnownInstalled). This project is derived from his own work. I just added packages.list patch and removed all JVM requirements by writing an ABX encoder/decoder from scratch.

### How to compile

`GOOS=android GOARCH=arm64 go build -o main.test main.go`
