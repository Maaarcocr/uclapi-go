package main

import (
	"fmt"
	"time"

	"github.com/Maaarcocr/uclapi-go"
)

func main() {
	room_bookings := uclapi.NewRoomBookingClient("your token")

	rooms, errs := room_bookings.GetRooms(uclapi.RoomOptList{})
	fmt.Println(rooms, errs)

	bookings, errs := room_bookings.GetBookings(uclapi.BookingOptList{
		Day: uclapi.Day(time.Now()),
	})
	fmt.Println(bookings.Bookings, errs)

	next, errs := room_bookings.NextPage(bookings)
	fmt.Println(next, errs)

	equipments, errs := room_bookings.GetEquipment("B05A", "044")
	fmt.Println(equipments, errs)
}
