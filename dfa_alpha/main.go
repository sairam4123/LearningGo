package main

import (
	"fmt"
)

type Alphabet struct {
	Element string
}

type State struct {
	Name      string
	IsFinal   bool
	IsInitial bool
}

type Transition struct {
	Name         string
	Current      *State
	Next         *State
	Input        Alphabet
	IsCustomName bool
}

func (s *State) To(nextState *State) Transition {
	return Transition{
		Name:         s.Name + "->" + nextState.Name,
		Current:      s,
		Next:         nextState,
		IsCustomName: false,
	}
}

func (s *State) GetName() string {
	var sName = s.Name
	if s.IsInitial {
		sName = "->" + sName
	}
	if s.IsFinal {
		sName = sName + "*"
	}
	return sName
}

func (s State) Final() State {
	return State{
		Name:      s.Name,
		IsFinal:   true,
		IsInitial: s.IsInitial,
	}
}

func (s State) Init() State {
	return State{
		Name:      s.Name,
		IsFinal:   s.IsFinal,
		IsInitial: true,
	}
}

func (t Transition) WithIp(input Alphabet) Transition {
	return Transition{
		Name:         t.Name,
		Input:        input,
		Next:         t.Next,
		Current:      t.Current,
		IsCustomName: false,
	}
}
func (t Transition) WithName(name string) Transition {
	return Transition{
		Name:         name,
		Input:        t.Input,
		Next:         t.Next,
		Current:      t.Current,
		IsCustomName: true,
	}
}

func (t *Transition) GetName() string {
	var tName = t.Name
	if !t.IsCustomName {
		tName = t.Current.GetName() + "--" + t.Input.Element + "->" + t.Next.GetName()
	}
	return tName
}

type DFA struct {
	states      []*State
	transitions []Transition
	alphabets   []Alphabet
}

func (d DFA) AddAlphabet(elements ...Alphabet) DFA {
	d.alphabets = append(d.alphabets, elements...)
	return d
}

func (d DFA) GetBestTransition(curState *State, input Alphabet) (Transition, error) {
	for _, v := range d.transitions {
		if v.Current.GetName() == curState.GetName() && v.Input == input {
			return v, nil
		}
	}
	return Transition{}, fmt.Errorf("Cannot find transition")
}
func (d DFA) GetInititialState() (State, error) {
	for _, v := range d.states {
		if v.IsInitial {
			return *v, nil
		}
	}
	return State{}, fmt.Errorf("Cannot find initial state")
}

type DFAResult struct {
	isAccepted bool
	steps      []DFARunStep
	err        string
}

type DFARunStep struct {
	fromState     State
	toState       State
	inputConsumed string
}

func (d DFA) Run(input string) (DFAResult, error) {
	var curState, err = d.GetInititialState()
	var steps = make([]DFARunStep, 0)
	if err != nil {
		return DFAResult{}, err
	}
	for _, ip := range input {
		var ip_a = Alphabet{Element: string(ip)}
		var trans, err = d.GetBestTransition(&curState, ip_a)
		if err != nil {
			return DFAResult{}, fmt.Errorf(
				"Cannot find any transition to reach final state, cur state %s i/p: %s", curState.Name, string(ip),
			)
		}
		curState = *trans.Next
		steps = append(steps, DFARunStep{
			fromState:     *trans.Current,
			toState:       *trans.Next,
			inputConsumed: string(ip),
		})

	}

	if curState.IsFinal {
		return DFAResult{
			isAccepted: true,
			steps:      steps,
		}, nil
	}

	return DFAResult{
		isAccepted: false,
		steps:      steps,
	}, nil
}

func CreateAB(element string) Alphabet {
	return Alphabet{
		Element: element,
	}
}

func CreateState(name string) State {
	return State{
		Name:      name,
		IsFinal:   false,
		IsInitial: false,
	}
}

func (d DFA) AddState(states ...*State) DFA {
	d.states = append(d.states, states...)
	return d
}

func (d DFA) AddTransition(trans ...Transition) DFA {
	d.transitions = append(d.transitions, trans...)
	return d
}

func (d DFA) Validate() (DFA, error) {
	// check for initial state
	var intitialStateCount = CountInitialStates(d)
	if intitialStateCount != 1 {
		return d, fmt.Errorf("Number of initial states must be exactly one. Initial States: %d", intitialStateCount)
	}
	if ContainsDuplicateTransitions(d) {
		return d, fmt.Errorf("Duplicate transitions detected.")
	}
	if !AllStatesReachable(d) {
		return d, fmt.Errorf("Not all states reachable from initial state")
	}

	return d, nil
}

