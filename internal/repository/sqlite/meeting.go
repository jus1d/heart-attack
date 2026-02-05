package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jus1d/kypidbot/internal/domain"
)

type MeetingRepo struct {
	db *sql.DB
}

func NewMeetingRepo(d *DB) *MeetingRepo {
	return &MeetingRepo{db: d.db}
}

func (r *MeetingRepo) SaveMeeting(ctx context.Context, m *domain.Meeting) (int64, error) {
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO meetings (pair_id, place_id, time)
		VALUES (?, ?, ?)`,
		m.PairID, m.PlaceID, m.Time,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *MeetingRepo) GetMeetingByID(ctx context.Context, id int64) (*domain.Meeting, error) {
	var m domain.Meeting
	var dillConf, doeConf, dillCanc, doeCanc int

	err := r.db.QueryRowContext(ctx, `
		SELECT id, pair_id, place_id, time, dill_confirmed, doe_confirmed, dill_cancelled, doe_cancelled
		FROM meetings WHERE id = ?`, id).Scan(
		&m.ID, &m.PairID, &m.PlaceID, &m.Time,
		&dillConf, &doeConf, &dillCanc, &doeCanc,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	m.DillConfirmed = dillConf == 1
	m.DoeConfirmed = doeConf == 1
	m.DillCancelled = dillCanc == 1
	m.DoeCancelled = doeCanc == 1

	return &m, nil
}

func (r *MeetingRepo) ConfirmMeeting(ctx context.Context, meetingID int64, isDill bool) error {
	col := "doe_confirmed"
	if isDill {
		col = "dill_confirmed"
	}
	_, err := r.db.ExecContext(ctx,
		`UPDATE meetings SET `+col+` = 1 WHERE id = ?`, meetingID)
	return err
}

func (r *MeetingRepo) CancelMeeting(ctx context.Context, meetingID int64, isDill bool) error {
	col := "doe_cancelled"
	if isDill {
		col = "dill_cancelled"
	}
	_, err := r.db.ExecContext(ctx,
		`UPDATE meetings SET `+col+` = 1 WHERE id = ?`, meetingID)
	return err
}
