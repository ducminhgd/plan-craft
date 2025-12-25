-- Initial database schema for Plan Craft
-- Creates all tables in dependency order

-- ============================================================================
-- 1. CLIENTS
-- ============================================================================
CREATE TABLE IF NOT EXISTS clients (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) UNIQUE,
    email VARCHAR(255),
    phone VARCHAR(50),
    address TEXT,
    contact_person VARCHAR(255),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    notes TEXT,

    CHECK (name != '')
);

CREATE INDEX IF NOT EXISTS idx_clients_is_active ON clients(is_active);
CREATE INDEX IF NOT EXISTS idx_clients_code ON clients(code);

-- ============================================================================
-- 2. PROJECTS
-- ============================================================================
CREATE TABLE IF NOT EXISTS projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Basic Information
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) UNIQUE,
    description TEXT,

    -- Project Classification
    type VARCHAR(50) NOT NULL DEFAULT 'product',
    status VARCHAR(50) NOT NULL DEFAULT 'not_started',

    -- Timeline
    start_date DATETIME,
    target_end_date DATETIME,
    actual_end_date DATETIME,

    -- Estimation (in hours)
    estimated_effort DECIMAL(10,2) DEFAULT 0,
    actual_effort DECIMAL(10,2) DEFAULT 0,

    -- Budget and Cost
    estimated_cost DECIMAL(15,2) DEFAULT 0,
    actual_cost DECIMAL(15,2) DEFAULT 0,
    currency VARCHAR(3) DEFAULT 'USD',

    -- Progress Tracking
    progress DECIMAL(5,2) DEFAULT 0,

    -- Work Time Configuration
    hours_per_day DECIMAL(5,2),
    days_per_week DECIMAL(5,2),
    days_per_month DECIMAL(5,2),

    -- Metadata
    assumptions JSON,
    constraints JSON,
    tags JSON,

    -- Ownership
    owner_id INTEGER,
    client_id INTEGER,

    FOREIGN KEY (client_id) REFERENCES clients(id),

    CHECK (progress >= 0 AND progress <= 100),
    CHECK (hours_per_day IS NULL OR hours_per_day > 0),
    CHECK (days_per_week IS NULL OR days_per_week > 0),
    CHECK (days_per_month IS NULL OR days_per_month > 0),
    CHECK (target_end_date IS NULL OR start_date IS NULL OR target_end_date >= start_date),
    CHECK (actual_end_date IS NULL OR start_date IS NULL OR actual_end_date >= start_date)
);

CREATE INDEX IF NOT EXISTS idx_projects_name ON projects(name);
CREATE INDEX IF NOT EXISTS idx_projects_status ON projects(status);
CREATE INDEX IF NOT EXISTS idx_projects_owner_id ON projects(owner_id);
CREATE INDEX IF NOT EXISTS idx_projects_client_id ON projects(client_id);

-- ============================================================================
-- 3. MILESTONES
-- ============================================================================
CREATE TABLE IF NOT EXISTS milestones (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    project_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'not_started',

    -- Timeline
    planned_start_date DATETIME,
    planned_end_date DATETIME,
    actual_start_date DATETIME,
    actual_end_date DATETIME,

    -- Estimation
    estimated_effort DECIMAL(10,2) DEFAULT 0,
    actual_effort DECIMAL(10,2) DEFAULT 0,

    -- Cost
    estimated_cost DECIMAL(15,2) DEFAULT 0,
    actual_cost DECIMAL(15,2) DEFAULT 0,

    -- Progress
    progress DECIMAL(5,2) DEFAULT 0,
    display_order INTEGER DEFAULT 0,

    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,

    CHECK (progress >= 0 AND progress <= 100),
    CHECK (planned_end_date IS NULL OR planned_start_date IS NULL OR planned_end_date >= planned_start_date)
);

CREATE INDEX IF NOT EXISTS idx_milestones_project ON milestones(project_id);
CREATE INDEX IF NOT EXISTS idx_milestones_status ON milestones(status);

