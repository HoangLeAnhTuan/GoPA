package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"gopa/internal/core/domain"
)

type NetworkRepository struct{ db *sqlx.DB }

func NewNetworkRepository(db *sqlx.DB) *NetworkRepository { return &NetworkRepository{db: db} }
func (r *NetworkRepository) List(ctx context.Context, userID uuid.UUID) ([]domain.NetworkNode, error) {
	nodes := make([]domain.NetworkNode, 0)
	if err := r.db.SelectContext(ctx, &nodes, `SELECT id, user_id, device_name, mac_address, static_ip, vlan_tag, connection_type, parent_node_id, location, notes, created_at, updated_at FROM network_nodes WHERE user_id = $1 ORDER BY created_at`, userID); err != nil {
		return nil, fmt.Errorf("list network nodes: %w", err)
	}
	return nodes, nil
}
func (r *NetworkRepository) Get(ctx context.Context, userID, id uuid.UUID) (domain.NetworkNode, error) {
	var node domain.NetworkNode
	if err := r.db.GetContext(ctx, &node, `SELECT id, user_id, device_name, mac_address, static_ip, vlan_tag, connection_type, parent_node_id, location, notes, created_at, updated_at FROM network_nodes WHERE id = $1 AND user_id = $2`, id, userID); err != nil {
		return domain.NetworkNode{}, mapNotFound(err, "get network node")
	}
	return node, nil
}
func (r *NetworkRepository) Create(ctx context.Context, node domain.NetworkNode) (domain.NetworkNode, error) {
	if _, err := r.db.NamedExecContext(ctx, `INSERT INTO network_nodes (id, user_id, device_name, mac_address, static_ip, vlan_tag, connection_type, parent_node_id, location, notes, created_at, updated_at) VALUES (:id, :user_id, :device_name, :mac_address, :static_ip, :vlan_tag, :connection_type, :parent_node_id, :location, :notes, :created_at, :updated_at)`, node); err != nil {
		return domain.NetworkNode{}, fmt.Errorf("create network node: %w", err)
	}
	return node, nil
}
func (r *NetworkRepository) Update(ctx context.Context, node domain.NetworkNode) (domain.NetworkNode, error) {
	result, err := r.db.NamedExecContext(ctx, `UPDATE network_nodes SET device_name = :device_name, mac_address = :mac_address, static_ip = :static_ip, vlan_tag = :vlan_tag, connection_type = :connection_type, parent_node_id = :parent_node_id, location = :location, notes = :notes, updated_at = :updated_at WHERE id = :id AND user_id = :user_id`, node)
	if err != nil {
		return domain.NetworkNode{}, fmt.Errorf("update network node: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return domain.NetworkNode{}, fmt.Errorf("count updated network node: %w", err)
	}
	if rows == 0 {
		return domain.NetworkNode{}, domain.ErrNotFound
	}
	return node, nil
}
func (r *NetworkRepository) Delete(ctx context.Context, userID, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM network_nodes WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("delete network node: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("count deleted network node: %w", err)
	}
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}
