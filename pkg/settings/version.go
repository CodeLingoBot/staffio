package settings

const (
	NAME = "staffio"
)

var (
	buildVersion = "dev"
)

func IsDevelop() bool {
	return "dev" == buildVersion
}
