package entities

import "testing"

func TestCost_TableName(t *testing.T) {
	cost := Cost{}
	if cost.TableName() != "costs" {
		t.Errorf("Expected table name 'costs', got '%s'", cost.TableName())
	}
}

func TestCost_BeforeSave_ValidCost(t *testing.T) {
	projectID := uint(1)
	cost := Cost{
		ProjectID: &projectID,
		Type:      CostTypeMaterial,
		Name:      "Materials cost",
		Amount:    1000,
		Quantity:  1,
	}

	err := cost.BeforeSave(nil)
	if err != nil {
		t.Errorf("Expected no error for valid cost, got: %v", err)
	}
}

func TestCost_BeforeSave_NoAssociation(t *testing.T) {
	cost := Cost{
		Type:   CostTypeLabor,
		Name:   "Development cost",
		Amount: 1000,
	}

	err := cost.BeforeSave(nil)
	if err == nil {
		t.Error("Expected error for cost with no project/milestone/task association, got nil")
	}
}

func TestCost_BeforeSave_InvalidCostType(t *testing.T) {
	projectID := uint(1)
	cost := Cost{
		ProjectID: &projectID,
		Type:      CostType("invalid"),
		Name:      "Development cost",
		Amount:    1000,
	}

	err := cost.BeforeSave(nil)
	if err == nil {
		t.Error("Expected error for invalid cost type, got nil")
	}
}

func TestCost_BeforeSave_NegativeValues(t *testing.T) {
	projectID := uint(1)
	tests := []struct {
		name     string
		amount   float64
		quantity float64
		unitCost float64
		hours    float64
		wantErr  bool
	}{
		{"Valid values", 1000, 1, 100, 8, false},
		{"Negative amount", -100, 1, 100, 8, true},
		{"Negative quantity", 1000, -1, 100, 8, true},
		{"Negative unit cost", 1000, 1, -100, 8, true},
		{"Negative hours", 1000, 1, 100, -8, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost := Cost{
				ProjectID: &projectID,
				Type:      CostTypeMaterial, // Use non-labor cost type for validation tests
				Name:      "Test cost",
				Amount:    tt.amount,
				Quantity:  tt.quantity,
				UnitCost:  tt.unitCost,
				Hours:     tt.hours,
			}

			err := cost.BeforeSave(nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("BeforeSave() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCost_BeforeSave_AutoCalculateAmount(t *testing.T) {
	projectID := uint(1)
	cost := Cost{
		ProjectID: &projectID,
		Type:      CostTypeMaterial,
		Name:      "Materials",
		Amount:    0,
		Quantity:  10,
		UnitCost:  50,
	}

	err := cost.BeforeSave(nil)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	expected := 500.0 // 10 * 50
	if cost.Amount != expected {
		t.Errorf("Expected amount to be auto-calculated to %v, got %v", expected, cost.Amount)
	}
}

func TestCost_BeforeSave_LaborCostValidation(t *testing.T) {
	projectID := uint(1)
	tests := []struct {
		name          string
		resourceID    *uint
		projectRoleID *uint
		wantErr       bool
	}{
		{"Labor with resource", uintPtr(1), nil, false},
		{"Labor with project role", nil, uintPtr(1), false},
		{"Labor with both", uintPtr(1), uintPtr(1), false},
		{"Labor with neither", nil, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost := Cost{
				ProjectID:     &projectID,
				Type:          CostTypeLabor,
				Name:          "Labor cost",
				Amount:        1000,
				ResourceID:    tt.resourceID,
				ProjectRoleID: tt.projectRoleID,
			}

			err := cost.BeforeSave(nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("BeforeSave() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCost_IsLabor(t *testing.T) {
	tests := []struct {
		name     string
		costType CostType
		want     bool
	}{
		{"Labor cost", CostTypeLabor, true},
		{"Material cost", CostTypeMaterial, false},
		{"Equipment cost", CostTypeEquipment, false},
		{"Infrastructure cost", CostTypeInfrastructure, false},
		{"Service cost", CostTypeService, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost := Cost{Type: tt.costType}
			if got := cost.IsLabor(); got != tt.want {
				t.Errorf("IsLabor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCost_IsActual(t *testing.T) {
	tests := []struct {
		name        string
		isEstimated bool
		want        bool
	}{
		{"Estimated cost", true, false},
		{"Actual cost", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost := Cost{IsEstimated: tt.isEstimated}
			if got := cost.IsActual(); got != tt.want {
				t.Errorf("IsActual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCost_TotalAmount(t *testing.T) {
	tests := []struct {
		name     string
		amount   float64
		quantity float64
		unitCost float64
		want     float64
	}{
		{"With unit cost and quantity", 0, 10, 50, 500},
		{"Without unit cost", 1000, 1, 0, 1000},
		{"Without quantity", 1000, 0, 50, 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost := Cost{
				Amount:   tt.amount,
				Quantity: tt.quantity,
				UnitCost: tt.unitCost,
			}
			if got := cost.TotalAmount(); got != tt.want {
				t.Errorf("TotalAmount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCost_CalculateLaborCost(t *testing.T) {
	tests := []struct {
		name       string
		costType   CostType
		hours      float64
		rateType   RateType
		rate       float64
		want       float64
	}{
		{"Hourly rate", CostTypeLabor, 40, RateTypeHourly, 50, 2000},
		{"Daily rate", CostTypeLabor, 40, RateTypeDaily, 400, 2000}, // 40 hours = 5 days, 5 * 400 = 2000
		{"Monthly rate", CostTypeLabor, 160, RateTypeMonthly, 8000, 8000}, // 160 hours = 1 month
		{"Non-labor cost", CostTypeMaterial, 40, RateTypeHourly, 50, 0},
		{"Zero hours", CostTypeLabor, 0, RateTypeHourly, 50, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost := Cost{
				Type:     tt.costType,
				Hours:    tt.hours,
				RateType: tt.rateType,
			}
			if got := cost.CalculateLaborCost(tt.rate); got != tt.want {
				t.Errorf("CalculateLaborCost() = %v, want %v", got, tt.want)
			}
		})
	}
}
