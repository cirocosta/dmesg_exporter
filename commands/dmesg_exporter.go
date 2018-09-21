package commands

var DmesgExporter struct {
	RunOnce runOnce `command:"run-once"`
	Start  start  `command:"start"`
}
