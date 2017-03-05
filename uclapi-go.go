// A wrapper library for uclapi.
package uclapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/parnurzeal/gorequest"
)

const base_url = "https://uclapi.com/roombookings/"
const header = "uclapi-roombookings-version"
const version = "1"

type RoomKind string

//Alias type for time.Time that follows the ISO 8601 formatting standard when is Encoded and Decoded.
type KloppTime time.Time

//Alias type for time.Time that follows the YYYYMMDD standard when is Encoded and Decoded.
type Day time.Time

const (
	ClassRoom      RoomKind = "CR"
	LectureTheatre RoomKind = "LT"
	SocialSpace    RoomKind = "SS"
	PublicCluster  RoomKind = "PC1"
)

//Wrapper for the room bookings api
type WrapperRoomBooking struct {
	token string
}

type RoomOptList struct {
	RoomID         string   `url:"roomid,omitempty"`
	RoomName       string   `url:"roomname,omitempty"`
	SiteId         string   `url:"siteid,omitempty"`
	SiteName       string   `url:"sitename,omitempty"`
	Classification RoomKind `url:"classification,omitempty"`
	Capacity       int      `url:"capacity,omitempty"`
}

type BookingOptList struct {
	RoomID        string    `url:"roomid,omitempty"`
	StartTime     KloppTime `url:"start_datetime,omitempty"`
	EndTime       KloppTime `url:"end_datetime,omitempty"`
	Day           Day       `url:"date,omitempty"`
	SiteId        string    `url:"siteid,omitempty"`
	Description   string    `url:"description,omitempty"`
	Contact       string    `url:"contact,omitempty"`
	ResultPerPage int       `url:"result_per_page,omitempty"`
	RoomName      string    `url:"roomname,omitempty"`
}

//Address is just an array of string that if concatened would create the entire address.
type Location struct {
	Address []string `json:"address"`
}

type Room struct {
	RoomID         string   `json:"roomid"`
	RoomName       string   `json:"roomname"`
	SiteId         string   `json:"siteid"`
	SiteName       string   `json:"sitename"`
	Classification RoomKind `json:"classification"`
	Capacity       int      `json:"capacity"`
	Automated      bool     `json:"automated"`
	Location       Location `json:"location"`
}

type Booking struct {
	SlotId      int       `json:"slotid,omitempty"`
	Contact     string    `json:"contact,omitempty"`
	StartTime   KloppTime `json:"start_time"`
	EndTime     KloppTime `json:"end_time"`
	RoomID      string    `json:"roomid,omitempty"`
	RoomName    string    `json:"roomname,omitempty"`
	siteid      string    `json:"siteid,omitempty"`
	WeekNumber  int       `json:"weeknumber,omitempty"`
	Phone       string    `json:"phone,omitempty"`
	Description string    `json:"description,omitempty"`
}

type Equipment struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Units       int    `json:"units"`
}

type ResponseRooms struct {
	Ok    bool   `json:"ok"`
	Rooms []Room `json:"rooms"`
}

type ResponseBookings struct {
	Ok             bool      `json:"ok"`
	NextPageExists bool      `json:"next_page_exists"`
	Count          int       `json:"count,omitempty"`
	PageToken      string    `json:"page_token"`
	Bookings       []Booking `json:"bookings"`
}

type ResponseEquipment struct {
	Ok         bool        `json:"ok"`
	Equipments []Equipment `json:"equipment"`
}

//NewRoomBookinsWrapper create a WrapperRoomBooking using the given token.
func NewRoomBookingWrapper(token string) WrapperRoomBooking {
	return WrapperRoomBooking{
		token: token,
	}
}

func performGetRequest(relative_url string, params string) ([]byte, []error) {
	_, body, errs := gorequest.New().Get(base_url+relative_url+params).Set(header, version).EndBytes()
	if errs != nil {
		return []byte{}, errs
	}
	return body, errs
}

func readRoomResponse(body []byte) (ResponseRooms, []error) {
	var response ResponseRooms
	json.Unmarshal(body, &response)
	if !response.Ok {
		api_error := make(map[string]string)
		json.Unmarshal(body, &api_error)
		return ResponseRooms{}, []error{errors.New(api_error["error"])}
	}
	return response, nil
}

func readBookingResponse(body []byte) (ResponseBookings, []error) {
	//	fmt.Println("RESPONSE: ", string(body))
	var response ResponseBookings
	json.Unmarshal(body, &response)
	if !response.Ok {
		api_error := make(map[string]string)
		json.Unmarshal(body, &api_error)
		return ResponseBookings{}, []error{errors.New(api_error["error"])}
	}
	return response, nil
}

