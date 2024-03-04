package awsipprefix

type GenericAWSPrefix struct {
	IPRange            string
	Region             string
	Service            string
	NetworkBorderGroup string
}

type AwsIpv4Prefix struct {
	IPPrefix           string `json:"ip_prefix"`
	Region             string `json:"region"`
	Service            string `json:"service"`
	NetworkBorderGroup string `json:"network_border_group"`
}

type AwsIpv6Prefix struct {
	IPv6Prefix         string `json:"ipv6_prefix"`
	Region             string `json:"region"`
	Service            string `json:"service"`
	NetworkBorderGroup string `json:"network_border_group"`
}

type RawAwsIPRangeJSON struct {
	SyncToken    string          `json:"syncToken"`
	CreateDate   string          `json:"createDate"`
	Prefixes     []AwsIpv4Prefix `json:"prefixes"`
	IPv6Prefixes []AwsIpv6Prefix `json:"ipv6_prefixes"`
}
