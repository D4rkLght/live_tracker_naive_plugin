package connectiontracker

import "sync"


type Storage struct {
	mu      sync.RWMutex
	Tracker *Tracker
}

func NewStorage() *Storage {
	return &Storage{
		Tracker: &Tracker{
			Users: make(map[string]*User),
		},
	}
}


func (s *Storage) AddConnection(username string, connID string, conn *UserConnection) bool {

	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.Tracker.Users[username]

	if !exists {
		user = &User{
			Active:      0,
			RejectedAttempts:      0,
			Connections: make(map[string]*UserConnection),
		}

		s.Tracker.Users[username] = user
	}

	// проверяем лимит
	if user.Active >= 3 {
		user.RejectedAttempts++
		return false
	}

	user.Connections[connID] = conn
	user.Active++

	return true
}

func (s *Storage) RemoveConnection(username string, id string) {

	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.Tracker.Users[username]

	if !exists {
		return
	}

	_, exists = user.Connections[id]

	if !exists {
		return
	}

	delete(user.Connections, id)

	if user.Active > 0 {
		user.Active--
	}

	// если больше нет соединений - удаляем пользователя
	if user.Active == 0 {
		delete(s.Tracker.Users, username)
	}
}


func (s *Storage) ListUsers() map[string]*User {

	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]*User)

	for username, user := range s.Tracker.Users {

		copyUser := &User{
			Active:      user.Active,
			RejectedAttempts:      user.RejectedAttempts,
			Connections: make(map[string]*UserConnection),
		}

		for id, conn := range user.Connections {
			copyUser.Connections[id] = conn
		}

		result[username] = copyUser
	}

	return result
}


var storage = NewStorage()