func readEquipmentResponse(body []byte) (ResponseEquipment, []error) {
	fmt.Println("RESPONSE: ", string(body))
	var response ResponseEquipment
	json.Unmarshal(body, &response)
	if !response.Ok {
		api_error := make(map[string]string)
		json.Unmarshal(body, &api_error)
		return ResponseEquipment{}, []error{errors.New(api_error["error"])}
	}
	return response, nil
}

func (date *KloppTime) UnmarshalJSON(b []byte) error {
	date_string := string(b)[1 : len(b)-1]
	parsed_date, err := time.Parse("2006-01-02T15:04:05-07:00", date_string)
	if err != nil {
		fmt.Println(err)
		return err
	}
	*date = KloppTime(parsed_date)
	return nil
}

func (date KloppTime) EncodeValues(key string, v *url.Values) error {
	time_value := time.Time(date)
	if time_value.IsZero() {
		return nil
	}
	v.Add(key, time_value.Format("2006-01-02T15:04:05-07:00"))
	return nil
}

func (date Day) EncodeValues(key string, v *url.Values) error {
	time_value := time.Time(date)
	if time_value.IsZero() {
		return nil
	}
	v.Add(key, time_value.Format("20060102"))
	return nil
}

//GetRooms takes a RoomOptList object that contains the optional fields that can be used for the query to the rooms api.
//In case of success it returns a ResponseRooms object and nil for error.
//In case of error it returns an empty ResponseRooms object and an array for errors.
func (ucl WrapperRoomBooking) GetRooms(opt RoomOptList) (ResponseRooms, []error) {
	request_parameters, err := query.Values(opt)
	if err != nil {
		return ResponseRooms{}, []error{err}
	}
	request_parameters["token"] = []string{ucl.token}
	url_query := request_parameters.Encode()
	body, errs := performGetRequest("rooms?", url_query)
	if errs != nil {
		return ResponseRooms{}, errs
	}
	return readRoomResponse(body)
}

//GetBookings takes a BookingOptList object that contains the optional fields that can be used for the query to the bookings api.
//In case of success it returns a ResponseBookings object and nil for error.
//In case of error it returns an empty ResponseBookings object and an array for errors.
func (ucl WrapperRoomBooking) GetBookings(opt BookingOptList) (ResponseBookings, []error) {
	request_parameters, err := query.Values(opt)
	if err != nil {
		return ResponseBookings{}, []error{err}
	}
	request_parameters["token"] = []string{ucl.token}
	url_query := request_parameters.Encode()
	body, errs := performGetRequest("bookings?", url_query)
	if errs != nil {
		return ResponseBookings{}, errs
	}
	return readBookingResponse(body)
}

//NextPage exists because the bookings api is paginated.
//This function takes a previous ResponseBookings object and it returns the next page as a new ResponseBookings object.
//In case of success it returns a ResponseBookings object and nil for error.
//In case of failure it returns an empty ResponseBookings object and an array of errors.
//This function fails in the case where the ResponseBookins object used as a parameter has the field NextPageExists set to false.
func (client WrapperRoomBooking) NextPage(prev ResponseBookings) (ResponseBookings, []error) {
	if !prev.NextPageExists {
		return ResponseBookings{}, []error{errors.New("The next page doesn't exist")}
	}
	new_query := url.Values{}
	new_query["token"] = []string{client.token}
	new_query["page_token"] = []string{prev.PageToken}
	url_query := new_query.Encode()
	body, errs := performGetRequest("bookings?", url_query)
	if errs != nil {
		return ResponseBookings{}, errs
	}
	return readBookingResponse(body)

}

//GetRooms takes a roomId and a siteId used for the query to the equipment api.
//In case of success it returns a ResponseEquipment object and nil for error.
//In case of error it returns an empty ResponseEquipment object and an array for errors.
func (client WrapperRoomBooking) GetEquipment(roomId string, siteId string) (ResponseEquipment, []error) {
	new_query := url.Values{}
	new_query["token"] = []string{client.token}
	new_query["roomid"] = []string{roomId}
	new_query["siteid"] = []string{siteId}
	url_query := new_query.Encode()
	body, errs := performGetRequest("equipment?", url_query)
	if errs != nil {
		return ResponseEquipment{}, errs
	}
	return readEquipmentResponse(body)
}
