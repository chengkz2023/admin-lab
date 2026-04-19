package response

type DirFilePipelineProfile struct {
	Title            string                           `json:"title"`
	Classification   string                           `json:"classification"`
	Summary          string                           `json:"summary"`
	Highlights       []string                         `json:"highlights"`
	QuickSteps       []string                         `json:"quickSteps"`
	Processors       []DirFilePipelineProcessorOption `json:"processors"`
	DefaultConfig    DirFilePipelineConfigSample      `json:"defaultConfig"`
	IntegrationNotes []string                         `json:"integrationNotes"`
}

type DirFilePipelineProcessorOption struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DirFilePipelineConfigSample struct {
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

type DirFilePipelineRunResult struct {
	EffectiveConfig DirFilePipelineConfigSample  `json:"effectiveConfig"`
	Scanned         int                          `json:"scanned"`
	Skipped         int                          `json:"skipped"`
	Processed       int                          `json:"processed"`
	Failed          int                          `json:"failed"`
	OutputFiles     []string                     `json:"outputFiles"`
	ArchivedFiles   []string                     `json:"archivedFiles"`
	SkippedFiles    []string                     `json:"skippedFiles"`
	FailureItems    []DirFilePipelineFailureItem `json:"failureItems"`
}

type DirFilePipelineFailureItem struct {
	File      string `json:"file"`
	Stage     string `json:"stage"`
	Reason    string `json:"reason"`
	ErrorFile string `json:"errorFile"`
}
