package connectiontracker

import "time"

type UserConnection struct {
    Started time.Time `json:"started"`
    Host    string    `json:"host,omitempty"`
}

type User struct {
    Active      int                           `json:"active"`
    RejectedAttempts  int                           `json:"rejected_attempts"`
    Connections map[string]*UserConnection   `json:"connections"`
}

type Tracker struct {
    Users map[string]*User `json:"users"`
}