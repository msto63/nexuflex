module github.com/nexuflex/nexuflex-client

go 1.20

require (
	github.com/gdamore/tcell/v2 v2.8.1
	github.com/golang/protobuf v1.5.4
	github.com/nexuflex/shared v0.0.0-00010101000000-000000000000
	github.com/rivo/tview v0.0.0-20241227133733-17b7edb88c57
	google.golang.org/grpc v1.71.0
)

require (
	github.com/gdamore/encoding v1.0.1 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/term v0.30.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250311190419-81fb87f6b8bf // indirect
	google.golang.org/protobuf v1.36.5 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
)

replace github.com/nexuflex/shared => ../shared 
