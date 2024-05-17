module github.com/c12s/cockpit

go 1.21

require (
	github.com/spf13/cobra v1.7.0
	golang.org/x/term v0.20.0
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.20.0 // indirect
)

replace github.com/c12s/magnetar => ../magnetar

replace github.com/c12s/kuiper => ../kuiper

replace github.com/c12s/oort => ../oort
