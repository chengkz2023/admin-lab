package request

// DirFilePipelineRunRequest describes one reusable run-once task.
type DirFilePipelineRunRequest struct {
	InputDir     string `json:"inputDir"`
	OutputDir    string `json:"outputDir"`
	ErrorDir     string `json:"errorDir"`
	ArchiveDir   string `json:"archiveDir"`
	FilePattern  string `json:"filePattern"`
	MaxFiles     int    `json:"maxFiles"`
	StableWaitMs int    `json:"stableWaitMs"`
	Processor    string `json:"processor"`
	OutputSuffix string `json:"outputSuffix"`
}
