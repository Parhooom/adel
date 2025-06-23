package judge

type LanguageConfig struct {
	Language          string
	Image             string
	PreCompileCommand []string
	CompileCommand    []string
	RunCommand        []string
	FileName          string
	IsCompiled        bool
}

func GetLanguageConfigs() map[string]*LanguageConfig {
	return map[string]*LanguageConfig{
		"c": {
			Language:       "c",
			Image:          "gcc:15.1-bookworm",
			CompileCommand: []string{"gcc", "-o", "main", "main.c"},
			RunCommand:     []string{"./main"},
			FileName:       "main.c",
			IsCompiled:     true,
		},
		"go": {
			Language:          "go",
			Image:             "golang:1.24.4-bookworm",
			PreCompileCommand: []string{"go", "mod", "init", "adel"},
			CompileCommand:    []string{"go", "build", "-o", "main", "main.go"},
			RunCommand:        []string{"./main"},
			FileName:          "main.go",
			IsCompiled:        true,
		},
		"python": {
			Language:   "python",
			Image:      "python:3.12-slim",
			RunCommand: []string{"python3", "main.py"},
			FileName:   "main.py",
			IsCompiled: false,
		},
	}
}

func GetLanguageConfig(language string) *LanguageConfig {
	languageConfig, ok := GetLanguageConfigs()[language]
	if ok {
		return languageConfig
	}

	return nil
}

func GetSupportedLanguages() []string {
	configs := GetLanguageConfigs()
	languages := make([]string, 0, len(configs))

	for lang := range configs {
		languages = append(languages, lang)
	}

	return languages
}
