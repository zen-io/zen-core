package target

import (
	"fmt"
	"os/exec"
	"strings"

	out "github.com/tiagoposse/go-tasklist-out"
	"github.com/zen-io/zen-core/utils"
)

type Target struct {
	Labels      []string
	Hashes      []string
	Tools       map[string]string
	Srcs        map[string][]string
	Outs        []string
	Env         map[string]string

	shouldInterpolate bool
	Cwd    string
	
	*QualifiedTargetName
	out.TaskLogger
}

func (t *Target) Exec(command []string, errorMsg string) error {
	t.Debugln(strings.Join(command, " "))
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Dir = t.Cwd
	cmd.Env = t.GetEnvironmentVariablesList()

	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%s: %w (%s)", errorMsg, err, out)
	}
	return nil
}

func (t *Target) Copy(from, to string, customVars ...map[string]string) error {
	t.Traceln("copying from %s to %s (interpol: %v)", from, to, t.shouldInterpolate)
	if t.shouldInterpolate {
		if err := utils.CopyWithInterpolate(from, to, append([]map[string]string{t.Env}, customVars...)...); err != nil {
			return fmt.Errorf("copying from %s to %s: %w", from, to, err)
		}
	} else if from != to {
		return utils.CopyFile(from, to)
	}

	return nil
}
