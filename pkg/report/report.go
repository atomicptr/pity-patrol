package report

type Reward struct {
	Name  string
	Count int
	Image string
}

type Field struct {
	Key   string
	Value string
}

type Report struct {
	WasClaimed   bool
	Reward       *Reward
	CustomFields []Field
}
