package generator

import (
	"fmt"
	"os"
)

// snapshot records the state of the files a generation step may touch so a
// failure part-way through leaves the project exactly as it was: files that
// did not exist are removed, files that did are restored.
type snapshot struct {
	original map[string][]byte // existing file -> content before the step
	created  []string          // files that did not exist before the step
}

// newSnapshot captures the files the step will modify and notes which of the
// files it will create are new. Paths that do not exist yet in modify are
// ignored; the step itself reports that error.
func newSnapshot(modify, create []string) (*snapshot, error) {
	s := &snapshot{original: map[string][]byte{}}
	for _, p := range modify {
		content, err := os.ReadFile(p)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}
		s.original[p] = content
	}
	for _, p := range create {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			s.created = append(s.created, p)
		}
	}
	return s, nil
}

// restore undoes the step. Errors are reported together rather than aborting,
// so as much as possible is rolled back.
func (s *snapshot) restore() error {
	var failed []error
	for p, content := range s.original {
		if err := os.WriteFile(p, content, 0644); err != nil {
			failed = append(failed, err)
		}
	}
	for _, p := range s.created {
		if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
			failed = append(failed, err)
		}
	}
	if len(failed) > 0 {
		return fmt.Errorf("rollback incomplete: %v", failed)
	}
	return nil
}

// atomically runs step, rolling back every file in modify/create if it fails.
func atomically(modify, create []string, step func() error) error {
	snap, err := newSnapshot(modify, create)
	if err != nil {
		return err
	}
	if err := step(); err != nil {
		if rbErr := snap.restore(); rbErr != nil {
			return fmt.Errorf("%w (and %v)", err, rbErr)
		}
		return err
	}
	return nil
}

// refuseOverwrite returns an error naming the first path that already exists,
// unless force is set.
func refuseOverwrite(force bool, paths ...string) error {
	if force {
		return nil
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return fmt.Errorf("%s already exists (use --force to overwrite)", p)
		}
	}
	return nil
}
