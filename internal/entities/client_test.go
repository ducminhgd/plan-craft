package entities

import "testing"

func TestClient_TableName(t *testing.T) {
	client := Client{}
	if client.TableName() != "clients" {
		t.Errorf("Expected table name 'clients', got '%s'", client.TableName())
	}
}

func TestClient_BeforeSave_ValidClient(t *testing.T) {
	client := Client{
		Name:  "Test Client Inc.",
		Email: "contact@testclient.com",
	}

	err := client.BeforeSave(nil)
	if err != nil {
		t.Errorf("Expected no error for valid client, got: %v", err)
	}
}

func TestClient_BeforeSave_EmptyName(t *testing.T) {
	client := Client{
		Name:  "",
		Email: "contact@testclient.com",
	}

	err := client.BeforeSave(nil)
	if err == nil {
		t.Error("Expected error for empty client name, got nil")
	}
}

func TestClient_IsActiveClient(t *testing.T) {
	tests := []struct {
		name     string
		isActive bool
		want     bool
	}{
		{"Active client", true, true},
		{"Inactive client", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := Client{IsActive: tt.isActive}
			if got := client.IsActiveClient(); got != tt.want {
				t.Errorf("IsActiveClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
