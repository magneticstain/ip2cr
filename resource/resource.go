package resource

type Resource struct {
	Id, RID, AccountID, Name, Status, CloudSvc                   string
	AccountAliases, NetworkMap, PublicIPv4Addrs, PublicIPv6Addrs []string
}
