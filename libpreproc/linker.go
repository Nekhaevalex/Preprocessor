package libpreproc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//LinkerScript - representation of linker information
type LinkerScript struct {
	ENTRY        string                     `json:"ENTRY"`
	ARCHITECTURE string                     `json:"ARCHITECTURE"`
	SEARCHDIR    string                     `json:"SEARCH_DIR"`
	OUTPUT       string                     `json:"OUTPUT"`
	RELJMP       bool                       `json:"REL_JMP"`
	MEMORY       map[string]MemoryPartition `json:"MEMORY"`
	SECTIONS     map[string][]string        `json:"SECTIONS"`
}

//MemoryPartition - representation of memory partition
type MemoryPartition struct {
	ORIGIN int `json:"ORIGIN"`
	LENGTH int `json:"LENGTH"`
}

//OpenLinkerScript - parses linker script
func OpenLinkerScript(filename string) (LinkerScript, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return LinkerScript{}, err
	}
	var script LinkerScript
	if err := json.Unmarshal(file, &script); err != nil {
		return LinkerScript{}, err
	}
	return script, nil
}

//GetPartitionList returns memory partition list
func (l *LinkerScript) GetPartitionList() ([]string, error) {
	if l == nil {
		return nil, fmt.Errorf("no parsed linkre script")
	}
	keys := make([]string, len(l.MEMORY))
	i := 0
	for k := range l.MEMORY {
		keys[i] = k
		i++
	}
	return keys, nil
}

//GetPartitionSections returns section order in specified partition
func (l *LinkerScript) GetPartitionSections(partition string) ([]string, error) {
	if l == nil {
		return nil, fmt.Errorf("no parsed linkre script")
	}
	return l.SECTIONS[partition], nil
}
