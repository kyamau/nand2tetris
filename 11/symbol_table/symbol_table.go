package symbol_table

type Entry struct {
	varName string
	varType string
	varKind string
	varNum  int
}

type SymbolTable struct {
	name        string
	entries     []Entry
	kindCounter map[string]int
}

func NewSymbolTable(name string) *SymbolTable {
	return &SymbolTable{name, []Entry{}, make(map[string]int)}
}

func (s *SymbolTable) Define(varName string, varType string, varKind string) {
	s.entries = append(s.entries, Entry{varName, varType, varKind, s.kindCounter[varKind]})
	s.kindCounter[varKind]++
}

func (s *SymbolTable) VarCount(varKind string) int {
	return s.kindCounter[varKind]
}

func (s *SymbolTable) KindOf(varName string) (string, bool) {
	for _, e := range s.entries {
		if e.varName == varName {
			return e.varKind, true
		}
	}
	return "", false
}

func (s *SymbolTable) TypeOf(varName string) (string, bool) {
	for _, e := range s.entries {
		if e.varName == varName {
			return e.varType, true
		}
	}
	return "", false
}

func (s *SymbolTable) IndexOf(varName string) (int, bool) {
	for _, e := range s.entries {
		if e.varName == varName {
			return e.varNum, true
		}
	}
	return -1, false
}

func (s *SymbolTable) Name() string {
	return s.name
}