-- ============================================================================
-- 4. TASKS
-- ============================================================================
CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Project and Hierarchy
    project_id INTEGER NOT NULL,
    milestone_id INTEGER,
    parent_task_id INTEGER,
    level INTEGER NOT NULL DEFAULT 1,

    -- Basic Information
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'not_started',
    priority VARCHAR(50) DEFAULT 'medium',

    -- Timeline
    planned_start_date DATETIME,
    planned_end_date DATETIME,
    actual_start_date DATETIME,
    actual_end_date DATETIME,

    -- Estimation (in hours)
    estimated_effort DECIMAL(10,2) DEFAULT 0,
    actual_effort DECIMAL(10,2) DEFAULT 0,

    -- Cost
    estimated_cost DECIMAL(15,2) DEFAULT 0,
    actual_cost DECIMAL(15,2) DEFAULT 0,

    -- Progress
    progress DECIMAL(5,2) DEFAULT 0,
    display_order INTEGER DEFAULT 0,

    -- Metadata
    tags JSON,
    notes TEXT,

    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (milestone_id) REFERENCES milestones(id) ON DELETE SET NULL,
    FOREIGN KEY (parent_task_id) REFERENCES tasks(id) ON DELETE CASCADE,

    CHECK (progress >= 0 AND progress <= 100),
    CHECK (level >= 1),
    CHECK (planned_end_date IS NULL OR planned_start_date IS NULL OR planned_end_date >= planned_start_date)
);

CREATE INDEX IF NOT EXISTS idx_tasks_project ON tasks(project_id);
CREATE INDEX IF NOT EXISTS idx_tasks_milestone ON tasks(milestone_id);
CREATE INDEX IF NOT EXISTS idx_tasks_parent ON tasks(parent_task_id);
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
CREATE INDEX IF NOT EXISTS idx_tasks_priority ON tasks(priority);

-- ============================================================================
-- 5. TASK DEPENDENCIES
-- ============================================================================
CREATE TABLE IF NOT EXISTS task_dependencies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    task_id INTEGER NOT NULL,
    depends_on_task_id INTEGER NOT NULL,
    dependency_type VARCHAR(50) NOT NULL DEFAULT 'finish_to_start',
    lag_days DECIMAL(5,2) DEFAULT 0,
    lead_days DECIMAL(5,2) DEFAULT 0,
    is_hard_dependency BOOLEAN NOT NULL DEFAULT TRUE,
    notes TEXT,

    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    FOREIGN KEY (depends_on_task_id) REFERENCES tasks(id) ON DELETE CASCADE,

    CHECK (task_id != depends_on_task_id),
    CHECK (NOT (lag_days > 0 AND lead_days > 0)),
    UNIQUE (task_id, depends_on_task_id)
);

CREATE INDEX IF NOT EXISTS idx_task_dependencies_task ON task_dependencies(task_id);
CREATE INDEX IF NOT EXISTS idx_task_dependencies_depends ON task_dependencies(depends_on_task_id);

-- ============================================================================
-- 6. RESOURCES
-- ============================================================================
CREATE TABLE IF NOT EXISTS resources (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Basic Information
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),

    -- Role Information
    role VARCHAR(100) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    -- Capacity
    default_hours_per_day DECIMAL(5,2) DEFAULT 8,
    default_days_per_week DECIMAL(5,2) DEFAULT 5,
    default_days_per_month DECIMAL(5,2) DEFAULT 20,

    -- Cost Information
    default_hourly_rate DECIMAL(10,2) DEFAULT 0,
    default_daily_rate DECIMAL(10,2) DEFAULT 0,
    default_monthly_rate DECIMAL(10,2) DEFAULT 0,
    currency VARCHAR(3) DEFAULT 'USD',

    -- Metadata
    skills JSON,
    notes TEXT,

    CHECK (default_hours_per_day >= 0 AND default_hours_per_day <= 24),
    CHECK (default_days_per_week >= 0 AND default_days_per_week <= 7),
    CHECK (default_days_per_month >= 0 AND default_days_per_month <= 31),
    CHECK (default_hourly_rate >= 0),
    CHECK (default_daily_rate >= 0),
    CHECK (default_monthly_rate >= 0)
);

CREATE INDEX IF NOT EXISTS idx_resources_email ON resources(email);
CREATE INDEX IF NOT EXISTS idx_resources_role ON resources(role);
CREATE INDEX IF NOT EXISTS idx_resources_is_active ON resources(is_active);

