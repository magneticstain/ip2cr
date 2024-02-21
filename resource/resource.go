package resource

type Resource struct {
	RID, AccountID, CloudSvc   string
	AccountAliases, NetworkMap []string
}
