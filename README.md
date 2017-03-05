# uclapi-go
A wrapper library in Go for uclapi
##Introduction
```go
import github.com/Maaarcocr/uclapi-go
```
The first thing you need is a Wrapper for the Room Booking API. You can create one with:
```go
wrapper := uclapi.NewRoomBookingWrapper("your token")
```
The wrapper has four methods:
* GetRooms
* GetBookings
* NextPage
* GetEquipment

##Documentation
You can find the documentation for this library at: https://godoc.org/github.com/Maaarcocr/uclapi-go