-- ============================================================================
-- 7. PROJECT ROLES
-- ============================================================================
CREATE TABLE IF NOT EXISTS project_roles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Project and Resource
    project_id INTEGER NOT NULL,
    resource_id INTEGER NOT NULL,

    -- Role Information
    role VARCHAR(100) NOT NULL,
    level VARCHAR(50) NOT NULL,

    -- Capacity in this Project
    hours_per_day DECIMAL(5,2),
    days_per_week DECIMAL(5,2),
    days_per_month DECIMAL(5,2),

    -- Estimated Allocation
    estimated_man_months DECIMAL(10,2) DEFAULT 0,

    -- Cost Rates in this Project
    hourly_rate DECIMAL(10,2),
    daily_rate DECIMAL(10,2),
    monthly_rate DECIMAL(10,2),

    -- Timeline
    start_date DATETIME,
    end_date DATETIME,

    -- Status
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    notes TEXT,

    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE,

    CHECK (hours_per_day IS NULL OR (hours_per_day >= 0 AND hours_per_day <= 24)),
    CHECK (days_per_week IS NULL OR (days_per_week >= 0 AND days_per_week <= 7)),
    CHECK (days_per_month IS NULL OR (days_per_month >= 0 AND days_per_month <= 31)),
    CHECK (hourly_rate IS NULL OR hourly_rate >= 0),
    CHECK (daily_rate IS NULL OR daily_rate >= 0),
    CHECK (monthly_rate IS NULL OR monthly_rate >= 0),
    CHECK (estimated_man_months >= 0),
    CHECK (end_date IS NULL OR start_date IS NULL OR end_date >= start_date)
);

CREATE INDEX IF NOT EXISTS idx_project_roles_project ON project_roles(project_id);
CREATE INDEX IF NOT EXISTS idx_project_roles_resource ON project_roles(resource_id);

-- ============================================================================
-- 8. RESOURCE ALLOCATIONS
-- ============================================================================
CREATE TABLE IF NOT EXISTS resource_allocations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Project Role Assignment
    project_role_id INTEGER NOT NULL,
    resource_id INTEGER NOT NULL,
    project_id INTEGER NOT NULL,

    -- Time Range
    start_date DATETIME NOT NULL,
    end_date DATETIME NOT NULL,

    -- Allocation Percentage (0-100)
    allocation_percent DECIMAL(5,2) NOT NULL,

    -- Capacity Override (optional)
    hours_per_day DECIMAL(5,2),
    days_per_week DECIMAL(5,2),
    days_per_month DECIMAL(5,2),

    -- Status
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    notes TEXT,

    FOREIGN KEY (project_role_id) REFERENCES project_roles(id) ON DELETE CASCADE,
    FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,

    CHECK (allocation_percent >= 0 AND allocation_percent <= 100),
    CHECK (hours_per_day IS NULL OR (hours_per_day >= 0 AND hours_per_day <= 24)),
    CHECK (days_per_week IS NULL OR (days_per_week >= 0 AND days_per_week <= 7)),
    CHECK (days_per_month IS NULL OR (days_per_month >= 0 AND days_per_month <= 31)),
    CHECK (end_date >= start_date)
);

CREATE INDEX IF NOT EXISTS idx_resource_allocation_role ON resource_allocations(project_role_id);
CREATE INDEX IF NOT EXISTS idx_resource_allocation_resource ON resource_allocations(resource_id);
CREATE INDEX IF NOT EXISTS idx_resource_allocation_project ON resource_allocations(project_id);
CREATE INDEX IF NOT EXISTS idx_resource_allocation_dates ON resource_allocations(start_date, end_date);
CREATE INDEX IF NOT EXISTS idx_resource_allocation_active ON resource_allocations(is_active);

-- ============================================================================
-- 9. TASK ASSIGNMENTS
-- ============================================================================
CREATE TABLE IF NOT EXISTS task_assignments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Task and Resource Assignment
    task_id INTEGER NOT NULL,
    project_role_id INTEGER NOT NULL,
    resource_id INTEGER NOT NULL,

    -- Effort Estimation (in man-days)
    estimated_man_days DECIMAL(10,2) NOT NULL,
    actual_man_days DECIMAL(10,2) DEFAULT 0,

    -- Allocation Percentage (0-100)
    allocation_percent DECIMAL(5,2) DEFAULT 100,

    -- Status
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    notes TEXT,

    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    FOREIGN KEY (project_role_id) REFERENCES project_roles(id) ON DELETE CASCADE,
    FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE,

    CHECK (estimated_man_days >= 0),
    CHECK (actual_man_days >= 0),
    CHECK (allocation_percent >= 0 AND allocation_percent <= 100)
);

CREATE INDEX IF NOT EXISTS idx_task_assignments_task ON task_assignments(task_id);
CREATE INDEX IF NOT EXISTS idx_task_assignments_role ON task_assignments(project_role_id);
CREATE INDEX IF NOT EXISTS idx_task_assignments_resource ON task_assignments(resource_id);

