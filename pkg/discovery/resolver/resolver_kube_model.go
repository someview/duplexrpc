package resolver

type (
	ProxyRequest struct {
		Service         string // 服务名称
		Namespace       string // 命名空间名称
		ResourceVersion string // endpoint版本，用来区分最新状态
		PortName        string // 选择的端口名字
	}

	ProxyResponse struct {
		ResourceVersion string         `json:"ResourceVersion"`
		Endpoints       []EndPointInfo `json:"Endpoints"`
	}

	EndPointInfo struct {
		Metadata map[string]string `json:"Metadata"`
		Url      string            `json:"Url"`
	}
)
