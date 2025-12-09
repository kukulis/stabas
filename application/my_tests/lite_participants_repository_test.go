package my_tests

import (
	"darbelis.eu/stabas/dao"
	"darbelis.eu/stabas/db"
	"darbelis.eu/stabas/entities"
	"testing"
)

func setupLiteTestDB(t *testing.T) *dao.LiteParticipantsRepository {
	t.Helper()
	database := db.NewDatabase(":memory:")
	repo, err := dao.NewLiteParticipantsRepository(database)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	return repo
}

func TestNewLiteParticipantsRepository(t *testing.T) {
	database := db.NewDatabase(":memory:")
	repo, err := dao.NewLiteParticipantsRepository(database)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if repo == nil {
		t.Fatal("Expected repository to be created, got nil")
	}

	defer func() { _ = repo.Close() }()
}

func TestLiteAddParticipant(t *testing.T) {
	repo := setupLiteTestDB(t)
	defer func() { _ = repo.Close() }()

	participant := &entities.Participant{
		Name:    "John Doe",
		Deleted: false,
	}

	added, err := repo.AddParticipant(participant)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if added.Id == 0 {
		t.Error("Expected participant ID to be set")
	}

	if added.Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", added.Name)
	}
}

func TestLiteGetParticipants(t *testing.T) {
	repo := setupLiteTestDB(t)
	defer func() { _ = repo.Close() }()

	_, _ = repo.AddParticipant(&entities.Participant{Name: "Alice"})
	_, _ = repo.AddParticipant(&entities.Participant{Name: "Bob"})
	_, _ = repo.AddParticipant(&entities.Participant{Name: "Charlie"})

	participants := repo.GetParticipants()

	if len(participants) != 3 {
		t.Errorf("Expected 3 participants, got %d", len(participants))
	}

	names := []string{"Alice", "Bob", "Charlie"}
	for i, p := range participants {
		if p.Name != names[i] {
			t.Errorf("Expected participant %d to be '%s', got '%s'", i, names[i], p.Name)
		}
	}
}

func TestLiteGetParticipantsExcludesDeleted(t *testing.T) {
	repo := setupLiteTestDB(t)
	defer func() { _ = repo.Close() }()

	p1, _ := repo.AddParticipant(&entities.Participant{Name: "Alice"})
	_, _ = repo.AddParticipant(&entities.Participant{Name: "Bob"})
	_, _ = repo.AddParticipant(&entities.Participant{Name: "Charlie"})

	_ = repo.RemoveParticipant(p1.Id)

	participants := repo.GetParticipants()

	if len(participants) != 2 {
		t.Errorf("Expected 2 participants (excluding deleted), got %d", len(participants))
	}

	for _, p := range participants {
		if p.Name == "Alice" {
			t.Error("Deleted participant 'Alice' should not be in results")
		}
	}
}

func TestLiteFindParticipant(t *testing.T) {
	repo := setupLiteTestDB(t)
	defer func() { _ = repo.Close() }()

	added, _ := repo.AddParticipant(&entities.Participant{Name: "John Doe"})

	found, err := repo.FindParticipant(added.Id)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if found.Id != added.Id {
		t.Errorf("Expected ID %d, got %d", added.Id, found.Id)
	}

	if found.Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", found.Name)
	}
}

func TestLiteFindParticipantNotFound(t *testing.T) {
	repo := setupLiteTestDB(t)
	defer func() { _ = repo.Close() }()

	_, err := repo.FindParticipant(999)

	if err == nil {
		t.Error("Expected error for non-existent participant, got nil")
	}
}

func TestLiteFindParticipantIncludesDeleted(t *testing.T) {
	repo := setupLiteTestDB(t)
	defer func() { _ = repo.Close() }()

	added, _ := repo.AddParticipant(&entities.Participant{Name: "John Doe"})
	_ = repo.RemoveParticipant(added.Id)

	found, err := repo.FindParticipant(added.Id)

	if err != nil {
		t.Fatalf("Expected to find deleted participant, got error: %v", err)
	}

	if !found.Deleted {
		t.Error("Expected participant to be marked as deleted")
	}
}

