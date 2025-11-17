package models

// Template 扫描模板
type ScanTemplate struct {
	ID                  string     `bson:"_id,omitempty" json:"id"`
	Ignore              string     `bson:"ignore" json:"ignore"`
	Target              string     `bson:"target" json:"target"`
	Type                string     `bson:"type" json:"type"`
	Duplicates          string     `bson:"duplicates" json:"duplicates"`
	IsStart             bool       `bson:"isStart" json:"isStart"`
	TaskName            string     `bson:"TaskName" json:"TaskName"`
	TargetHandler       []string   `bson:"TargetHandler" json:"TargetHandler"`
	Parameters          Parameters `bson:"Parameters" json:"Parameters"`
	SubdomainScan       []string   `bson:"SubdomainScan" json:"SubdomainScan"`
	SubdomainSecurity   []string   `bson:"SubdomainSecurity" json:"SubdomainSecurity"`
	PortScanPreparation []string   `bson:"PortScanPreparation" json:"PortScanPreparation"`
	PortScan            []string   `bson:"PortScan" json:"PortScan"`
	PortFingerprint     []string   `bson:"PortFingerprint" json:"PortFingerprint"`
	AssetMapping        []string   `bson:"AssetMapping" json:"AssetMapping"`
	AssetHandle         []string   `bson:"AssetHandle" json:"AssetHandle"`
	URLScan             []string   `bson:"URLScan" json:"URLScan"`
	WebCrawler          []string   `bson:"WebCrawler" json:"WebCrawler"`
	URLSecurity         []string   `bson:"URLSecurity" json:"URLSecurity"`
	DirScan             []string   `bson:"DirScan" json:"DirScan"`
	VulnerabilityScan   []string   `bson:"VulnerabilityScan" json:"VulnerabilityScan"`
	Name                string     `bson:"name" json:"name"`
	VulList             []string   `bson:"vullist" json:"vullist"`
	PassiveScan         []string   `bson:"PassiveScan" json:"PassiveScan"`
}

type Parameters struct {
	TargetHandler       map[string]string `bson:"TargetHandler" json:"TargetHandler"`
	SubdomainScan       map[string]string `bson:"SubdomainScan" json:"SubdomainScan"`
	SubdomainSecurity   map[string]string `bson:"SubdomainSecurity" json:"SubdomainSecurity"`
	PortScanPreparation map[string]string `bson:"PortScanPreparation" json:"PortScanPreparation"`
	PortScan            map[string]string `bson:"PortScan" json:"PortScan"`
	PortFingerprint     map[string]string `bson:"PortFingerprint" json:"PortFingerprint"`
	AssetMapping        map[string]string `bson:"AssetMapping" json:"AssetMapping"`
	AssetHandle         map[string]string `bson:"AssetHandle" json:"AssetHandle"`
	URLScan             map[string]string `bson:"URLScan" json:"URLScan"`
	WebCrawler          map[string]string `bson:"WebCrawler" json:"WebCrawler"`
	URLSecurity         map[string]string `bson:"URLSecurity" json:"URLSecurity"`
	DirScan             map[string]string `bson:"DirScan" json:"DirScan"`
	VulnerabilityScan   map[string]string `bson:"VulnerabilityScan" json:"VulnerabilityScan"`
	PassiveScan         map[string]string `bson:"PassiveScan" json:"PassiveScan"`
}

// TemplateList 模板列表
type TemplateList struct {
	List  []ScanTemplate `json:"list"`
	Total int64          `json:"total"`
}

// TemplateListItem 模板列表项
type TemplateListItem struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}
