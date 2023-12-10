package plugin

type ELBTarget struct {
	ListenerArns, TgtGrpArns, TgtIds []string
}
