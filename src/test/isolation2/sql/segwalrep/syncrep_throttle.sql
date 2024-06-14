-- set wait_for_replication_threshold to 1kB for quicker test
!\retcode gpconfig -c wait_for_replication_threshold -v 1;
!\retcode gpstop -u;

----------
-- INSERT
----------

CREATE TABLE insert_throttle(i int);

-- Suspend walsender
SELECT gp_inject_fault_infinite('wal_sender_loop', 'suspend', dbid) FROM
    gp_segment_configuration WHERE role = 'p' and content = 1;

-- This should wait for syncrep since its WAL size greater than wait_for_replication_threshold
1&:INSERT INTO insert_throttle SELECT 1 FROM generate_series(1, 1000000);

SELECT is_query_waiting_for_syncrep(50,
    'INSERT INTO insert_throttle SELECT 1 FROM generate_series(1, 1000000);');

-- Smoke test: ensure CHECKPOINTs are not blocked while we are waiting on syncrep.
CHECKPOINT;

-- reset walsender
SELECT gp_inject_fault_infinite('wal_sender_loop', 'reset', dbid) FROM
    gp_segment_configuration WHERE role = 'p' and content = 1;

1<:

----------
-- DELETE
----------

CREATE TABLE del_throttle(i int);
INSERT INTO del_throttle SELECT 1 FROM generate_series(1, 1000000);

-- Suspend walsender
SELECT gp_inject_fault_infinite('wal_sender_loop', 'suspend', dbid) FROM
    gp_segment_configuration WHERE role = 'p' and content = 1;

-- This should wait for syncrep since its WAL size greater than wait_for_replication_threshold
1&:DELETE FROM del_throttle;

SELECT is_query_waiting_for_syncrep(50, 'DELETE FROM del_throttle;');

-- Smoke test: ensure CHECKPOINTs are not blocked while we are waiting on syncrep.
CHECKPOINT;

-- reset walsender
SELECT gp_inject_fault_infinite('wal_sender_loop', 'reset', dbid) FROM
    gp_segment_configuration WHERE role = 'p' and content = 1;

1<:

----------
-- UPDATE
----------

CREATE TABLE update_throttle(i int, j int);
INSERT INTO update_throttle SELECT 1, j FROM generate_series(1, 1000000) j;

-- Suspend walsender
SELECT gp_inject_fault_infinite('wal_sender_loop', 'suspend', dbid) FROM
    gp_segment_configuration WHERE role = 'p' and content = 1;

SET optimizer TO off;
EXPLAIN (COSTS OFF) UPDATE update_throttle SET j = j + 1;

-- This should wait for syncrep since its WAL size greater than wait_for_replication_threshold
1&:UPDATE update_throttle SET j = j + 1;

SELECT is_query_waiting_for_syncrep(50, 'UPDATE update_throttle SET j = j + 1;');

-- Smoke test: ensure CHECKPOINTs are not blocked while we are waiting on syncrep.
CHECKPOINT;

-- reset walsender
SELECT gp_inject_fault_infinite('wal_sender_loop', 'reset', dbid) FROM
    gp_segment_configuration WHERE role = 'p' and content = 1;

1<:

RESET optimizer;

-------------------
-- UPDATE (Split)
-------------------

CREATE TABLE split_update_throttle(i int);
INSERT INTO split_update_throttle SELECT 1 FROM generate_series(1, 1000000);

-- Suspend walsender
SELECT gp_inject_fault_infinite('wal_sender_loop', 'suspend', dbid) FROM
    gp_segment_configuration WHERE role = 'p' and content = 1;

SET optimizer TO off;
EXPLAIN (COSTS OFF) UPDATE split_update_throttle SET i = 20;

-- This should wait for syncrep since its WAL size greater than wait_for_replication_threshold
1&:UPDATE split_update_throttle SET i = 20;

SELECT is_query_waiting_for_syncrep(50, 'UPDATE split_update_throttle SET i = 20;');

-- Smoke test: ensure CHECKPOINTs are not blocked while we are waiting on syncrep.
CHECKPOINT;

-- reset walsender
SELECT gp_inject_fault_infinite('wal_sender_loop', 'reset', dbid) FROM
    gp_segment_configuration WHERE role = 'p' and content = 1;

1<:

RESET optimizer;

!\retcode gpconfig -r wait_for_replication_threshold;
!\retcode gpstop -u;
