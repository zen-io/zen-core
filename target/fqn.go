package target

import (
	"fmt"
	"regexp"
	"strings"
)

type QualifiedTargetName struct {
	project string
	pkg     string
	name    string
	script  string
}

func (fqn *QualifiedTargetName) Project() string {
	return fqn.project
}

func (fqn *QualifiedTargetName) Package() string {
	return fqn.pkg
}

func (fqn *QualifiedTargetName) Name() string {
	return fqn.name
}

func (fqn *QualifiedTargetName) Script() string {
	return fqn.script
}

func (fqn *QualifiedTargetName) BuildFqn() string {
	return fmt.Sprintf("//%s/%s:%s:build", fqn.project, fqn.pkg, fqn.name)
}

func (fqn *QualifiedTargetName) Fqn() string {
	return fmt.Sprintf("//%s/%s:%s:%s", fqn.project, fqn.pkg, fqn.name, fqn.script)
}

func (fqn *QualifiedTargetName) Qn() string {
	return fmt.Sprintf("//%s/%s:%s", fqn.project, fqn.pkg, fqn.name)
}

func (fqn *QualifiedTargetName) SetDefaultScript(defaultScript string) {
	if fqn.script == "" {
		fqn.script = defaultScript
	}
}

func NewFqnFromStr(stepFqn string) (*QualifiedTargetName, error) {
	re := regexp.MustCompile(`^(?:\/\/([\w\d\_\.\-]+)\/([\w\d\_\.\-\/]+))(?::([\w\d\.\_\-\.]+))?(?::([\w\.\d_\-]+))?$`)
	matches := re.FindStringSubmatch(stepFqn)
	if len(matches) == 0 {
		return nil, fmt.Errorf("%s doesnt match an fqn format", stepFqn)
	}
	if matches[3] == "" {
		matches[3] = "all"
	}

	if matches[4] == "" {
		matches[4] = ""
	}

	fqn := &QualifiedTargetName{
		project: matches[1],
		pkg:     matches[2],
		name:    matches[3],
		script:  matches[4],
	}
	return fqn, nil
}

func NewFqnFromStrWithDefault(stepFqn, defaultScript string) (*QualifiedTargetName, error) {
	fqn, err := NewFqnFromStr(stepFqn)
	if err != nil {
		return nil, err
	}
	fqn.SetDefaultScript(defaultScript)

	return fqn, nil
}

func NewFqnFromParts(proj, pkg, name, script string) *QualifiedTargetName {
	return &QualifiedTargetName{
		project: proj,
		pkg:     pkg,
		name:    name,
		script:  script,
	}
}

func InferFqn(target, proj, pkg, defaultScript string) (fqn *QualifiedTargetName, err error) {
	if strings.HasPrefix(target, ":") {
		target = fmt.Sprintf("//%s/%s%s", proj, pkg, target)
	}

	fqn, err = NewFqnFromStr(target)
	fqn.SetDefaultScript(defaultScript)
	return
}
