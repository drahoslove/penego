package gui

type Icon rune

const ( //  font awesome corresponding characters
	QuitIcon  = '\uf00d'
	FileIcon  = '\uf15b'
	PlayIcon  = '\uf04b'
	PauseIcon = '\uf04c'
	ResetIcon = '\uf021'
)

func AlwaysIcon(icon Icon) func() Icon {
	return func() Icon {
		return icon
	}
}
