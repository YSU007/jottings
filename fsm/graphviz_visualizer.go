package fsm

import (
	"bytes"
	"fmt"
	"sort"
)

// Visualize outputs a visualization of FSM in Graphviz format.
func Visualize(fsm *FSM) string {
	var buf bytes.Buffer

	// we sort the key alphabetically to have a reproducible graph output
	sortedEKeys := getSortedTransitionKeys(fsm.transitions)
	sortedStateKeys, _ := getSortedStates(fsm.transitions)

	writeHeaderLine(&buf)
	writeTransitions(&buf, sortedEKeys, fsm.transitions)
	writeStates(&buf, fsm.current, sortedStateKeys)
	writeFooter(&buf)

	return buf.String()
}

func writeHeaderLine(buf *bytes.Buffer) {
	buf.WriteString(`digraph fsm {`)
	buf.WriteString("\n")
}

func writeTransitions(buf *bytes.Buffer, sortedEKeys []eKey, transitions map[eKey]string) {
	for _, k := range sortedEKeys {
		v := transitions[k]
		buf.WriteString(fmt.Sprintf(`    "%s" -> "%s" [ label = "%s" ];`, k.src, v, k.event))
		buf.WriteString("\n")
	}

	buf.WriteString("\n")
}

func writeStates(buf *bytes.Buffer, current string, sortedStateKeys []string) {
	for _, k := range sortedStateKeys {
		if k == current {
			buf.WriteString(fmt.Sprintf(`    "%s" [color = "red"];`, k))
		} else {
			buf.WriteString(fmt.Sprintf(`    "%s";`, k))
		}
		buf.WriteString("\n")
	}
}

func writeFooter(buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintln("}"))
}

func getSortedTransitionKeys(transitions map[eKey]string) []eKey {
	// we sort the key alphabetically to have a reproducible graph output
	sortedTransitionKeys := make([]eKey, 0)

	for transition := range transitions {
		sortedTransitionKeys = append(sortedTransitionKeys, transition)
	}
	sort.Slice(sortedTransitionKeys, func(i, j int) bool {
		if sortedTransitionKeys[i].src == sortedTransitionKeys[j].src {
			return sortedTransitionKeys[i].event < sortedTransitionKeys[j].event
		}
		return sortedTransitionKeys[i].src < sortedTransitionKeys[j].src
	})

	return sortedTransitionKeys
}

func getSortedStates(transitions map[eKey]string) ([]string, map[string]string) {
	statesToIDMap := make(map[string]string)
	for transition, target := range transitions {
		if _, ok := statesToIDMap[transition.src]; !ok {
			statesToIDMap[transition.src] = ""
		}
		if _, ok := statesToIDMap[target]; !ok {
			statesToIDMap[target] = ""
		}
	}

	sortedStates := make([]string, 0, len(statesToIDMap))
	for state := range statesToIDMap {
		sortedStates = append(sortedStates, state)
	}
	sort.Strings(sortedStates)

	for i, state := range sortedStates {
		statesToIDMap[state] = fmt.Sprintf("id%d", i)
	}
	return sortedStates, statesToIDMap
}
