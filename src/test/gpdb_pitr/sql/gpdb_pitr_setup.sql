-- Set a checkpoint and then disable automatic checkpoints
CHECKPOINT;
CREATE EXTENSION IF NOT EXISTS gp_inject_fault;
SELECT gp_inject_fault_infinite('checkpoint', 'skip', dbid) FROM gp_segment_configuration WHERE role = 'p';

-- Create some tables and load some data
CREATE TABLE gpdb_pitr_table(num int);
CREATE TABLE gpdb_restore_points(rpname text, rplsn pg_lsn);
INSERT INTO gpdb_pitr_table SELECT generate_series(1, 10);

-- Create restore point and do some DML that will not be recovered
INSERT INTO gpdb_restore_points VALUES ('foo', gp_create_restore_point('foo'));
DELETE FROM gpdb_pitr_table;
