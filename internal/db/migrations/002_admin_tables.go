package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var AdminTables = &gormigrate.Migration{
	ID: "20251118_002_admin_tables",
	Migrate: func(db *gorm.DB) error {
		// First create all tables without data
		sql := `
		-- Create tables if they don't exist
		CREATE TABLE IF NOT EXISTS goadmin_users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(100) UNIQUE NOT NULL,
			password VARCHAR(100) NOT NULL,
			name VARCHAR(100),
			avatar VARCHAR(255),
			remember_token VARCHAR(100),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS goadmin_roles (
			id SERIAL PRIMARY KEY,
			name VARCHAR(50) UNIQUE NOT NULL,
			slug VARCHAR(50) UNIQUE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS goadmin_permissions (
			id SERIAL PRIMARY KEY,
			name VARCHAR(50) NOT NULL,
			slug VARCHAR(50) UNIQUE NOT NULL,
			http_method VARCHAR(255),
			http_path TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS goadmin_role_users (
			role_id INT REFERENCES goadmin_roles(id) ON DELETE CASCADE,
			user_id INT REFERENCES goadmin_users(id) ON DELETE CASCADE,
			PRIMARY KEY (role_id, user_id)
		);

		CREATE TABLE IF NOT EXISTS goadmin_role_permissions (
			role_id INT REFERENCES goadmin_roles(id) ON DELETE CASCADE,
			permission_id INT REFERENCES goadmin_permissions(id) ON DELETE CASCADE,
			PRIMARY KEY (role_id, permission_id)
		);

		CREATE TABLE IF NOT EXISTS goadmin_user_permissions (
			user_id INT REFERENCES goadmin_users(id) ON DELETE CASCADE,
			permission_id INT REFERENCES goadmin_permissions(id) ON DELETE CASCADE,
			PRIMARY KEY (user_id, permission_id)
		);

		CREATE TABLE IF NOT EXISTS goadmin_menu (
			id SERIAL PRIMARY KEY,
			parent_id INT DEFAULT 0,
			type INT DEFAULT 0,
			ordering INT DEFAULT 0,
			title VARCHAR(50) NOT NULL,
			icon VARCHAR(50),
			uri VARCHAR(255),
			header VARCHAR(150),
			plugin_name VARCHAR(150) DEFAULT '',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS goadmin_role_menu (
			role_id INT NOT NULL,
			menu_id INT NOT NULL,
			PRIMARY KEY (role_id, menu_id)
		);

		CREATE TABLE IF NOT EXISTS goadmin_operation_log (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL,
			path VARCHAR(255) NOT NULL,
			method VARCHAR(10) NOT NULL,
			ip VARCHAR(15) NOT NULL,
			input TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS goadmin_session (
			id VARCHAR(100) PRIMARY KEY,
			"values" TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`

		// Execute table creation
		err := db.Exec(sql).Error
		if err != nil {
			return err
		}

		// Now insert data with proper error handling
		insertSQL := []string{
			// Add default admin user (password: admin, hashed with bcrypt)
			`INSERT INTO goadmin_users (username, password, name, avatar)
			SELECT 'admin', '$2a$10$M3xA3LK5aLzVyF5ykYJlYe7C3OPKbPqY8Kp1xKLgxK9pY2P5Q6M0e', 'Administrator', ''
			WHERE NOT EXISTS (SELECT 1 FROM goadmin_users WHERE username = 'admin');`,

			// Add admin role
			`INSERT INTO goadmin_roles (name, slug)
			SELECT 'Administrator', 'administrator'
			WHERE NOT EXISTS (SELECT 1 FROM goadmin_roles WHERE slug = 'administrator');`,

			// Add all permissions permission
			`INSERT INTO goadmin_permissions (name, slug, http_method, http_path)
			SELECT 'All permission', '*', '', '*'
			WHERE NOT EXISTS (SELECT 1 FROM goadmin_permissions WHERE slug = '*');`,

			// Link admin user to admin role
			`INSERT INTO goadmin_role_users (role_id, user_id)
			SELECT r.id, u.id 
			FROM goadmin_roles r, goadmin_users u 
			WHERE r.slug = 'administrator' AND u.username = 'admin'
			AND NOT EXISTS (
				SELECT 1 FROM goadmin_role_users ru 
				WHERE ru.role_id = r.id AND ru.user_id = u.id
			);`,

			// Link admin role to all permissions
			`INSERT INTO goadmin_role_permissions (role_id, permission_id)
			SELECT r.id, p.id 
			FROM goadmin_roles r, goadmin_permissions p 
			WHERE r.slug = 'administrator' AND p.slug = '*'
			AND NOT EXISTS (
				SELECT 1 FROM goadmin_role_permissions rp 
				WHERE rp.role_id = r.id AND rp.permission_id = p.id
			);`,

			// Add menu items
			`INSERT INTO goadmin_menu (parent_id, type, ordering, title, icon, uri)
			SELECT * FROM (VALUES 
				(0, 1, 1, 'Dashboard', 'fa-bar-chart', '/'),
				(0, 1, 2, 'Users', 'fa-users', '/info/users'),
				(0, 1, 3, 'Workflows', 'fa-project-diagram', '/info/workflows'),
				(0, 1, 4, 'Executions', 'fa-play-circle', '/info/executions'),
				(0, 1, 5, 'Connections', 'fa-plug', '/info/connections'),
				(0, 1, 6, 'Subscriptions', 'fa-credit-card', '/info/subscriptions')
			) AS menu_items (parent_id, type, ordering, title, icon, uri)
			WHERE NOT EXISTS (SELECT 1 FROM goadmin_menu m WHERE m.title = menu_items.title);`,

			// Link all menus to admin role
			`INSERT INTO goadmin_role_menu (role_id, menu_id)
			SELECT r.id, m.id 
			FROM goadmin_roles r, goadmin_menu m 
			WHERE r.slug = 'administrator'
			AND NOT EXISTS (
				SELECT 1 FROM goadmin_role_menu rm 
				WHERE rm.role_id = r.id AND rm.menu_id = m.id
			);`,

			// Add additional columns to existing tables if they don't exist
			`DO $$
			BEGIN
				IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
					WHERE table_name = 'users' AND column_name = 'subscription_plan') THEN
					ALTER TABLE users ADD COLUMN subscription_plan VARCHAR(20) DEFAULT 'free';
				END IF;

				IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
					WHERE table_name = 'workflows' AND column_name = 'total_executions') THEN
					ALTER TABLE workflows ADD COLUMN total_executions INT DEFAULT 0;
				END IF;

				IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
					WHERE table_name = 'workflows' AND column_name = 'success_count') THEN
					ALTER TABLE workflows ADD COLUMN success_count INT DEFAULT 0;
				END IF;

				IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
					WHERE table_name = 'workflows' AND column_name = 'error_count') THEN
					ALTER TABLE workflows ADD COLUMN error_count INT DEFAULT 0;
				END IF;

				IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
					WHERE table_name = 'workflows' AND column_name = 'trigger_type') THEN
					ALTER TABLE workflows ADD COLUMN trigger_type VARCHAR(20);
				END IF;

				IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
					WHERE table_name = 'executions' AND column_name = 'duration_seconds') THEN
					ALTER TABLE executions ADD COLUMN duration_seconds INT;
				END IF;
			END $$;`,
		}

		// Execute all insert statements
		for _, stmt := range insertSQL {
			err = db.Exec(stmt).Error
			if err != nil {
				return err
			}
		}

		return nil
	},
	Rollback: func(db *gorm.DB) error {
		// We don't want to drop tables in rollback to avoid data loss
		return nil
	},
}
