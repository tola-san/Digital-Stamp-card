DELETE FROM rewards
WHERE name = 'Free Coffee'
  AND description = 'One free coffee after collecting 10 stamps.'
  AND required_stamps = 10;
