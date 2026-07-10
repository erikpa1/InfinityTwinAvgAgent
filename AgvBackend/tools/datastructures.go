package tools

// StringSet of strings using map
type StringSet map[string]struct{}

// Add adds an element to the set
func (s StringSet) Add(element string) {
	s[element] = struct{}{} // struct{}{} occupies no memory
}

// Remove removes an element from the set
func (s StringSet) Remove(element string) {
	delete(s, element)
}

// Contains checks if an element is in the set
func (s StringSet) Contains(element string) bool {
	_, exists := s[element]
	return exists
}

// Size returns the number of elements in the set
func (s StringSet) Size() int {
	return len(s)
}

// Clear removes all elements from the set
func (s StringSet) Clear() {
	for k := range s {
		delete(s, k)
	}
}
func (s StringSet) ToArray() []string {
	tmp := make([]string, s.Size())

	index := 0
	for key, _ := range s {
		tmp[index] = key
		index++
	}

	return tmp
}
