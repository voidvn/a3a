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
    values TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- Добавление дефолтного админа (пароль: admin, хешированный bcrypt)
INSERT INTO goadmin_users (username, password, name, avatar)
VALUES ('admin', '$2a$10$M3xA3LK5aLzVyF5ykYJlYe7C3OPKbPqY8Kp1xKLgxK9pY2P5Q6M0e', 'Администратор', '');

-- Создание роли администратора
INSERT INTO goadmin_roles (name, slug) VALUES ('Administrator', 'administrator');

-- Связь пользователя с ролью
INSERT INTO goadmin_role_users (role_id, user_id) VALUES (1, 1);

-- Создание пермишенов для всех
INSERT INTO goadmin_permissions (name, slug, http_method, http_path)
VALUES ('All permission', '*', '', '*');

-- Связь роли с пермишеном
INSERT INTO goadmin_role_permissions (role_id, permission_id) VALUES (1, 1);

-- Создание меню
INSERT INTO goadmin_menu (parent_id, type, ordering, title, icon, uri) VALUES
                                                                           (0, 1, 1, 'Dashboard', 'fa-bar-chart', '/'),
                                                                           (0, 1, 2, 'Пользователи', 'fa-users', '/info/users'),
                                                                           (0, 1, 3, 'Workflows', 'fa-project-diagram', '/info/workflows'),
                                                                           (0, 1, 4, 'Запуски', 'fa-play-circle', '/info/executions'),
                                                                           (0, 1, 5, 'Интеграции', 'fa-plug', '/info/connections'),
                                                                           (0, 1, 6, 'Подписки', 'fa-credit-card', '/info/subscriptions');

-- Связь меню с ролью
INSERT INTO goadmin_role_menu (role_id, menu_id)
SELECT 1, id FROM goadmin_menu;


-- Для users добавьте подсчет workflows
ALTER TABLE users ADD COLUMN IF NOT EXISTS subscription_plan VARCHAR(20) DEFAULT 'free';

-- Для workflows добавьте счетчики
ALTER TABLE workflows ADD COLUMN IF NOT EXISTS total_executions INT DEFAULT 0;
ALTER TABLE workflows ADD COLUMN IF NOT EXISTS success_count INT DEFAULT 0;
ALTER TABLE workflows ADD COLUMN IF NOT EXISTS error_count INT DEFAULT 0;
ALTER TABLE workflows ADD COLUMN IF NOT EXISTS trigger_type VARCHAR(20);

-- Для executions добавьте длительность
ALTER TABLE executions ADD COLUMN IF NOT EXISTS duration_seconds INT;
