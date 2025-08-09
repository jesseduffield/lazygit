package helpers

type SuspendResumeHelper struct {
	c *HelperCommon
}

func NewSuspendResumeHelper(c *HelperCommon) *SuspendResumeHelper {
	return &SuspendResumeHelper{
		c: c,
	}
}

func (s *SuspendResumeHelper) CanSuspendApp() bool {
	return canSuspendApp()
}

func (s *SuspendResumeHelper) SuspendApp() error {
	if !canSuspendApp() {
		return nil
	}

	if err := s.c.Suspend(); err != nil {
		return err
	}

	return sendStopSignal()
}

func (s *SuspendResumeHelper) InstallResumeSignalHandler() {
	installResumeSignalHandler(s.c.Log, s.c.Resume)
}
