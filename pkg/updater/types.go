package updater

type UpgradeEnvelope struct {
	SchemaVersion int                   `json:"schemaVersion"`
	Command       string                `json:"command"`
	Outcome       string                `json:"outcome"` // "upgraded", "up-to-date", "dry-run", "rolled-back", "failed"
	Plan          *UpgradePlan          `json:"plan,omitempty"`
	Events        []UpgradeEvent        `json:"events"`
	Result        *UpgradeResult        `json:"result,omitempty"`
	Recovery      RecoveryInstructions  `json:"recovery"`
	Error         *UpgradeError         `json:"error,omitempty"`
}

type UpgradePlan struct {
	CurrentVersion string `json:"currentVersion"`
	TargetVersion  string `json:"targetVersion"`
	OS             string `json:"os"`
	Arch           string `json:"arch"`
	AssetName      string `json:"assetName"`
	AssetURL       string `json:"assetUrl"`
	ExecutablePath string `json:"executablePath"`
}

type UpgradeEvent struct {
	Sequence int    `json:"sequence"`
	Phase    string `json:"phase"`  // "plan", "download", "validate", "replace", "verify", "cache", "cleanup"
	Status   string `json:"status"` // "started", "succeeded", "failed", "skipped"
	Message  string `json:"message,omitempty"`
}

type UpgradeResult struct {
	InstalledVersion string `json:"installedVersion"`
	CleanupDeferred  bool   `json:"cleanupDeferred"`
}

type RecoveryInstructions struct {
	TargetVersion      string `json:"targetVersion"`
	TargetSource       string `json:"targetSource"` // "resolved" or "fallback-current"
	InstallCommand     string `json:"installCommand"`
	StopProcessCommand string `json:"stopProcessCommand"`
}

type UpgradeError struct {
	Code    string `json:"code"`
	Phase   string `json:"phase"`
	Message string `json:"message"`
}

type FileOperations interface {
	Rename(from, to string) error
	Remove(path string) error
	MakeExecutable(path string) error
}

type UpgradeOptions struct {
	DryRun           bool
	CurrentVersion   string
	RequestedVersion string
	ExecutablePath   string
	GOOS             string
	GOARCH           string
	APIURL           string
	FileOps          FileOperations
	OnPlan           func(plan UpgradePlan)
	OnEvent          func(event UpgradeEvent)
	VerifyFunc       func(execPath, expectedVersion string) error
	DownloadFunc     func(url string) ([]byte, error)
}
