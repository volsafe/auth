package storage

import (
	"context"
	"database/sql"
	"auth/db"
)

const columns = ""
type Profile struct {
	UserID             uint
	PhoneNumber        string
	Latitude           float64
	Longitude          float64
	HomeSize           float64
	BuildingFloors     int
	FloorLivingOn      int
	WindowModel        string
	AdultsCount        int
	ChildrenCount      int
	ElectricityCompany string
	MeterType          string
	BillNumber         string
}

type Storage struct {
	db *db.DB
}

func (s *Storage) Close() error {
	if s.db != nil {
		s.db.Close()
	}
	return nil
}

func NewStorage() (*Storage, error) {
	dbConn, err := db.NewDB()
	if err != nil {
		return nil, err
	}
	return &Storage{db: dbConn}, nil
}

func (s *Storage) CreateUserProfile(ctx context.Context, p Profile) error {
	query := `INSERT INTO profiles (` + columns + `)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	_, err := s.db.DB.ExecContext(ctx, query, p.UserID, p.PhoneNumber, p.Latitude, p.Longitude, p.HomeSize, p.BuildingFloors, p.FloorLivingOn, p.WindowModel, p.AdultsCount, p.ChildrenCount, p.ElectricityCompany, p.MeterType, p.BillNumber)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetUserProfile(ctx context.Context, userID uint) (*Profile, error) {
	var p Profile
	query := "SELECT " + columns + " FROM profiles WHERE user_id = $1"

	row := s.db.DB.QueryRowContext(ctx, query, userID)
	err := row.Scan(&p.UserID, &p.PhoneNumber, &p.Latitude, &p.Longitude, &p.HomeSize, &p.BuildingFloors, &p.FloorLivingOn, &p.WindowModel, &p.AdultsCount, &p.ChildrenCount, &p.ElectricityCompany, &p.MeterType, &p.BillNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &p, nil
}

func (s *Storage) UpdateUserProfileByID(ctx context.Context, p Profile) error {
	query := `UPDATE profiles SET phone_number = $1, latitude = $2, longitude = $3, home_size = $4, building_floors = $5, floor_living_on = $6, window_model = $7, adults_count = $8, children_count = $9, electricity_company = $10, meter_type = $11, bill_number = $12
              WHERE user_id = $13`
	_, err := s.db.DB.ExecContext(ctx, query, p.PhoneNumber, p.Latitude, p.Longitude, p.HomeSize, p.BuildingFloors, p.FloorLivingOn, p.WindowModel, p.AdultsCount, p.ChildrenCount, p.ElectricityCompany, p.MeterType, p.BillNumber, p.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteUserProfile(ctx context.Context, userID uint) error {
	query := `DELETE FROM profiles WHERE user_id = $1`
	_, err := s.db.DB.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}
