package application

const (
	SLUG_COMMON_FUNCTION = "dzn_commonFunctions"
	SLUG_GEAR            = "dzn_gear"
	SLUG_DYNAI           = "dzn_dynai"
	SLUG_TSFRAMEWORK     = "dzn_tSFramework"
	SLUG_PATH_INPUT      = "pathInput"
	SLUG_BRV             = "dzn_brv"
	SLUG_ROOT            = "root"
	SLUG_CONFIG          = "Config"
	SLUG_CONFIG_DYNAI    = "Dynai"
	SLUG_CONFIG_GEAR     = "Gear"
	SLUG_CONFIG_HELPERS  = "Helpers"

	LOG_TAG_INFO  = "INFO"
	LOG_TAG_WARN  = "WARN"
	LOG_TAG_ERROR = "ERROR"

	APP_STATE_NEW         = "NEW"
	APP_STATE_INSTALL     = "INSTALL"
	APP_STATE_DOWNLOADED  = "DOWNLOADED"
	APP_STATE_DIFF        = "DIFF"        // Нужно разрешить разницу в файлах (через UI)
	APP_STATE_DIFF_VSCODE = "DIFF_VSCODE" // Открыт VSCode для разрешения разницы в файле, блокируем UI
	APP_STATE_DONE        = "DONE"
	APP_STATE_ERROR       = "ERROR"

	DIFF_RESOLUTION_BACKUP = "backup"
	DIFF_RESOLUTION_NEW    = "new"
	DIFF_RESOLUTION_OLD    = "old"
	DIFF_RESOLUTION_VSC    = "vsc"

	TEMP_DIR = "FrameworkLatest"

	HTTPS_PREFIX           = "https://"
	GITHUB_DOWNLOAD_SUFFIX = "archive/refs/heads/%s.zip"
	GITHUB_DEFAULT_BRANCH  = "master"
)

var (
	DefaultOptions = map[string]InstallOption{
		SLUG_COMMON_FUNCTION: {
			Slug:        SLUG_COMMON_FUNCTION,
			Value:       "https://github.com/10Dozen/dzn_commonFunctions/tree/v1.8", // github.com/10Dozen/dzn_commonFunctions",
			Label:       "dzn Common Functions",
			Description: "Набор общих функций, которые используются другими компонентами. Обязателен для установки.",
			Error:       "",
			Enabled:     false,
			Checked:     true,
		},
		SLUG_GEAR: {
			Slug:        SLUG_GEAR,
			Value:       "https://github.com/10Dozen/dzn_gear/tree/v2.11", // github.com/10Dozen/dzn_gear",
			Label:       "dzn Gear",
			Description: "Позволяет настраивать выдачу снаряжения игрокам и ботам.",
			Error:       "",
			Enabled:     true,
			Checked:     true,
		},
		SLUG_DYNAI: {
			Slug:        SLUG_DYNAI,
			Value:       "https://github.com/10Dozen/dzn_dynai/tree/v1.3.3", // github.com/10Dozen/dzn_dynai",
			Label:       "dzn DynAI",
			Description: "Создает AI-противников в выбранных зонах и с заданными настройками. Зависит от dzn_gear.",
			Error:       "",
			Enabled:     true,
			Checked:     true,
		},
		SLUG_TSFRAMEWORK: {
			Slug:        SLUG_TSFRAMEWORK,
			Value:       "https://github.com/10Dozen/dzn_tSFramework/tree/v2.14", // github.com/10Dozen/dzn_tSFramework",
			Label:       "Tactical Shift Framework",
			Description: "Базовый фреймворк для создания миссий. Включает множество модулей обеспечивающих геймплей и тактические возможности. Зависит от dzn_gear и dzn_dynai.",
			Error:       "",
			Enabled:     true,
			Checked:     true,
		},
	}

	componentsOrder = [4]string{
		SLUG_COMMON_FUNCTION,
		SLUG_GEAR,
		SLUG_DYNAI,
		SLUG_TSFRAMEWORK,
	}
)

type InstallOption struct {
	Slug        string
	Value       string
	Label       string
	Description string
	Error       string
	Enabled     bool
	Checked     bool
}

type InstallOptions struct {
	Components    []InstallOption
	PathValue     string
	PathError     string
	BackupEnabled bool
	Verified      bool
}

type LogMessage struct {
	Message string
	Level   string
}

type Diff struct {
	OriginalFile string
	NewFile      string
	Diff         string
	Resolution   string
}
