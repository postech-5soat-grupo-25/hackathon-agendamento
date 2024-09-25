CREATE TABLE working_hours (
    id SERIAL PRIMARY KEY,
    doctor_id INT NOT NULL UNIQUE,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    days_of_week BOOLEAN[] NOT NULL
);

INSERT INTO working_hours (doctor_id, start_time, end_time, days_of_week)
VALUES (
    123,                              
    '2024-09-25 09:00:00+00',         
    '2024-09-25 17:00:00+00',         
    '{true, true, true, true, true, false, false}'
);

--psql -U postgres -d agendamento