package models

// Field 定义所有可导出字段
var Field = map[string][]string{
	"asset": {
		"time", "lastScanTime", "host", "ip", "port", "service", "tls", "transport",
		"version", "metadata", "project", "type", "tags", "taskName", "rootDomain",
		"urlPath", "hash", "cdnname", "url", "title", "error", "body", "screenshot",
		"faviconmmh3", "faviconpath", "rawheaders", "jarm", "technologies",
		"statuscode", "contentlength", "cdn", "webcheck", "iconcontent", "domain",
		"webServer",
	},
	"crawler": {
		"url", "method", "body", "project", "taskName", "resultId", "rootDomain",
		"time", "tags",
	},
	"DirScanResult": {
		"url", "status", "msg", "project", "length", "taskName", "rootDomain", "tags",
	},
	"SensitiveResult": {
		"url", "urlid", "sid", "match", "project", "color", "time", "md5",
		"taskName", "rootDomain", "tags", "status",
	},
	"subdomain": {
		"host", "type", "value", "ip", "time", "tags", "project", "taskName",
		"rootDomain",
	},
	"SubdomainTakerResult": {
		"input", "value", "cname", "response", "project", "taskName", "rootDomain",
		"tags",
	},
	"UrlScan": {
		"input", "source", "outputtype", "output", "status", "length", "time",
		"body", "project", "taskName", "resultId", "rootDomain", "tags",
	},
	"vulnerability": {
		"url", "vulnid", "vulname", "matched", "project", "level", "time",
		"request", "response", "taskName", "rootdomain", "tags", "status",
	},
	"httpAsset": {
		"time", "lastScanTime", "tls", "hash", "cdnname", "port", "url", "title",
		"type", "error", "body", "host", "ip", "screenshot", "faviconmmh3",
		"faviconpath", "rawheaders", "jarm", "technologies", "statuscode",
		"contentlength", "cdn", "webcheck", "project", "iconcontent", "domain",
		"taskName", "webServer", "service", "rootDomain", "tags",
	},
	"otherAsset": {
		"time", "lastScanTime", "host", "ip", "port", "service", "tls", "transport",
		"version", "metadata", "project", "type", "tags", "taskName", "rootDomain",
		"urlPath",
	},
	"RootDomain": {
		"domain", "company", "icp", "project", "tags", "time",
	},
	"app": {
		"name", "apk", "bundleID", "category", "company", "description", "icp",
		"project", "tags", "time", "url", "version",
	},
	"mp": {
		"name", "category", "company", "description", "icp", "project", "tags", "time",
	},
}

// IsHTTPAssetField 检查字段是否属于HTTP资产
func IsHTTPAssetField(field string) bool {
	for _, f := range Field["httpAsset"] {
		if f == field {
			return true
		}
	}
	return false
}

// IsOtherAssetField 检查字段是否属于其他资产
func IsOtherAssetField(field string) bool {
	for _, f := range Field["otherAsset"] {
		if f == field {
			return true
		}
	}
	return false
}
