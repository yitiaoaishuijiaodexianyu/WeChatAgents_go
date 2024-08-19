package plug_in

import _struct "WeChatAgents_go/struct"

type PlugInFunc func(_struct.Message) (_struct.PlugInResult, error)

func PlugIn() map[string]PlugInFunc {
	var PlugInMap = map[string]PlugInFunc{}
	return PlugInMap
}
