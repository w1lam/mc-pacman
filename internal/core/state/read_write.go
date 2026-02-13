package state

// Read safely reads state
func (s *State) Read(fn func(*State)) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	fn(s)
}

// Write safley writes to state and saves manifest
func (s *State) Write(fn func(*State) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := fn(s); err != nil {
		return err
	}

	return s.repo.Save(s.manifest)
}
