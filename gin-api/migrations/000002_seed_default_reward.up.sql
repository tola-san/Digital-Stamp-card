INSERT INTO rewards (name, description, required_stamps, active)
SELECT 'Free Coffee', 'One free coffee after collecting 10 stamps.', 10, TRUE
WHERE NOT EXISTS (SELECT 1 FROM rewards);
