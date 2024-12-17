package plugins

type HookType string

const (
	ENTER_RULE HookType = "ENTER_RULE"
	ENTER_LINE HookType = "ENTER_LINE"
	EXIT_LINE  HookType = "EXIT_LINE"
	EXIT_RULE  HookType = "EXIT_RULE"
)

type Plugin interface {
	OnEnterRule(data map[string]interface{}) (map[string]interface{}, error)
	OnEnterLine(data map[string]interface{}) (map[string]interface{}, error)
	OnExitLine(data map[string]interface{}) (map[string]interface{}, error)
	OnExitRule(data map[string]interface{}) (map[string]interface{}, error)
}
