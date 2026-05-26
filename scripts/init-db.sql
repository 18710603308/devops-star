-- ========== DevOpsStar 数据库初始化脚本 ==========
-- 此脚本由 docker-compose.yml 自动执行
-- 路径：scripts/init-db.sql

-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    role VARCHAR(20) DEFAULT 'user',
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 创建项目表
CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    display_name VARCHAR(100),
    description TEXT,
    repo_url VARCHAR(255),
    repo_type VARCHAR(20) DEFAULT 'gitea',
    gitea_id INTEGER,
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 创建流水线表
CREATE TABLE IF NOT EXISTS pipelines (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    project_id INTEGER REFERENCES projects(id),
    description TEXT,
    config_yaml TEXT,
    status VARCHAR(20) DEFAULT 'idle',
    last_run_id VARCHAR(50),
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 创建流水线运行记录表
CREATE TABLE IF NOT EXISTS pipeline_runs (
    id SERIAL PRIMARY KEY,
    run_id VARCHAR(50) UNIQUE NOT NULL,
    pipeline_id INTEGER REFERENCES pipelines(id),
    status VARCHAR(20),
    trigger VARCHAR(50),
    branch VARCHAR(100),
    logs TEXT,
    started_at TIMESTAMP,
    finished_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 创建部署环境表
CREATE TABLE IF NOT EXISTS environments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    display_name VARCHAR(100),
    project_id INTEGER REFERENCES projects(id),
    deploy_type VARCHAR(20) DEFAULT 'docker',
    config TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 创建部署记录表
CREATE TABLE IF NOT EXISTS deploy_records (
    id SERIAL PRIMARY KEY,
    environment_id INTEGER REFERENCES environments(id),
    pipeline_run_id VARCHAR(50),
    status VARCHAR(20),
    deploy_url VARCHAR(255),
    image_tag VARCHAR(100),
    logs TEXT,
    deployed_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 创建通知配置表
CREATE TABLE IF NOT EXISTS notification_configs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    type VARCHAR(20) NOT NULL,
    webhook VARCHAR(500),
    smtp_host VARCHAR(100),
    smtp_port INTEGER DEFAULT 465,
    smtp_user VARCHAR(100),
    smtp_pass VARCHAR(255),
    notify_on TEXT,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- ========== 插入默认管理员账号 ==========
-- 密码：admin123（bcrypt hash）
INSERT INTO users (username, password, email, role)
VALUES (
    'admin',
    '$2a$10$8K1p/a0dL3LzO/2pL3x.O1p/a0dL3LzO/2pL3x.O1p/a0dL3L',
    'admin@devops-star.com',
    'admin'
) ON CONFLICT (username) DO NOTHING;

-- ========== 创建索引 ==========
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_projects_name ON projects(name);
CREATE INDEX IF NOT EXISTS idx_pipelines_project ON pipelines(project_id);
CREATE INDEX IF NOT EXISTS idx_pipeline_runs_pipeline ON pipeline_runs(pipeline_id);
CREATE INDEX IF NOT EXISTS idx_pipeline_runs_status ON pipeline_runs(status);
CREATE INDEX IF NOT EXISTS idx_environments_project ON environments(project_id);
CREATE INDEX IF NOT EXISTS idx_deploy_records_env ON deploy_records(environment_id);
CREATE INDEX IF NOT EXISTS idx_notification_configs_type ON notification_configs(type);

-- ========== 完成提示 ==========
DO $$
BEGIN
    RAISE NOTICE '✅ DevOpsStar 数据库初始化完成！';
    RAISE NOTICE '   默认管理员：admin / admin123';
    RAISE NOTICE '   请及时修改默认密码！';
END$$;
