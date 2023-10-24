module github.com/c12s/cockpit

go 1.21

require (
	github.com/c12s/kuiper v1.0.0
	github.com/c12s/magnetar v1.0.0
	github.com/c12s/oort v1.0.0
	github.com/fatih/color v1.15.0
	github.com/jedib0t/go-pretty/v6 v6.4.6
	github.com/spf13/cobra v1.7.0
	google.golang.org/grpc v1.57.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/nats-io/nats.go v1.28.0 // indirect
	github.com/nats-io/nkeys v0.4.4 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/crypto v0.11.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230525234030-28d5490b6b19 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)

replace github.com/c12s/magnetar => ../magnetar

replace github.com/c12s/kuiper => ../kuiper

replace github.com/c12s/oort => ../oort
