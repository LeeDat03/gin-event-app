package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type AttendeeModel struct {
	DB *sql.DB
}

type Attendee struct {
	ID      int `json:"id"`
	UserId  int `json:"userId"`
	EventId int `json:"eventId"`
}

func (m *AttendeeModel) Insert(attend *Attendee) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		INSERT INTO attendees (user_id, event_id)
		VALUES ($1, $2)
		RETURNING id;	
	`

	fmt.Println(attend.UserId, attend.EventId)

	err := m.DB.QueryRowContext(ctx, stmt, attend.UserId, attend.EventId).Scan(&attend.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m *AttendeeModel) GetByEventAndAttendee(eventId, userId int) (*Attendee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT * FROM attendees WHERE event_id=$1 AND user_id=$2
	`
	var attendee Attendee
	err := m.DB.QueryRowContext(ctx, query, eventId, userId).Scan(&attendee.ID, &attendee.UserId, &attendee.EventId)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &attendee, nil
}

func (m *AttendeeModel) GetAttendeesByEvent(id int) ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	SELECT u.id, u.name, u.email
	FROM users u
	JOIN attendees a ON u.id = a.user_id
	WHERE a.event_id = $1
	`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*User{}
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (m *AttendeeModel) Delete(eventId, userId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		DELETE FROM attendees
		WHERE event_id = $1 AND user_id = $2;
	`

	_, err := m.DB.ExecContext(ctx, stmt, eventId, userId)
	if err != nil {
		return err
	}
	return nil
}
