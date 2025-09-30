-- roles
CREATE TABLE IF NOT EXISTS roles (
  id SERIAL PRIMARY KEY,
  name VARCHAR(50) NOT NULL UNIQUE
);

INSERT INTO roles (name) VALUES ('engineer'), ('manager'), ('observer') ON CONFLICT DO NOTHING;

-- users
CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email VARCHAR(255) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  full_name VARCHAR(255),
  role_id INT REFERENCES roles(id),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- projects / objects
CREATE TABLE IF NOT EXISTS projects (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(255) NOT NULL,
  description TEXT,
  created_by UUID REFERENCES users(id),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- defects
CREATE TABLE IF NOT EXISTS defects (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  priority INT DEFAULT 3, -- 1..5
  status VARCHAR(50) DEFAULT 'new',
  assignee UUID REFERENCES users(id),
  due_date DATE,
  created_by UUID REFERENCES users(id),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- comments
CREATE TABLE IF NOT EXISTS comments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  defect_id UUID REFERENCES defects(id) ON DELETE CASCADE,
  author UUID REFERENCES users(id),
  body TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- attachments (simplified)
CREATE TABLE IF NOT EXISTS attachments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  defect_id UUID REFERENCES defects(id) ON DELETE CASCADE,
  filename TEXT,
  content_type TEXT,
  url TEXT,
  uploaded_by UUID REFERENCES users(id),
  uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- history
CREATE TABLE IF NOT EXISTS defect_history (
  id SERIAL PRIMARY KEY,
  defect_id UUID REFERENCES defects(id) ON DELETE CASCADE,
  changed_by UUID REFERENCES users(id),
  change JSONB,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
