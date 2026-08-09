package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"gopa/internal/core/domain"
)

type VehicleRepository struct{ db *sqlx.DB }

func NewVehicleRepository(db *sqlx.DB) *VehicleRepository { return &VehicleRepository{db: db} }
func (r *VehicleRepository) List(ctx context.Context, userID uuid.UUID) ([]domain.Vehicle, error) {
	vehicles := make([]domain.Vehicle, 0)
	if err := r.db.SelectContext(ctx, &vehicles, `SELECT id, user_id, make, model, year, nickname, current_mileage, created_at, updated_at FROM vehicles WHERE user_id = $1 ORDER BY created_at DESC`, userID); err != nil {
		return nil, fmt.Errorf("list vehicles: %w", err)
	}
	return vehicles, nil
}
func (r *VehicleRepository) Get(ctx context.Context, userID, id uuid.UUID) (domain.Vehicle, error) {
	var vehicle domain.Vehicle
	if err := r.db.GetContext(ctx, &vehicle, `SELECT id, user_id, make, model, year, nickname, current_mileage, created_at, updated_at FROM vehicles WHERE id = $1 AND user_id = $2`, id, userID); err != nil {
		return domain.Vehicle{}, mapNotFound(err, "get vehicle")
	}
	return vehicle, nil
}
func (r *VehicleRepository) Create(ctx context.Context, vehicle domain.Vehicle) (domain.Vehicle, error) {
	if _, err := r.db.NamedExecContext(ctx, `INSERT INTO vehicles (id, user_id, make, model, year, nickname, current_mileage, created_at, updated_at) VALUES (:id, :user_id, :make, :model, :year, :nickname, :current_mileage, :created_at, :updated_at)`, vehicle); err != nil {
		return domain.Vehicle{}, fmt.Errorf("create vehicle: %w", err)
	}
	return vehicle, nil
}
func (r *VehicleRepository) Update(ctx context.Context, vehicle domain.Vehicle) (domain.Vehicle, error) {
	result, err := r.db.NamedExecContext(ctx, `UPDATE vehicles SET make = :make, model = :model, year = :year, nickname = :nickname, current_mileage = :current_mileage, updated_at = :updated_at WHERE id = :id AND user_id = :user_id`, vehicle)
	if err != nil {
		return domain.Vehicle{}, fmt.Errorf("update vehicle: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return domain.Vehicle{}, fmt.Errorf("count updated vehicle: %w", err)
	}
	if rows == 0 {
		return domain.Vehicle{}, domain.ErrNotFound
	}
	return vehicle, nil
}
func (r *VehicleRepository) ListLogs(ctx context.Context, userID, vehicleID uuid.UUID) ([]domain.VehicleLog, error) {
	logs := make([]domain.VehicleLog, 0)
	if err := r.db.SelectContext(ctx, &logs, `SELECT id, vehicle_id, user_id, log_type, mileage, description, cost::text AS cost, log_date, created_at, updated_at FROM vehicle_logs WHERE vehicle_id = $1 AND user_id = $2 ORDER BY log_date DESC, created_at DESC`, vehicleID, userID); err != nil {
		return nil, fmt.Errorf("list vehicle logs: %w", err)
	}
	return logs, nil
}
func (r *VehicleRepository) CreateLog(ctx context.Context, log domain.VehicleLog) (domain.VehicleLog, error) {
	var vehicleExists bool
	if err := r.db.GetContext(ctx, &vehicleExists, `SELECT EXISTS(SELECT 1 FROM vehicles WHERE id = $1 AND user_id = $2)`, log.VehicleID, log.UserID); err != nil {
		return domain.VehicleLog{}, fmt.Errorf("verify vehicle ownership: %w", err)
	}
	if !vehicleExists {
		return domain.VehicleLog{}, domain.ErrNotFound
	}
	if _, err := r.db.NamedExecContext(ctx, `INSERT INTO vehicle_logs (id, vehicle_id, user_id, log_type, mileage, description, cost, log_date, created_at, updated_at) VALUES (:id, :vehicle_id, :user_id, :log_type, :mileage, :description, :cost, :log_date, :created_at, :updated_at)`, log); err != nil {
		return domain.VehicleLog{}, fmt.Errorf("create vehicle log: %w", err)
	}
	if _, err := r.db.ExecContext(ctx, `UPDATE vehicles SET current_mileage = GREATEST(current_mileage, $1), updated_at = $2 WHERE id = $3 AND user_id = $4`, log.Mileage, log.UpdatedAt, log.VehicleID, log.UserID); err != nil {
		return domain.VehicleLog{}, fmt.Errorf("update vehicle mileage: %w", err)
	}
	return log, nil
}
func (r *VehicleRepository) LatestMaintenanceMileage(ctx context.Context, userID, vehicleID uuid.UUID) (*int, error) {
	var mileage *int
	if err := r.db.GetContext(ctx, &mileage, `SELECT mileage FROM vehicle_logs WHERE vehicle_id = $1 AND user_id = $2 AND log_type = 'MAINTENANCE' ORDER BY mileage DESC LIMIT 1`, vehicleID, userID); err != nil {
		if mapNotFound(err, "latest maintenance") == domain.ErrNotFound {
			return nil, nil
		}
		return nil, mapNotFound(err, "latest maintenance")
	}
	return mileage, nil
}
