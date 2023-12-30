package plugin

type ELBTarget struct {
	ListenerArn, TgtGrpArn string
	TgtIds                 []string
}
