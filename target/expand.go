package target

import (
	"fmt"
	"path/filepath"
)

func (target *Target) ExpandTools(getRefOuts func(ref string) (map[string]string, error)) error {
	tools := map[string]string{}

	for toolName, tool := range target.Tools {
		if IsTargetReference(tool) {
			outs, err := getRefOuts(tool)
			if err != nil {
				return err
			}

			for _, v := range outs {
				tools[toolName] = v
				break
			}
		} else {
			tools[toolName] = tool
		}

		target.Env["PATH"] = fmt.Sprintf("%s:%s", filepath.Dir(tools[toolName]), target.Env["PATH"])
		target.Env["ZEN_TOOL_"+toolName] = tools[toolName]
	}

	target.Tools = tools

	return nil
}
