package booking

import (
	"testing"

	"time"

	"strconv"

	"github.com/atuldaemon/rct/parking"
)

func TestBook(t *testing.T) {

	pInMemStore, err := parking.NewInMemParkingStore()

	if err != nil {
		t.Error("Failed to create parking inmem store")
	}
	t.Log("Created parking inmem store")

	pService := parking.NewService(pInMemStore)
	t.Log("Created parking service")

	bInMemStore, err := NewInMemBookingStore()

	if err != nil {
		t.Error("Failed to create booking inmem store")
	}
	t.Log("Created inmem booking store")

	bService := NewService(bInMemStore, pService)
	t.Log("Created booking service")

	b, err := bService.Book(nil, "1", time.Now(), time.Duration(30*time.Minute))

	if err != nil {
		t.Error("Error in booking")
	}
	if b.SpotId != 1 {
		t.Error("Incorrect spot booked")
	}
	t.Log("Booked spot")
}

func TestBookMultiple(t *testing.T) {

	pInMemStore, err := parking.NewInMemParkingStore()

	if err != nil {
		t.Error("Failed to create parking inmem store")
	}
	t.Log("Created parking inmem store")

	pService := parking.NewService(pInMemStore)
	t.Log("Created parking service")

	bInMemStore, err := NewInMemBookingStore()

	if err != nil {
		t.Error("Failed to create booking inmem store")
	}
	t.Log("Created inmem booking store")

	bService := NewService(bInMemStore, pService)
	t.Log("Created booking service")

	b, err := bService.Book(nil, "1", time.Now(), time.Duration(30*time.Minute))

	if err != nil {
		t.Error("Error in booking")
	}
	if b.SpotId != 1 {
		t.Error("Incorrect spot booked")
	}
	t.Log("Booked spot")

	_, err = bService.Book(nil, "1", time.Now(), time.Duration(30*time.Minute))

	if err == nil {
		t.Error("Expecting error in booking the same spot again")
	} else {
		t.Log("Denied booking the same spot again")
	}
}

func TestDeleteBooking(t *testing.T) {

	pInMemStore, err := parking.NewInMemParkingStore()

	if err != nil {
		t.Error("Failed to create parking inmem store")
	}
	t.Log("Created parking inmem store")

	pService := parking.NewService(pInMemStore)
	t.Log("Created parking service")

	bInMemStore, err := NewInMemBookingStore()

	if err != nil {
		t.Error("Failed to create booking inmem store")
	}
	t.Log("Created inmem booking store")

	bService := NewService(bInMemStore, pService)
	t.Log("Created booking service")

	b, err := bService.Book(nil, "1", time.Now(), time.Duration(30*time.Minute))

	if err != nil {
		t.Error("Error in booking")
	}
	if b.SpotId != 1 {
		t.Error("Incorrect spot booked")
	}
	t.Log("Booked spot")

	err = bService.Delete(nil, strconv.Itoa(b.ID))
	if err != nil {
		t.Error("Could not free spot")
	}

	// book the same spot again
	_, err = bService.Book(nil, "1", time.Now(), time.Duration(30*time.Minute))

	if err != nil {
		t.Error("Could not book a free spot")
	} else {
		t.Log("booked a spot which was released")
	}
}
