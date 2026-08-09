package services

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"gopa/internal/core/domain"
	"gopa/internal/core/ports"
)

type VehicleInput struct {
	Make, Model    string
	Year           int16
	Nickname       *string
	CurrentMileage int
}
type VehicleLogInput struct {
	LogType           domain.VehicleLogType
	Mileage           int
	Description, Cost string
	LogDate           time.Time
}
type AssetService struct {
	vehicles ports.VehicleRepository
	nodes    ports.NetworkRepository
	now      func() time.Time
}

func NewAssetService(vehicles ports.VehicleRepository, nodes ports.NetworkRepository) *AssetService {
	return &AssetService{vehicles: vehicles, nodes: nodes, now: time.Now}
}
func (s *AssetService) ListVehicles(ctx context.Context, userID uuid.UUID) ([]domain.Vehicle, error) {
	return s.vehicles.List(ctx, userID)
}
func (s *AssetService) GetVehicle(ctx context.Context, userID, id uuid.UUID) (domain.Vehicle, error) {
	return s.vehicles.Get(ctx, userID, id)
}
func (s *AssetService) CreateVehicle(ctx context.Context, userID uuid.UUID, input VehicleInput) (domain.Vehicle, error) {
	now := s.now().UTC()
	vehicle := domain.Vehicle{ID: uuid.New(), UserID: userID, Make: strings.TrimSpace(input.Make), Model: strings.TrimSpace(input.Model), Year: input.Year, Nickname: input.Nickname, CurrentMileage: input.CurrentMileage, CreatedAt: now, UpdatedAt: now}
	if err := vehicle.Validate(); err != nil {
		return domain.Vehicle{}, err
	}
	return s.vehicles.Create(ctx, vehicle)
}
func (s *AssetService) UpdateVehicle(ctx context.Context, userID, id uuid.UUID, input VehicleInput) (domain.Vehicle, error) {
	vehicle, err := s.vehicles.Get(ctx, userID, id)
	if err != nil {
		return domain.Vehicle{}, err
	}
	vehicle.Make = strings.TrimSpace(input.Make)
	vehicle.Model = strings.TrimSpace(input.Model)
	vehicle.Year = input.Year
	vehicle.Nickname = input.Nickname
	vehicle.CurrentMileage = input.CurrentMileage
	vehicle.UpdatedAt = s.now().UTC()
	if err := vehicle.Validate(); err != nil {
		return domain.Vehicle{}, err
	}
	return s.vehicles.Update(ctx, vehicle)
}
func (s *AssetService) ListVehicleLogs(ctx context.Context, userID, vehicleID uuid.UUID) ([]domain.VehicleLog, error) {
	if _, err := s.vehicles.Get(ctx, userID, vehicleID); err != nil {
		return nil, err
	}
	return s.vehicles.ListLogs(ctx, userID, vehicleID)
}
func (s *AssetService) CreateVehicleLog(ctx context.Context, userID, vehicleID uuid.UUID, input VehicleLogInput) (domain.VehicleLog, error) {
	now := s.now().UTC()
	log := domain.VehicleLog{ID: uuid.New(), VehicleID: vehicleID, UserID: userID, LogType: input.LogType, Mileage: input.Mileage, Description: strings.TrimSpace(input.Description), Cost: input.Cost, LogDate: input.LogDate.UTC(), CreatedAt: now, UpdatedAt: now}
	if err := log.Validate(); err != nil {
		return domain.VehicleLog{}, err
	}
	if err := domain.ParseCost(log.Cost); err != nil {
		return domain.VehicleLog{}, err
	}
	return s.vehicles.CreateLog(ctx, log)
}
func (s *AssetService) Maintenance(ctx context.Context, userID, vehicleID uuid.UUID) (domain.MaintenanceResult, error) {
	vehicle, err := s.vehicles.Get(ctx, userID, vehicleID)
	if err != nil {
		return domain.MaintenanceResult{}, err
	}
	lastMileage, err := s.vehicles.LatestMaintenanceMileage(ctx, userID, vehicleID)
	if err != nil {
		return domain.MaintenanceResult{}, err
	}
	return domain.CalculateMaintenance(vehicle.CurrentMileage, lastMileage), nil
}
func (s *AssetService) ListNodes(ctx context.Context, userID uuid.UUID) ([]domain.NetworkNode, error) {
	return s.nodes.List(ctx, userID)
}
func (s *AssetService) CreateNode(ctx context.Context, userID uuid.UUID, node domain.NetworkNode) (domain.NetworkNode, error) {
	now := s.now().UTC()
	node.ID = uuid.New()
	node.UserID = userID
	node.CreatedAt = now
	node.UpdatedAt = now
	if err := node.Validate(); err != nil {
		return domain.NetworkNode{}, err
	}
	if err := s.validateParent(ctx, userID, node.ID, node.ParentNodeID); err != nil {
		return domain.NetworkNode{}, err
	}
	return s.nodes.Create(ctx, node)
}
func (s *AssetService) UpdateNode(ctx context.Context, userID, id uuid.UUID, input domain.NetworkNode) (domain.NetworkNode, error) {
	node, err := s.nodes.Get(ctx, userID, id)
	if err != nil {
		return domain.NetworkNode{}, err
	}
	node.DeviceName = strings.TrimSpace(input.DeviceName)
	node.MACAddress = input.MACAddress
	node.StaticIP = input.StaticIP
	node.VLAN = input.VLAN
	node.ConnectionType = input.ConnectionType
	node.ParentNodeID = input.ParentNodeID
	node.Location = input.Location
	node.Notes = input.Notes
	node.UpdatedAt = s.now().UTC()
	if err := node.Validate(); err != nil {
		return domain.NetworkNode{}, err
	}
	if err := s.validateParent(ctx, userID, node.ID, node.ParentNodeID); err != nil {
		return domain.NetworkNode{}, err
	}
	return s.nodes.Update(ctx, node)
}
func (s *AssetService) DeleteNode(ctx context.Context, userID, id uuid.UUID) error {
	return s.nodes.Delete(ctx, userID, id)
}
func (s *AssetService) validateParent(ctx context.Context, userID, nodeID uuid.UUID, parentID *uuid.UUID) error {
	visited := map[uuid.UUID]bool{nodeID: true}
	current := parentID
	for current != nil {
		if visited[*current] {
			return domain.ErrValidation
		}
		visited[*current] = true
		parent, err := s.nodes.Get(ctx, userID, *current)
		if err != nil {
			return err
		}
		current = parent.ParentNodeID
	}
	return nil
}
