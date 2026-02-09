package database

import (
	"context"
	"fmt"
	"log"
	"time"
)

// VerifySchema проверяет наличие всех необходимых таблиц
func (d *Database) VerifySchema(ctx context.Context) error {
	requiredTables := []string{
		"users",
		"resources",
		"bookings",
		"sessions",
		"audit_logs",
		"notifications",
	}

	for _, table := range requiredTables {
		exists, err := d.tableExists(ctx, table)
		if err != nil {
			return fmt.Errorf("ошибка проверки таблицы %s: %w", table, err)
		}
		if !exists {
			return fmt.Errorf("таблица %s не найдена в БД", table)
		}
	}

	log.Println("✓ Все необходимые таблицы найдены в БД")
	return nil
}

// tableExists проверяет существование таблицы
func (d *Database) tableExists(ctx context.Context, tableName string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_schema = 'public'
			AND table_name = $1
		);
	`

	var exists bool
	err := d.DB.QueryRowContext(ctx, query, tableName).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// GetDatabaseStats возвращает статистику по БД
func (d *Database) GetDatabaseStats(ctx context.Context) (*DatabaseStats, error) {
	stats := &DatabaseStats{}

	// Считаем количество записей в каждой таблице
	tables := map[string]*int64{
		"users":         &stats.UsersCount,
		"resources":     &stats.ResourcesCount,
		"bookings":      &stats.BookingsCount,
		"sessions":      &stats.SessionsCount,
		"audit_logs":    &stats.AuditLogsCount,
		"notifications": &stats.NotificationsCount,
	}

	for table, count := range tables {
		query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
		err := d.DB.QueryRowContext(ctx, query).Scan(count)
		if err != nil {
			// Если таблицы нет, просто пропускаем
			*count = 0
		}
	}

	return stats, nil
}

// DatabaseStats содержит статистику по БД
type DatabaseStats struct {
	UsersCount         int64
	ResourcesCount     int64
	BookingsCount      int64
	SessionsCount      int64
	AuditLogsCount     int64
	NotificationsCount int64
}

// CleanupExpiredSessions удаляет истекшие сессии
func (d *Database) CleanupExpiredSessions(ctx context.Context) (int64, error) {
	query := `DELETE FROM sessions WHERE expires_at < $1`
	result, err := d.DB.ExecContext(ctx, query, time.Now())
	if err != nil {
		return 0, fmt.Errorf("ошибка удаления сессий: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rows, nil
}

// HealthCheck проверяет здоровье БД
func (d *Database) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Проверяем подключение
	if err := d.DB.PingContext(ctx); err != nil {
		return fmt.Errorf("БД недоступна: %w", err)
	}

	// Проверяем что можем выполнить простой запрос
	var result int
	err := d.DB.QueryRowContext(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса: %w", err)
	}

	return nil
}
