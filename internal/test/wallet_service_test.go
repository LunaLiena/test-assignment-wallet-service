package test

import (
	"context"
	"testing"
	"time"
	"wallet-service/internal/entities"
	"wallet-service/internal/service"

	"github.com/stretchr/testify/assert"
	testcontainers "github.com/testcontainers/testcontainers-go"
	postgresContainer "github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	ctx := context.Background()

	pgContainer, err := postgresContainer.Run(ctx,
		"postgres:15",
		postgresContainer.WithDatabase("testdb"),
		postgresContainer.WithUsername("testuser"),
		postgresContainer.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)

	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")

	if err != nil {
		t.Fatal(err)
	}

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	if err := db.AutoMigrate(&entities.Wallet{}); err != nil {
		t.Fatal("Failed to migrate test database:", err)
	}

	return db
}

func TestWalletService_Deposit(t *testing.T) {
	db := setupTestDB(t)
	service := service.NewWalletService(db)

	walletId := "550e8400-e29b-41d4-a716-446655440000"

	// 🔥 Добавьте инициализацию кошелька
	err := service.CreateWallet(walletId)
	assert.NoError(t, err)

	amount := 1000.0
	err = service.Deposit(walletId, amount)
	assert.NoError(t, err)

	balance, err := service.GetBalance(walletId)
	assert.NoError(t, err)
	assert.Equal(t, amount, balance)
}

func TestWalletService_Withdraw(t *testing.T) {
	db := setupTestDB(t)
	service := service.NewWalletService(db)

	walletID := "550e8400-e29b-41d4-a716-446655440001"
	depositAmount := 1000.0
	withdrawAmount := 300.0

	// 🔥 Обязательно создаём кошелек перед первой операцией
	err := service.CreateWallet(walletID)
	assert.NoError(t, err)

	// Сначала депозит
	err = service.Deposit(walletID, depositAmount)
	assert.NoError(t, err)

	// Затем списание
	err = service.WithDraw(walletID, withdrawAmount)
	assert.NoError(t, err)

	// Проверяем баланс
	balance, err := service.GetBalance(walletID)
	assert.NoError(t, err)
	assert.Equal(t, depositAmount-withdrawAmount, balance)
}
