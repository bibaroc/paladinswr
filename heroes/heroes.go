package heroes

func All() []string {
	return append(append(append(
		Flanker(),
		FrontLine()...),
		Damage()...),
		Support()...)
}

func Flanker() []string {
	return []string{
		"androxus", "buck", "evie", "koga",
		"lex", "maeve", "moji", "skye",
		"talus", "vora", "zhin",
	}
}

func FrontLine() []string {
	return []string{
		"ash", "atlas", "barik", "fernando",
		"inara", "khan", "makoa", "raum",
		"ruckus", "terminus", "torvald", "yagorath"}
}

func Damage() []string {
	return []string{
		"bomb-king", "cassie", "dredge", "drogoz",
		"imani", "kinessa", "lian", "sha-lin",
		"strix", "tiberius", "tyra", "viktor",
		"vivian", "willo",
	}
}

func Support() []string {
	return []string{
		"corvus", "furia", "grohk", "grover",
		"io", "jenos", "mal-damba", "pip",
		"seris", "ying",
	}
}
