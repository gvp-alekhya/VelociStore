package core

type RespCmd struct {
	Cmd  string
	Args []string
}

type RespCmds []*RespCmd
