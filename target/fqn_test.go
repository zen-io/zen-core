package target

import "testing"

func TestFqn(t *testing.T) {
	fqn, err := NewFqnFromStr("//project/path/to/pkg:name:script")
	if err != nil {
		t.Error(err)
		return
	}

	compareFqn(fqn, t)
}

func TestInferFqnNameAndScript(t *testing.T) {
	fqn, err := InferFqn(":name:script", "project", "path/to/pkg", "error")
	if err != nil {
		t.Error(err)
		return
	}

	compareFqn(fqn, t)
}

func TestInferFqnName(t *testing.T) {
	fqn, err := InferFqn(":name", "project", "path/to/pkg", "script")
	if err != nil {
		t.Error(err)
		return
	}

	compareFqn(fqn, t)
}

func TestInferFqnFullTarget(t *testing.T) {
	fqn, err := InferFqn("//project/path/to/pkg:name:script", "error", "error", "error")
	if err != nil {
		t.Error(err)
		return
	}

	compareFqn(fqn, t)
}
func TestInferFqnFullTargetNoScript(t *testing.T) {
	fqn, err := InferFqn("//project/path/to/pkg:name", "error", "error", "script")
	if err != nil {
		t.Error(err)
		return
	}

	compareFqn(fqn, t)
}

func compareFqn(fqn *QualifiedTargetName, t *testing.T) {
	compareQn(fqn, t)

	if fqn.Script() != "script" {
		t.Errorf("Script not correct: name vs %s\n", fqn.Script())
	}
}

func compareQn(fqn *QualifiedTargetName, t *testing.T) {
	if fqn.Project() != "project" {
		t.Errorf("Project not correct: project vs %s\n", fqn.Project())
	}
	if fqn.Package() != "path/to/pkg" {
		t.Errorf("Package not correct: path/to/pkg vs %s\n", fqn.Package())
	}
	if fqn.Name() != "name" {
		t.Errorf("Name not correct: name vs %s\n", fqn.Name())
	}
}
