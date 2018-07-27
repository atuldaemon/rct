package parking

import (
	"testing"
)

func TestFindById(t *testing.T) {
	inMemStore, err := NewInMemParkingStore()

	if err != nil {
		t.Error("Failed to create inmem store")
	}
	t.Log("Created inmem store")

	service := NewService(inMemStore)
	t.Log("Created parking service")

	s, err := service.FindById(nil, "1")

	if err != nil {
		t.Error("Error in Find")
	}
	if s.ID != 1 {
		t.Error("Incorrect spot returned")
	}
	t.Log("Found spot")
}

func TestGetAll(t *testing.T) {
	inMemStore, err := NewInMemParkingStore()

	if err != nil {
		t.Error("Failed to create inmem store")
	}
	t.Log("Created inmem store")

	service := NewService(inMemStore)
	t.Log("Created parking service")

	ss, err := service.GetAll(nil)
	if err != nil {
		t.Error("Error in Find")
	}
	if len(ss) != 5 {
		t.Error("Incorrect results returned")
	}
	t.Log("Got all spots")
}

func TestSearchByCost(t *testing.T) {
	inMemStore, err := NewInMemParkingStore()

	if err != nil {
		t.Error("Failed to create inmem store")
	}
	t.Log("Created inmem store")

	service := NewService(inMemStore)
	t.Log("Created parking service")

	curLat := "33.755787"
	curLon := "-116.359998"

	ss, err := service.Search(nil, curLat, curLon, "10000", "cost")
	if err != nil {
		t.Error("Error in Search")
	}
	if len(ss) != 1 && ss[0].ID != 3 {
		t.Error("Incorrect results returned")
	}
	t.Log("Search spot by cost")
}

func TestSearchByDist(t *testing.T) {
	inMemStore, err := NewInMemParkingStore()

	if err != nil {
		t.Error("Failed to create inmem store")
	}
	t.Log("Created inmem store")

	service := NewService(inMemStore)
	t.Log("Created parking service")

	curLat := "33.755787"
	curLon := "-116.359998"

	ss, err := service.Search(nil, curLat, curLon, "10000", "dist")
	if err != nil {
		t.Error("Error in Search")
	}
	if len(ss) != 2 && ss[0].ID != 3 && ss[1].ID != 4 {
		t.Error("Incorrect results returned")
	}
	t.Log("Search spot by dist")
}
