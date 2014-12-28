package main

type Humanoid struct {
    id        int64
    name      string
    starCount int64
}

type HumanoidSlice []*Humanoid

func SortHumanoids (m map[int64]*Humanoid) HumanoidSlice {
    slice := HumanoidSliceFromMap(m)
    sort.Sort(sort.Reverse(slice))
    return slice
}

func HumanoidSliceFromMap (m map[int64]*Humanoid) HumanoidSlice {
	slice := make(HumanoidSlice, 0, len(m))

	for _, human := range m {
		slice = append(slice, human)
	}

	return slice
}

// Hurray internet! Implement the Sort interface. % godoc sort

func (hs HumanoidSlice) Len() int {
	return len(hs)
}

func (hs HumanoidSlice) Less(i, j int) bool {
	return hs[i].starCount < hs[j].starCount
}

func (hs HumanoidSlice) Swap(i, j int) {
	hs[i], hs[j] = hs[j], hs[i]
}