-- ============================================================================
-- 10. COSTS
-- ============================================================================
CREATE TABLE IF NOT EXISTS costs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Project, Milestone, and Task Association
    project_id INTEGER,
    milestone_id INTEGER,
    task_id INTEGER,

    -- Cost Classification
    type VARCHAR(50) NOT NULL,
    category VARCHAR(100),
    name VARCHAR(255) NOT NULL,

    -- Cost Details
    amount DECIMAL(15,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    quantity DECIMAL(10,2) DEFAULT 1,
    unit_cost DECIMAL(15,2) DEFAULT 0,

    -- For Labor Costs
    resource_id INTEGER,
    project_role_id INTEGER,
    rate_type VARCHAR(50),
    hours DECIMAL(10,2) DEFAULT 0,

    -- Status
    is_estimated BOOLEAN NOT NULL DEFAULT TRUE,
    date DATETIME,

    -- Additional Information
    notes TEXT,

    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (milestone_id) REFERENCES milestones(id) ON DELETE CASCADE,
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE,
    FOREIGN KEY (project_role_id) REFERENCES project_roles(id) ON DELETE CASCADE,

    CHECK (project_id IS NOT NULL OR milestone_id IS NOT NULL OR task_id IS NOT NULL),
    CHECK (amount >= 0),
    CHECK (quantity >= 0),
    CHECK (unit_cost >= 0),
    CHECK (hours >= 0)
);

CREATE INDEX IF NOT EXISTS idx_costs_project ON costs(project_id);
CREATE INDEX IF NOT EXISTS idx_costs_milestone ON costs(milestone_id);
CREATE INDEX IF NOT EXISTS idx_costs_task ON costs(task_id);
CREATE INDEX IF NOT EXISTS idx_costs_type ON costs(type);
CREATE INDEX IF NOT EXISTS idx_costs_category ON costs(category);
CREATE INDEX IF NOT EXISTS idx_costs_resource ON costs(resource_id);
CREATE INDEX IF NOT EXISTS idx_costs_project_role ON costs(project_role_id);
CREATE INDEX IF NOT EXISTS idx_costs_is_estimated ON costs(is_estimated);
CREATE INDEX IF NOT EXISTS idx_costs_date ON costs(date);

-- ============================================================================
-- TRIGGERS FOR AUTOMATIC TIMESTAMP UPDATES
-- ============================================================================

CREATE TRIGGER IF NOT EXISTS update_clients_updated_at
    AFTER UPDATE ON clients FOR EACH ROW
    BEGIN UPDATE clients SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id; END;

CREATE TRIGGER IF NOT EXISTS update_projects_updated_at
    AFTER UPDATE ON projects FOR EACH ROW
    BEGIN UPDATE projects SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id; END;

CREATE TRIGGER IF NOT EXISTS update_milestones_updated_at
    AFTER UPDATE ON milestones FOR EACH ROW
    BEGIN UPDATE milestones SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id; END;

CREATE TRIGGER IF NOT EXISTS update_tasks_updated_at
    AFTER UPDATE ON tasks FOR EACH ROW
    BEGIN UPDATE tasks SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id; END;

CREATE TRIGGER IF NOT EXISTS update_task_dependencies_updated_at
    AFTER UPDATE ON task_dependencies FOR EACH ROW
    BEGIN UPDATE task_dependencies SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id; END;

CREATE TRIGGER IF NOT EXISTS update_resources_updated_at
    AFTER UPDATE ON resources FOR EACH ROW
    BEGIN UPDATE resources SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id; END;

CREATE TRIGGER IF NOT EXISTS update_project_roles_updated_at
    AFTER UPDATE ON project_roles FOR EACH ROW
    BEGIN UPDATE project_roles SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id; END;

CREATE TRIGGER IF NOT EXISTS update_resource_allocations_updated_at
    AFTER UPDATE ON resource_allocations FOR EACH ROW
    BEGIN UPDATE resource_allocations SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id; END;

CREATE TRIGGER IF NOT EXISTS update_task_assignments_updated_at
    AFTER UPDATE ON task_assignments FOR EACH ROW
    BEGIN UPDATE task_assignments SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id; END;

CREATE TRIGGER IF NOT EXISTS update_costs_updated_at
    AFTER UPDATE ON costs FOR EACH ROW
    BEGIN UPDATE costs SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id; END;
