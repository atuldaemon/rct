package parking

import (
	"errors"
	"sync"

	"strconv"

	"sort"

	"github.com/umahmood/haversine"
)

// The parking store which stores the information about the spots

type SpotType int

const (
	all      SpotType = 0
	free     SpotType = 1
	reserved SpotType = 2
)

type ParkingStore interface {
	Get(t SpotType) ([]Spot, error)
	Create(Spot) (Spot, error)
	Update(Spot) (Spot, error)
	Delete(id int) error
	Search(lat, lon, radius string, metric SearchMetric) ([]ExtendedSpot, error)
	FindById(id int) (Spot, error)
}

type Spot struct {
	ID         int    `json:"id"`
	Lat        string `json:"lat"`
	Lon        string `json:"lon"`
	Cost       string `json:"cost"`
	IsReserved bool   `json:"isReserved"`
	Address    string `json:"address,omitempty"`
}

// ExtendedSpot stores the distance of the spot from the searched location
type ExtendedSpot struct {
	Spot
	// Distance in meters
	Distance float64 `json:"distance"`
}

func MakeNewExtendedSpot(spot Spot, distanceKM float64) ExtendedSpot {
	esp := ExtendedSpot{Distance: distanceKM * 1000}
	esp.ID = spot.ID
	esp.IsReserved = spot.IsReserved
	esp.Lat = spot.Lat
	esp.Lon = spot.Lon
	esp.Address = spot.Address
	esp.Cost = spot.Cost
	return esp
}

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrNotFound        = errors.New("not found")
	ErrInvalidReq      = errors.New("invalid request")
	ErrInternal        = errors.New("internal data error")
)

// In memory store that stores the parking database in memory
type InMemStore struct {
	mtx sync.RWMutex // controls access to the map m
	m   map[int]Spot
}

func NewInMemParkingStore() (ParkingStore, error) {
	s := &InMemStore{m: make(map[int]Spot, 0)}
	ss := createDefaultSpots()
	for _, sp := range ss {
		s.m[sp.ID] = sp
	}
	return s, nil
}

func (s *InMemStore) Get(t SpotType) ([]Spot, error) {
	switch t {
	case all:
		return s.getAll()
	case free:
		return s.getFree()
	case reserved:
		return s.getReserved()
	default:
		return nil, ErrInvalidReq

	}
	return nil, nil
}

// CRUD ops on Parking store

func (s *InMemStore) Create(st Spot) (Spot, error) {
	// TODO: implement this - not required for the test since we will be only using the dummy data
	return Spot{}, nil
}

func (s *InMemStore) Update(st Spot) (Spot, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	sp, ok := s.m[st.ID]
	if !ok {
		return Spot{}, ErrInconsistentIDs
	}
	sp.IsReserved = st.IsReserved
	s.m[sp.ID] = sp

	return sp, nil
}

func (s *InMemStore) Delete(id int) error {
	// TODO: implement this - perhaps not required... Definitely not required for the test
	return nil
}

func (s *InMemStore) getAll() ([]Spot, error) {
	ss := make([]Spot, 0)
	for _, sp := range s.m {
		ss = append(ss, sp)
	}
	return ss, nil
}

func (s *InMemStore) getFree() ([]Spot, error) {
	ss := make([]Spot, 0)
	for _, sp := range s.m {
		if sp.IsReserved == false {
			ss = append(ss, sp)
		}
	}
	return ss, nil
}

func (s *InMemStore) getReserved() ([]Spot, error) {
	ss := make([]Spot, 0)
	for _, sp := range s.m {
		if sp.IsReserved == true {
			ss = append(ss, sp)
		}
	}
	return ss, nil
}

func (s *InMemStore) FindById(id int) (Spot, error) {
	if sp, ok := s.m[id]; ok {
		return sp, nil
	}
	return Spot{}, ErrNotFound
}

// Search searches for the neighbouring spots based on the searchmetric
// SearchMetric can be one of cost and distance
// The search results will be ordered based on the metric
func (s *InMemStore) Search(lat, lon, radius string, metric SearchMetric) ([]ExtendedSpot, error) {
	ess := make([]ExtendedSpot, 0)
	latFloat, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return nil, ErrInvalidReq
	}
	lonFloat, err := strconv.ParseFloat(lon, 64)
	if err != nil {
		return nil, ErrInvalidReq
	}
	radFloat, err := strconv.ParseFloat(radius, 64)
	if err != nil {
		return nil, ErrInvalidReq
	}

	// Make use of the third party havesine library for computing the distance between two spots
	p1 := haversine.Coord{Lat: latFloat, Lon: lonFloat}
	for _, sp := range s.m {
		p2LatFloat, err := strconv.ParseFloat(sp.Lat, 64)
		if err != nil {
			return nil, ErrInternal
		}
		p2LonFloat, err := strconv.ParseFloat(sp.Lon, 64)
		p2 := haversine.Coord{Lat: p2LatFloat, Lon: p2LonFloat}
		_, km := haversine.Distance(p1, p2)
		if km < radFloat/1000 {
			esp := MakeNewExtendedSpot(sp, km)
			ess = append(ess, esp)
		}
	}
	return SortSpots(ess, metric)
}

func SortSpots(ess []ExtendedSpot, metric SearchMetric) ([]ExtendedSpot, error) {
	switch metric {
	case "dist":
		sort.Slice(ess, func(i, j int) bool {
			return ess[i].Distance < ess[j].Distance
		})
		return ess, nil
	case "cost":
		sort.Slice(ess, func(i, j int) bool {
			return ess[i].Cost < ess[j].Cost
		})
		return ess, nil
	}
	return ess, nil
}



// Dummy data for testing
func createDefaultSpots() []Spot {
	ss := []Spot{
		Spot{1, "44.968046", "-94.420307", "100", false, "address 1"},
		Spot{2, "44.33328", "-89.132008", "10", false, "address 2"},
		Spot{3, "33.755787", "-116.359998", "80", false, "address 3"},
		Spot{4, "33.844843", "-116.54911", "70", false, "address 4"},
		Spot{5, "44.92057", "-93.44786", "90", false, "address 5"},
	}
	return ss
}
