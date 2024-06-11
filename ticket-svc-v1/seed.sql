CREATE TABLE int_ranges (
    id SERIAL PRIMARY KEY,
    start_range INT NOT NULL,
    end_range INT NOT NULL,
    current_number INT NOT NULL,
    UNIQUE(start_range, end_range)
);

INSERT INTO int_ranges (start_range, end_range, current_number) VALUES 
(1, 1000, 1),
(1001, 2000, 1001),
(2001, 3000, 2001),
(3001, 4000, 3001),
(4001, 5000, 4001),
(5001, 6000, 5001),
(6001, 7000, 6001),
(7001, 8000, 7001),
(8001, 9000, 8001),
(9001, 10000, 9001);
