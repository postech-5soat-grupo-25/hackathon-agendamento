CREATE TABLE IF NOT EXISTS working_hours (
    id SERIAL PRIMARY KEY,
    doctor_id INT NOT NULL UNIQUE,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    days_of_week BOOLEAN[] NOT NULL
);

CREATE TABLE IF NOT EXISTS appointments (
    id SERIAL PRIMARY KEY,
    doctor_id INT NOT NULL REFERENCES working_hours(doctor_id) ON DELETE CASCADE,
    client_id INT NOT NULL,
    appointment_time TIMESTAMPTZ NOT NULL,
    description TEXT,
    UNIQUE(doctor_id, appointment_time),
    UNIQUE(client_id, appointment_time)
);

INSERT INTO working_hours (doctor_id, start_time, end_time, days_of_week)
VALUES (
    123,                              
    '2024-09-25 09:00:00+00',         
    '2024-09-25 17:00:00+00',         
    '{true, true, true, true, true, false, false}'
);

--psql -U postgres -d agendamento