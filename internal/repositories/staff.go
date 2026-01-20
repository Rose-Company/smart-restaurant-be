package repositories

import (
	"app-noti/internal/models"

	"gorm.io/gorm"
)

// StaffProfileRepository interface
type StaffProfileRepository interface {
	BaseRepository[models.StaffProfile]
}

type staffProfileRepository struct {
	*baseRepository[models.StaffProfile]
}

func NewStaffProfileRepository(db *gorm.DB) StaffProfileRepository {
	return &staffProfileRepository{
		baseRepository: &baseRepository[models.StaffProfile]{db: db},
	}
}

// WaiterTableAssignmentRepository interface
type WaiterTableAssignmentRepository interface {
	BaseRepository[models.WaiterTableAssignment]
}

type waiterTableAssignmentRepository struct {
	*baseRepository[models.WaiterTableAssignment]
}

func NewWaiterTableAssignmentRepository(db *gorm.DB) WaiterTableAssignmentRepository {
	return &waiterTableAssignmentRepository{
		baseRepository: &baseRepository[models.WaiterTableAssignment]{db: db},
	}
}

// StaffInvitationRepository interface
type StaffInvitationRepository interface {
	BaseRepository[models.StaffInvitation]
}

type staffInvitationRepository struct {
	*baseRepository[models.StaffInvitation]
}

func NewStaffInvitationRepository(db *gorm.DB) StaffInvitationRepository {
	return &staffInvitationRepository{
		baseRepository: &baseRepository[models.StaffInvitation]{db: db},
	}
}
