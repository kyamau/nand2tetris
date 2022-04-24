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

func NewEmptyClassSymbolTable(name string) *SymbolTable {
	return &SymbolTable{name, []Entry{}, make(map[string]int)}
}

func NewEmptyMethodSymbolTable(name string) *SymbolTable {
	return &SymbolTable{name, []Entry{}, make(map[string]int)}
}

func (s *SymbolTable) Define(varName string, varType string, varKind string) {
	s.entries = append(s.entries, Entry{varName, varType, varKind, s.kindCounter[varKind]})
	s.kindCounter[varKind]++
}

func (s *SymbolTable) VarCount(varKind string) int {
	return s.kindCounter[varKind]
}

func (s *SymbolTable) KindOf(varName string) (bool, string) {
	for _, e := range s.entries {
		if e.varName == varName {
			return true, e.varKind
		}
	}
	return false, ""
}

func (s *SymbolTable) TypeOf(varName string) (bool, string) {
	for _, e := range s.entries {
		if e.varName == varName {
			return true, e.varType
		}
	}
	return false, ""
}

func (s *SymbolTable) IndexOf(varName string) (bool, int) {
	for _, e := range s.entries {
		if e.varName == varName {
			return true, e.varNum
		}
	}
	return false, -1
}
