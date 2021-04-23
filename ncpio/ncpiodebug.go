package ncpio

type (
	// Debuger interface allows implementations to provide to this package any
	// object that implements the methods defined in it.
	Debuger interface {
		Println(v ...interface{})
		Printf(format string, v ...interface{})
	}

	// NoDebuger implements the logger that does not perform any operation
	// by default. This allows us to efficiently discard the unwanted messages.
	NoDebuger struct{}
)

// Println is the library provided NoDebuger's
// implementation of the required interface function()
func (NoDebuger) Println(v ...interface{}) {}

// Printf is the library provided NoDebuger's
// implementation of the required interface function(){}
func (NoDebuger) Printf(format string, v ...interface{}) {}
