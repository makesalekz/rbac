package enum

type DummyType string

const (
	Type1 DummyType = "TYPE_1"
	Type2 DummyType = "TYPE_2"
)

func chatTypeValues() []DummyType {
	return []DummyType{Type1, Type2}
}

func (DummyType) Values() (kinds []string) {
	for _, value := range chatTypeValues() {
		kinds = append(kinds, string(value))
	}
	return
}

func (m DummyType) Value() string {
	return string(m)
}

func (m DummyType) IsValid() bool {
	for _, value := range chatTypeValues() {
		if m == value {
			return true
		}
	}
	return false
}
