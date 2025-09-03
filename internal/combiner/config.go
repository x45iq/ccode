package combiner

type Config struct {
	RootDir    string
	Output     string
	Force      bool
	DryRun     bool
	StripEmpty bool
}