func TestLiteUpdateParticipant(t *testing.T) {
	repo := setupLiteTestDB(t)
	defer func() { _ = repo.Close() }()

	added, _ := repo.AddParticipant(&entities.Participant{Name: "John Doe"})

	added.Name = "Jane Smith"
	err := repo.UpdateParticipant(added)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	found, _ := repo.FindParticipant(added.Id)

	if found.Name != "Jane Smith" {
		t.Errorf("Expected name 'Jane Smith', got '%s'", found.Name)
	}
}

func TestLiteUpdateParticipantNotFound(t *testing.T) {
	repo := setupLiteTestDB(t)
	defer func() { _ = repo.Close() }()

	err := repo.UpdateParticipant(&entities.Participant{Id: 999, Name: "Ghost"})

	if err == nil {
		t.Error("Expected error when updating non-existent participant, got nil")
	}
}

func TestLiteRemoveParticipant(t *testing.T) {
	repo := setupLiteTestDB(t)
	defer func() { _ = repo.Close() }()

	added, _ := repo.AddParticipant(&entities.Participant{Name: "John Doe"})

	err := repo.RemoveParticipant(added.Id)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	participants := repo.GetParticipants()
	if len(participants) != 0 {
		t.Errorf("Expected 0 active participants, got %d", len(participants))
	}

	found, _ := repo.FindParticipant(added.Id)
	if !found.Deleted {
		t.Error("Expected participant to be marked as deleted")
	}
}

func TestLiteRemoveParticipantNotFound(t *testing.T) {
	repo := setupLiteTestDB(t)
	defer func() { _ = repo.Close() }()

	err := repo.RemoveParticipant(999)

	if err == nil {
		t.Error("Expected error when removing non-existent participant, got nil")
	}
}

func TestLiteFindParticipantByName(t *testing.T) {
	repo := setupLiteTestDB(t)
	defer func() { _ = repo.Close() }()

	_, _ = repo.AddParticipant(&entities.Participant{Name: "Alice"})
	_, _ = repo.AddParticipant(&entities.Participant{Name: "Bob"})

	found := repo.FindParticipantByName("Alice")

	if found == nil {
		t.Fatal("Expected to find participant 'Alice', got nil")
	}

	if found.Name != "Alice" {
		t.Errorf("Expected name 'Alice', got '%s'", found.Name)
	}
}

func TestLiteFindParticipantByNameNotFound(t *testing.T) {
	repo := setupLiteTestDB(t)
	defer func() { _ = repo.Close() }()

	_, _ = repo.AddParticipant(&entities.Participant{Name: "Alice"})

	found := repo.FindParticipantByName("NonExistent")

	if found != nil {
		t.Error("Expected nil for non-existent participant, got result")
	}
}

func TestLiteFindParticipantByNameExcludesDeleted(t *testing.T) {
	repo := setupLiteTestDB(t)
	defer func() { _ = repo.Close() }()

	added, _ := repo.AddParticipant(&entities.Participant{Name: "Alice"})
	_ = repo.RemoveParticipant(added.Id)

	found := repo.FindParticipantByName("Alice")

	if found != nil {
		t.Error("Expected nil for deleted participant, got result")
	}
}

func TestLiteUpdateParticipantToken(t *testing.T) {
	repo := setupLiteTestDB(t)
	defer func() { _ = repo.Close() }()

	added, _ := repo.AddParticipant(&entities.Participant{Name: "John Doe"})

	token := "test-token-123"
	err := repo.UpdateParticipantToken(added.Id, token)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	found, _ := repo.FindParticipant(added.Id)

	if found.Token != token {
		t.Errorf("Expected token '%s', got '%s'", token, found.Token)
	}
}

