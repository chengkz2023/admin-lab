package response

type ReliableUploadProfile struct {
	Title            string                        `json:"title"`
	Classification   string                        `json:"classification"`
	Summary          string                        `json:"summary"`
	PackagePath      string                        `json:"packagePath"`
	DocumentPath     string                        `json:"documentPath"`
	CapabilityPoints []string                      `json:"capabilityPoints"`
	TaskModels       []ReliableUploadTaskModel     `json:"taskModels"`
	EngineEntries    []ReliableUploadEntry         `json:"engineEntries"`
	Interfaces       []ReliableUploadInterface     `json:"interfaces"`
	IntegrationSteps []string                      `json:"integrationSteps"`
	MigrationNotes   []string                      `json:"migrationNotes"`
	ExampleConfigs   []ReliableUploadExampleConfig `json:"exampleConfigs"`
	IncludedFiles    []ReliableUploadIncludedFile  `json:"includedFiles"`
	CodeSnippets     []ReliableUploadCodeSnippet   `json:"codeSnippets"`
}

type ReliableUploadTaskModel struct {
	Key         string   `json:"key"`
	Name        string   `json:"name"`
	Scene       string   `json:"scene"`
	TriggerMode string   `json:"triggerMode"`
	Highlights  []string `json:"highlights"`
}

type ReliableUploadEntry struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ReliableUploadInterface struct {
	Name      string   `json:"name"`
	Signature string   `json:"signature"`
	Notes     []string `json:"notes"`
}

type ReliableUploadExampleConfig struct {
	TaskCode        string `json:"taskCode"`
	TaskType        string `json:"taskType"`
	IntervalMinutes int    `json:"intervalMinutes,omitempty"`
	DelaySeconds    int    `json:"delaySeconds,omitempty"`
	BatchSize       int    `json:"batchSize"`
	MaxRetry        int    `json:"maxRetry"`
	SFTPSubdir      string `json:"sftpSubdir"`
	FilePrefix      string `json:"filePrefix"`
}

type ReliableUploadIncludedFile struct {
	Path string `json:"path"`
	Role string `json:"role"`
}

type ReliableUploadCodeSnippet struct {
	Title    string `json:"title"`
	Language string `json:"language"`
	Code     string `json:"code"`
}
