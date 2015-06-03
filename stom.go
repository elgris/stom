package stom

type Policy uint8

const (
	PolicyUseDefault Policy = iota
	PolicyExclude

	defaultTag = "db"
)

var defaultStom = stom{[]string{defaultTag}, PolicyUseDefault, nil}

type stom struct {
	tags         []string
	policy       Policy
	defaultValue interface{}
}

func (this *stom) SetTags(tags ...string) {
	if len(tags) == 0 {
		this.tags = []string{defaultTag}
		return
	}

	this.tags = tags
}

func (this *stom) SetDefault(defaultValue interface{}) {
	this.defaultValue = defaultValue
}

func (this *stom) SetPolicy(policy Policy) {
	this.policy = policy
}

func SetDefault(defaultValue interface{}) {
	defaultStom.SetDefault(defaultValue)
}

func SetTags(tags ...string) {
	defaultStom.SetTags(tags...)
}

func SetPolicy(policy Policy) {
	defaultStom.SetPolicy(policy)
}

func ToMap(s interface{}) map[string]interface{} {
	return map[string]interface{}{}
}
