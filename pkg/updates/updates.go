package updates

// Update checks for updates and does updates
type Update struct {
	LastChecked string
}

// Updater implements the check and update methods
type Updater interface {
	Check()
	Update()
}

// NewUpdater creates a new updater
func NewUpdater() *Update {

	update := &Update{
		LastChecked: "today",
	}
	return update
}

// Check checks if there is an available update
func (u *Update) Check() {

}