func TestLiteUpdateParticipantTokenNotFound(t *testing.T) {
	repo := setupLiteTestDB(t)
	defer func() { _ = repo.Close() }()

	err := repo.UpdateParticipantToken(999, "test-token")

	if err == nil {
		t.Error("Expected error when updating token for non-existent participant, got nil")
	}
}

func TestLiteFindParticipantByToken(t *testing.T) {
	repo := setupLiteTestDB(t)
	defer func() { _ = repo.Close() }()

	added, _ := repo.AddParticipant(&entities.Participant{Name: "John Doe"})
	token := "test-token-456"
	_ = repo.UpdateParticipantToken(added.Id, token)

	found := repo.FindParticipantByToken(token)

	if found == nil {
		t.Fatal("Expected to find participant by token, got nil")
	}

	if found.Id != added.Id {
		t.Errorf("Expected ID %d, got %d", added.Id, found.Id)
	}

	if found.Token != token {
		t.Errorf("Expected token '%s', got '%s'", token, found.Token)
	}
}

func TestLiteFindParticipantByTokenNotFound(t *testing.T) {
	repo := setupLiteTestDB(t)
	defer func() { _ = repo.Close() }()

	found := repo.FindParticipantByToken("non-existent-token")

	if found != nil {
		t.Error("Expected nil for non-existent token, got result")
	}
}

func TestLiteFindParticipantByTokenExcludesDeleted(t *testing.T) {
	repo := setupLiteTestDB(t)
	defer func() { _ = repo.Close() }()

	added, _ := repo.AddParticipant(&entities.Participant{Name: "John Doe"})
	token := "test-token-789"
	_ = repo.UpdateParticipantToken(added.Id, token)
	_ = repo.RemoveParticipant(added.Id)

	found := repo.FindParticipantByToken(token)

	if found != nil {
		t.Error("Expected nil for deleted participant, got result")
	}
}

func TestLiteUpdateParticipantPassword(t *testing.T) {
	repo := setupLiteTestDB(t)
	defer func() { _ = repo.Close() }()

	added, _ := repo.AddParticipant(&entities.Participant{Name: "John Doe"})

	password := "securepassword123"
	err := repo.UpdateParticipantPassword(added.Id, password)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	found, _ := repo.FindParticipant(added.Id)

	if found.Password != password {
		t.Errorf("Expected password '%s', got '%s'", password, found.Password)
	}
}

func TestLiteUpdateParticipantPasswordNotFound(t *testing.T) {
	repo := setupLiteTestDB(t)
	defer func() { _ = repo.Close() }()

	err := repo.UpdateParticipantPassword(999, "password")

	if err == nil {
		t.Error("Expected error when updating password for non-existent participant, got nil")
	}
}

func TestLiteMultipleOperations(t *testing.T) {
	repo := setupLiteTestDB(t)
	defer func() { _ = repo.Close() }()

	p1, _ := repo.AddParticipant(&entities.Participant{Name: "Alice"})
	p2, _ := repo.AddParticipant(&entities.Participant{Name: "Bob"})
	p3, _ := repo.AddParticipant(&entities.Participant{Name: "Charlie"})

	_ = repo.UpdateParticipantToken(p1.Id, "alice-token")
	_ = repo.UpdateParticipantPassword(p1.Id, "alice-pass")

	_ = repo.UpdateParticipant(&entities.Participant{Id: p2.Id, Name: "Bob Updated"})

	_ = repo.RemoveParticipant(p3.Id)

	participants := repo.GetParticipants()
	if len(participants) != 2 {
		t.Errorf("Expected 2 active participants, got %d", len(participants))
	}

	foundByToken := repo.FindParticipantByToken("alice-token")
	if foundByToken == nil || foundByToken.Name != "Alice" {
		t.Error("Failed to find Alice by token")
	}

	foundByName := repo.FindParticipantByName("Bob Updated")
	if foundByName == nil {
		t.Error("Failed to find updated Bob by name")
	}
}
