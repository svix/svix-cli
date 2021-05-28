package cfg

type Config struct {
	// Spinner *spinner.Spinner
}

func New() *Config {
	return &Config{
		// Spinner: spinner.New(spinner.CharSets[14], 100*time.Millisecond),
	}
}
