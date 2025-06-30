-- DDL for schedules table
CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; -- Required for UUID generation

CREATE TABLE IF NOT EXISTS schedules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    client_name VARCHAR(255) NOT NULL,
    shift_time TIMESTAMPTZ NOT NULL,
    location VARCHAR(255) NOT NULL, -- General location string, e.g., "123 Main St, Anytown"
    status VARCHAR(50) NOT NULL DEFAULT 'upcoming', -- e.g., 'upcoming', 'in-progress', 'completed', 'missed'
    start_time TIMESTAMPTZ NULL,
    start_latitude NUMERIC(10, 8) NULL,
    start_longitude NUMERIC(11, 8) NULL,
    end_time TIMESTAMPTZ NULL,
    end_latitude NUMERIC(10, 8) NULL,
    end_longitude NUMERIC(11, 8) NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- DDL for tasks table
CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    schedule_id UUID NOT NULL,
    description TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- e.g., 'pending', 'completed', 'not_completed'
    reason TEXT NULL, -- Optional reason if not completed
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_schedule
        FOREIGN KEY(schedule_id)
            REFERENCES schedules(id)
            ON DELETE CASCADE
);

-- Index for faster lookup by schedule_id in tasks table
CREATE INDEX IF NOT EXISTS idx_tasks_schedule_id ON tasks (schedule_id);

-- Test Data (Optional: You can run these inserts after creating tables)

