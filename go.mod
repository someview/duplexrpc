module gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git

go 1.22

toolchain go1.22.1

require gitlab.dev.wiqun.com/tl/goserver/chat/l1/tl.gobase.git v1.0.6

require (
	github.com/cloudwego/netpoll v0.6.3
	github.com/golang/protobuf v1.5.0
	github.com/jhue58/latency v0.3.0
	github.com/json-iterator/go v1.1.12
	github.com/orcaman/concurrent-map/v2 v2.0.1
	github.com/panjf2000/ants/v2 v2.10.0
	github.com/stretchr/testify v1.8.2
	go.uber.org/mock v0.4.0
)

require (
	github.com/bytedance/gopkg v0.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sync v0.3.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/cloudwego/netpoll => github.com/someview/netpoll v0.5.2-0.20240819022335-ca3d1998d221
