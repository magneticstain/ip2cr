package awsfqdnregexmap

func GetRegexMap() map[string]string {
	return map[string]string{
		"CLOUDFRONT": "^[a-z0-9\\-]+\\.[a-z0-9\\-]+\\.[a-z0-9\\-]+\\.cloudfront\\.net\\.$",                            // EX: server-65-8-191-186.bos50.r.cloudfront.net.
		"EC2":        "^ec2\\-[\\d]{1,3}\\-[\\d]{1,3}\\-[\\d]{1,3}\\-[\\d]{1,3}\\.[a-z0-9\\-]+\\.amazonaws\\.com\\.$", // EX: ec2-35-170-192-9.compute-1.amazonaws.com.
	}
}
