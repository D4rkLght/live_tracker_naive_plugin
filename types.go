package connectiontracker

import "time"

type UserConnection struct {
    IP      string    `json:"ip"`
    Started time.Time `json:"started"`
    Host    string    `json:"host,omitempty"`
}

type User struct {
    Connections map[string]*UserConnection `json:"connections"`
}

type Tracker struct {
	Users map[string]*User
}
