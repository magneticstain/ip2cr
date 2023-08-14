package awsipprefix

type GenericAWSPrefix struct {
	IpRange            string
	Region             string
	Service            string
	NetworkBorderGroup string
}

type AwsIpv4Prefix struct {
	IpPrefix           string `json:"ip_prefix"`
	Region             string `json:"region"`
	Service            string `json:"service"`
	NetworkBorderGroup string `json:"network_border_group"`
}

type AwsIpv6Prefix struct {
	Ipv6Prefix         string `json:"ipv6_prefix"`
	Region             string `json:"region"`
	Service            string `json:"service"`
	NetworkBorderGroup string `json:"network_border_group"`
}

type RawAwsIpRangeJSON struct {
	SyncToken    string          `json:"syncToken"`
	CreateDate   string          `json:"createDate"`
	Prefixes     []AwsIpv4Prefix `json:"prefixes"`
	IPv6Prefixes []AwsIpv6Prefix `json:"ipv6_prefixes"`
}
