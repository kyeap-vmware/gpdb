-- The DELETE was not replayed so should show all 10 rows.
SELECT * FROM gpdb_pitr_table;

-- The INSERT happened after the restore point was created
-- so this table should be empty.
SELECT * FROM gpdb_restore_points;
