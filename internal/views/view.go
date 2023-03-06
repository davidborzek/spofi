package views

// View represents a view of the application
// which will be used to show in rofi.
type View interface {
	Show(payload ...interface{})
	SetParent(parent View)
}