func Diff[T comparable](a map[T]struct{}, b map[T]struct{}) map[T]struct{} {
	var diff = make(map[T]struct{}, 0)
	for k := range a {
		if _, ok := b[k]; !ok {
			diff[k] = struct{}{}
		}
	}
	return diff
}

func (d DFA) GetTransitionStates(curState *State) []*State {
	var reachableStates = make([]*State, 0)
	for _, v := range d.transitions {
		if v.Current.GetName() == curState.GetName() {
			reachableStates = append(reachableStates, v.Next)
		}
	}

	return reachableStates

}

func (d DFA) GetInputTransitionStates(curState *State, input Alphabet) []*State {
	var reachableStates = make([]*State, 0)
	for _, v := range d.transitions {
		if v.Current.GetName() == curState.GetName() && v.Input == input {
			reachableStates = append(reachableStates, v.Next)
		}
	}

	return reachableStates
}

func AllStatesReachable(d DFA) bool {
	var initialState, err = d.GetInititialState()
	if err != nil {
		return false
	}

	var queue = make([]*State, 0)
	queue = append(queue, &initialState)
	fmt.Println("Queue", queue)
	var visitedMap = make(map[string]struct{}, 0)
	for len(queue) != 0 {
		var elem = queue[0]
		fmt.Println("Queue", queue)
		queue = queue[1:] // discard top element

		if _, ok := visitedMap[elem.GetName()]; !ok {
			visitedMap[elem.GetName()] = struct{}{}
			var reachableStates = d.GetTransitionStates(elem)
			for _, v := range reachableStates {
				queue = append(queue, v)
			}
		}
	}

	var statesAsSet = make(map[string]struct{}, 0)
	for _, v := range d.states {
		statesAsSet[v.GetName()] = struct{}{}
	}
	fmt.Println("States", statesAsSet)
	fmt.Println("Reachable States", visitedMap)

	var diff = Diff(statesAsSet, visitedMap)
	fmt.Println("Diff", diff)
	return len(diff) == 0
}

func ContainsDuplicateTransitions(d DFA) bool {
	for i, t1 := range d.transitions {
		for j, t2 := range d.transitions {
			if i == j {
				continue
			}
			if t1.Input == t2.Input && t1.Current.GetName() == t2.Current.GetName() {
				return true
			}
		}
	}
	return false
}

func CountInitialStates(d DFA) int {
	var count = 0
	for _, v := range d.states {
		if v.IsInitial {
			count++
		}
	}
	return count
}

func main() {

	var a = CreateAB("a")
	var b = CreateAB("b")
	var c = CreateAB("c")
	var d = CreateAB("d")

	// Build a DFA for a*bbc*dd
	var q1 = CreateState("q1").Init()
	var q2 = CreateState("q2")
	var q3 = CreateState("q3")
	var q4 = CreateState("q4")
	var q5 = CreateState("q5").Final()
	var q6 = CreateState("q6")
	// var q7 = CreateState("q7")

	var t1 = q1.To(&q1).WithIp(a)
	var t2 = q1.To(&q2).WithIp(b)
	var t3 = q2.To(&q3).WithIp(b)
	var t4 = q3.To(&q3).WithIp(c)
	var t6 = q3.To(&q4).WithIp(d)
	var t7 = q4.To(&q5).WithIp(d)
	var t8 = q5.To(&q6).WithIp(d)
	var t9 = q6.To(&q6).WithIp(d)

	var dfa, derr = DFA{}.AddAlphabet(a, b, c, d).AddState(&q1, &q2, &q3, &q4, &q5, &q6).AddTransition(t1, t2, t3, t4, t6, t7, t8, t9).Validate()
	if derr != nil {
		fmt.Println("DFA is invalid", derr)
		return
	}
	var res, err = dfa.Run("bbdd")

	if err != nil {
		fmt.Println("error occured", err)
		return
	}

	for _, v := range res.steps {
		fmt.Printf("Transitioning from %s -> %s consuming %s as input\n", v.fromState.GetName(), v.toState.GetName(), v.inputConsumed)
	}
	fmt.Println("Input Accepted:", res.isAccepted)

}
