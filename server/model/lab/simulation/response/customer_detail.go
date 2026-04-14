package response

type CustomerDetailSimulation struct {
	Breadcrumb    []string                   `json:"breadcrumb"`
	Status        string                     `json:"status"`
	Customer      CustomerSummary            `json:"customer"`
	Metrics       []CustomerMetric           `json:"metrics"`
	BasicInfo     []KeyValueItem             `json:"basicInfo"`
	SecurityOwner []KeyValueItem             `json:"securityOwner"`
	RoomNodes     []CustomerRoomNode         `json:"roomNodes"`
	IPSegments    []CustomerIPSegment        `json:"ipSegments"`
	Applications  []CustomerApplicationEntry `json:"applications"`
}

type CustomerSummary struct {
	Name         string `json:"name"`
	ShortName    string `json:"shortName"`
	CustomerNo   string `json:"customerNo"`
	UnifiedCode  string `json:"unifiedCode"`
	ServiceType  string `json:"serviceType"`
	ServiceStart string `json:"serviceStart"`
}

type CustomerMetric struct {
	Label string `json:"label"`
	Value string `json:"value"`
	Sub   string `json:"sub"`
}

type KeyValueItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type CustomerRoomNode struct {
	Code        string            `json:"code"`
	Name        string            `json:"name"`
	SubName     string            `json:"subName"`
	Tags        []string          `json:"tags"`
	Racks       []string          `json:"racks"`
	Connections []RoomConnection  `json:"connections"`
	Details     []RoomDetailGroup `json:"details"`
}

type RoomConnection struct {
	Label   string `json:"label"`
	Content string `json:"content"`
}

type RoomDetailGroup struct {
	Label string         `json:"label"`
	Items []KeyValueItem `json:"items"`
}

type CustomerIPSegment struct {
	Source      string       `json:"source"`
	Range       string       `json:"range"`
	Mask        string       `json:"mask"`
	AddressNum  int          `json:"addressNum"`
	Usage       string       `json:"usage"`
	NodeCode    string       `json:"nodeCode"`
	AllocatedAt string       `json:"allocatedAt"`
	Hidden      bool         `json:"hidden"`
	NATMappings []NATMapping `json:"natMappings"`
}

type NATMapping struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type CustomerApplicationEntry struct {
	ServiceType string `json:"serviceType"`
	Permit      string `json:"permit"`
	AccessMode  string `json:"accessMode"`
	Content     string `json:"content"`
	Domain      string `json:"domain"`
}
