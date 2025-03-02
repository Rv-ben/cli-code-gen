package codeEditor

type BaseAction interface {
	ToString() string
	GetType() string
}