-- Insert sample schedules
INSERT INTO schedules (id, client_name, shift_time, location, status) VALUES
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Alice Johnson', NOW() + INTERVAL '2 hour', '123 Oak Ave, City, ST', 'upcoming'),
('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'Bob Smith', NOW() - INTERVAL '1 hour', '456 Pine St, Town, ST', 'in-progress'),
('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'Charlie Brown', NOW() - INTERVAL '2 day', '789 Elm St, Village, ST', 'completed'),
('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', 'Diana Miller', NOW() + INTERVAL '1 day', '101 Birch Ln, Hamlet, ST', 'upcoming');

-- Update the 'in-progress' schedule with start details
UPDATE schedules
SET start_time = NOW() - INTERVAL '1 hour',
    start_latitude = 34.0522,
    start_longitude = -118.2437
WHERE id = 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12';

-- Update the 'completed' schedule with start and end details
UPDATE schedules
SET start_time = NOW() - INTERVAL '2 day' - INTERVAL '4 hour',
    start_latitude = 33.99,
    start_longitude = -118.45,
    end_time = NOW() - INTERVAL '2 day' - INTERVAL '2 hour',
    end_latitude = 33.99,
    end_longitude = -118.45,
    status = 'completed'
WHERE id = 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13';


-- Insert sample tasks for schedule 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12' (Bob Smith - in-progress)
INSERT INTO tasks (id, schedule_id, description, status) VALUES
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a15', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'Check vital signs', 'pending'),
('f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a16', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'Administer medication (2:00 PM)', 'pending'),
('60eebc99-9c0b-4ef8-bb6d-6bb9bd380a17', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'Assist with light meal preparation', 'completed');

-- Insert sample tasks for schedule 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13' (Charlie Brown - completed)
INSERT INTO tasks (id, schedule_id, description, status) VALUES
('80eebc99-9c0b-4ef8-bb6d-6bb9bd380a18', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'Help with bathing', 'completed'),
('90eebc99-9c0b-4ef8-bb6d-6bb9bd380a19', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'Perform physical therapy exercises', 'completed'),
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a20', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'Record daily observations', 'completed');

-- Insert more sample schedules
INSERT INTO schedules (id, client_name, shift_time, location, status) VALUES
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a21', 'Eve Geller', NOW() + INTERVAL '5 hour', '707 Cedar Rd, Suburb, ST', 'upcoming'),
('f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'Frank White', NOW() - INTERVAL '3 day', '888 Maple Dr, Rural, ST', 'upcoming'),
('11eebc99-9c0b-4ef8-bb6d-6bb9bd380a23', 'Grace Lee', NOW() + INTERVAL '10 hour', '999 Willow Ct, Uptown, ST', 'upcoming'),
('22eebc99-9c0b-4ef8-bb6d-6bb9bd380a24', 'Henry Adams', NOW() - INTERVAL '1 day', '111 Elm St, Downtown, ST', 'completed'),
('33eebc99-9c0b-4ef8-bb6d-6bb9bd380a25', 'Ivy King', NOW() + INTERVAL '3 day', '222 Oak St, Westside, ST', 'upcoming');

-- Update 'missed' schedule (no start/end times typically)
UPDATE schedules
SET updated_at = NOW()
WHERE id = 'f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22';

-- Update 'completed' schedule (Henry Adams)
UPDATE schedules
SET start_time = NOW() - INTERVAL '1 day' - INTERVAL '3 hour',
    start_latitude = 40.7128,
    start_longitude = -74.0060,
    end_time = NOW() - INTERVAL '1 day' - INTERVAL '1 hour',
    end_latitude = 40.7128,
    end_longitude = -74.0060,
    status = 'completed'
WHERE id = '22eebc99-9c0b-4ef8-bb6d-6bb9bd380a24';

-- Insert tasks for new schedules

-- Tasks for 'e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a21' (Eve Geller - upcoming)
INSERT INTO tasks (id, schedule_id, description, status) VALUES
('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a26', 'e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a21', 'Prepare breakfast', 'pending'),
('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a27', 'e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a21', 'Light housekeeping', 'pending');

-- Tasks for 'f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22' (Frank White - missed)
INSERT INTO tasks (id, schedule_id, description, status, reason) VALUES
('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a28', 'f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'Check on pets', 'not_completed', 'Client unreachable'),
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a29', 'f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'Deliver groceries', 'not_completed', 'Access denied');

-- Tasks for '11eebc99-9c0b-4ef8-bb6d-6bb9bd380a23' (Grace Lee - upcoming)
INSERT INTO tasks (id, schedule_id, description, status) VALUES
('f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a30', '11eebc99-9c0b-4ef8-bb6d-6bb9bd380a23', 'Escort to doctor appointment', 'pending'),
('60eebc99-9c0b-4ef8-bb6d-6bb9bd380a31', '11eebc99-9c0b-4ef8-bb6d-6bb9bd380a23', 'Medication reminder', 'pending');

-- Tasks for '22eebc99-9c0b-4ef8-bb6d-6bb9bd380a24' (Henry Adams - completed)
INSERT INTO tasks (id, schedule_id, description, status) VALUES
('70eebc99-9c0b-4ef8-bb6d-6bb9bd380a32', '22eebc99-9c0b-4ef8-bb6d-6bb9bd380a24', 'Read aloud for 30 minutes', 'completed'),
('80eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', '22eebc99-9c0b-4ef8-bb6d-6bb9bd380a24', 'Organize pantry', 'completed');

-- Tasks for '33eebc99-9c0b-4ef8-bb6d-6bb9bd380a25' (Ivy King - upcoming)
INSERT INTO tasks (id, schedule_id, description, status) VALUES
('90eebc99-9c0b-4ef8-bb6d-6bb9bd380a34', '33eebc99-9c0b-4ef8-bb6d-6bb9bd380a25', 'Assist with grocery shopping list', 'pending'),
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a35', '33eebc99-9c0b-4ef8-bb6d-6bb9bd380a25', 'Water plants', 'pending');

INSERT INTO schedules (id, client_name, shift_time, location, status) VALUES
('30eebc99-9c0b-4ef8-bb6d-6bb9bd380a32', 'Olivia Green', NOW() + INTERVAL '3 day', '222 Cedar Dr, Austin, TX', 'upcoming'),
('40eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'Peter Black', NOW() - INTERVAL '6 hour', '333 Maple Ave, Denver, CO', 'in-progress'),
('50eebc99-9c0b-4ef8-bb6d-6bb9bd380a34', 'Quinn Taylor', NOW() + INTERVAL '1 week', '444 Spruce St, Chicago, IL', 'upcoming'),
('60eebc99-9c0b-4ef8-bb6d-6bb9bd380a35', 'Rachel King', NOW() - INTERVAL '5 day', '555 Walnut Blvd, Boston, MA', 'completed'),
('70eebc99-9c0b-4ef8-bb6d-6bb9bd380a36', 'Sam Clark', NOW() - INTERVAL '1 day', '666 Pine St, Dallas, TX', 'completed'); 

UPDATE schedules
SET start_time = NOW() - INTERVAL '6 hour',
    start_latitude = 39.7392,
    start_longitude = -104.9903
WHERE id = '40eebc99-9c0b-4ef8-bb6d-6bb9bd380a33';

UPDATE schedules
SET start_time = NOW() - INTERVAL '5 day' - INTERVAL '8 hour',
    start_latitude = 42.3601,
    start_longitude = -71.0589,
    end_time = NOW() - INTERVAL '5 day' - INTERVAL '6 hour',
    end_latitude = 42.3601,
    end_longitude = -71.0589
WHERE id = '60eebc99-9c0b-4ef8-bb6d-6bb9bd380a35';

UPDATE schedules
SET start_time = NOW() - INTERVAL '1 day' - INTERVAL '3 hour',
    start_latitude = 32.7767,
    start_longitude = -96.7970,
    end_time = NOW() - INTERVAL '1 day' - INTERVAL '1 hour',
    end_latitude = 32.7767,
    end_longitude = -96.7970,
    status = 'completed'
WHERE id = '70eebc99-9c0b-4ef8-bb6d-6bb9bd380a36';

INSERT INTO tasks (id, schedule_id, description, status) VALUES
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a37', '30eebc99-9c0b-4ef8-bb6d-6bb9bd380a32', 'Help with grocery shopping', 'pending'),
('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a38', '30eebc99-9c0b-4ef8-bb6d-6bb9bd380a32', 'Light cleaning of living room', 'pending');

-- Insert tasks for Peter Black (in-progress) - no change
INSERT INTO tasks (id, schedule_id, description, status) VALUES
('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a39', '40eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'Monitor blood pressure', 'pending'),
('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a40', '40eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'Prepare lunch', 'completed');

-- Insert tasks for Rachel King (completed) - no change
INSERT INTO tasks (id, schedule_id, description, status) VALUES
('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a41', '60eebc99-9c0b-4ef8-bb6d-6bb9bd380a35', 'Assist with walking exercise', 'completed'),
('f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a42', '60eebc99-9c0b-4ef8-bb6d-6bb9bd380a35', 'Organize medication for the week', 'completed');

INSERT INTO tasks (id, schedule_id, description, status, reason) VALUES
('01eebc99-9c0b-4ef8-bb6d-6bb9bd380a43', '70eebc99-9c0b-4ef8-bb6d-6bb9bd380a36', 'Pick up prescription', 'completed', NULL), 
('02eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', '70eebc99-9c0b-4ef8-bb6d-6bb9bd380a36', 'Companionship visit', 'completed', NULL); 