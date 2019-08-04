package navigator

//Address stores the name of a Address
type Address struct {
	countries map[string]Country
}

//Country stores the name of a country
type Country struct {
	name   string
	cities map[string]City
}

//City stores name and streets
type City struct {
	ID      int
	Streets map[string]Street
}

//Street stores name of the street
type Street struct {
	ID      int64
	Name    string
	City    string
	Country string
	Pos     []*Pos
	Con     []int64
}

//Pos test
type Pos struct {
	Lat float64
	Lon float64
}

//Navigator is
type Navigator struct {
	streets map[int64]Street
}

func NewNavigator() Navigator {
	return Navigator{
		streets: make(map[int64]Street),
	}
}

//AddStreet is
func (n Navigator) AddStreet(s Street) {
	n.streets[s.ID] = s
}
