package turbine

import "testing"

func TestLoader(t *testing.T) {
	t.Parallel()

	check := func(t *testing.T, p *Package, err error, found bool) {
		t.Helper()
		if err != nil {
			t.Error("error", err)
		}
		if found {
			if p == nil {
				t.Error("package is nil")
			} else {
				if p.Files == nil {
					t.Error("no file nodes")
				}
				if p.Positions == nil {
					t.Error("no file positions")
				}
			}
		} else {
			if p != nil {
				t.Error("found package that should not exist", p)
			}
		}
	}

	p, err := LoadPackage("time")
	check(t, p, err, true)

	p, err = LoadTestPackage("time")
	check(t, p, err, true)
	p, err = LoadTestPackage("encoding")
	check(t, p, err, false)

	p, err = LoadExternalTestPackage("time")
	check(t, p, err, true)
	p, err = LoadExternalTestPackage("encoding")
	check(t, p, err, false)
}
