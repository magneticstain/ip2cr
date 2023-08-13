package awsfqdnregexmap

func GetRegexMap() map[string]string {
	return map[string]string{
		"CLOUDFRONT": "^[a-z0-9\\-]+\\.[a-z0-9\\-]+\\.[a-z0-9\\-]+\\.cloudfront\\.net\\.$",
		"EC2":        "^ec2\\-[\\d]{1,3}\\-[\\d]{1,3}\\-[\\d]{1,3}\\-[\\d]{1,3}\\.[a-z0-9\\-]+\\.amazonaws\\.com\\.$",
	}
}
