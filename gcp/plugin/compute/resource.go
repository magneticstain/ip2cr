package compute

type ComputeResource struct {
	Id, Name, Status string
	PublicIPv4Addrs  []string
	PublicIPv6Addrs  []string
}
