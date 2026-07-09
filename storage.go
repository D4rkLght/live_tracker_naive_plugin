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


func (s *Storage) AddConnection(username string, connID string, conn *UserConnection) {

	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.Tracker.Users[username]

	if !exists {
		user = &User{
			Connections: make(map[string]*UserConnection),
		}

		s.Tracker.Users[username] = user
	}

	user.Connections[connID] = conn
}


func (s *Storage) RemoveConnection(username string, id string) {

	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.Tracker.Users[username]

	if !exists {
		return
	}


	delete(user.Connections, id)


	// если у пользователя больше нет соединений,
	// можно удалить его из списка

	if len(user.Connections) == 0 {
		delete(s.Tracker.Users, username)
	}
}


func (s *Storage) ListUsers() map[string]*User {

	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]*User)

	for username, user := range s.Tracker.Users {

		copyUser := &User{
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
