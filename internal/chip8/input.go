package chip8

type Keypad interface {
	GetKeyStates() [16]bool
}
