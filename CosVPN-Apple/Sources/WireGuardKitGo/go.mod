module github.com/ochernishov/cosvpn/apple

go 1.23.1

require (
	github.com/ochernishov/cosvpn v0.0.0
	golang.org/x/sys v0.32.0
)

require (
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/net v0.39.0 // indirect
	golang.zx2c4.com/wintun v0.0.0-20230126152724-0fa3db229ce2 // indirect
)

replace github.com/ochernishov/cosvpn => ../../../CosVPN-Go
