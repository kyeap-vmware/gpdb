-- Tests for restore point (RP) based transaction isolation for hot standby

!\retcode gpconfig -c gp_hot_standby_snapshot_mode -v restorepoint;
!\retcode gpconfig -c default_transaction_isolation -v "'repeatable read'";
-- we will smoke testing the related logging under Debug_print_snapshot_dtm too
!\retcode gpconfig -c debug_print_snapshot_dtm -v true --skipvalidation;
!\retcode gpstop -u;

----------------------------------------------------------------
-- Basic transaction isolation test
----------------------------------------------------------------

1: create table hs_rp_basic(a int, b text);

-- in-progress transaction won't be visible
2: begin;
2: insert into hs_rp_basic select i,'in_progress' from generate_series(1,5) i;
-- transactions completed before the RP: all would be visible on standby, including 1PC and 2PC
1: insert into hs_rp_basic select i,'complete_before_rp1_2pc' from generate_series(1,5) i;
1: insert into hs_rp_basic values(1, 'complete_before_rp1_1pc');

-- take the RP, ignore return value which includes LSN
1: select sum(1) from gp_create_restore_point('rp1');
-- set the RP to use at system-level, this is what would typically be done in real life (e.g. in GPDR)
!\retcode gpconfig -c gp_hot_standby_snapshot_restore_point_name -v rp1;
!\retcode gpstop -u;

-- transactions after the RP: won't be visible on standby
1: insert into hs_rp_basic select i,'complete_after_rp1_2pc' from generate_series(1,5) i;
1: insert into hs_rp_basic values(1, 'complete_after_rp1_1pc');

-- standby sees expected result
-1S: select * from hs_rp_basic;

-- more completed transactions, including completing the in-progress one
1: insert into hs_rp_basic select i,'complete_after_rp1' from generate_series(1,5) i;
2: update hs_rp_basic set b = 'in_progress_at_rp1_complete_after_rp1' where b = 'in_progress';
2: end;
-- still won't be seen on the standby
-1S: select * from hs_rp_basic;

-- a new RP is created
1: select sum(1) from gp_create_restore_point('rp2');
1: insert into hs_rp_basic select i,'complete_after_rp2_2pc' from generate_series(1,5) i;
1: insert into hs_rp_basic select i,'complete_after_rp2_1pc' from generate_series(1,5) i;

-- the standby uses it, and sees all that completed before 'rp2';
-- also testing that we can set the GUC over the session.
-1S: set gp_hot_standby_snapshot_restore_point_name = 'rp2';
-1S: select * from hs_rp_basic;

-- using an earlier RP, and sees result according to that
-1S: set gp_hot_standby_snapshot_restore_point_name = 'rp1';
-1S: select * from hs_rp_basic;

----------------------------------------------------------------
-- Coverage for incorrect transaction isolation level
----------------------------------------------------------------

!\retcode gpconfig -c default_transaction_isolation -v "'read committed'";
!\retcode gpstop -u;

-- error out
-1S: select * from hs_rp_basic;

!\retcode gpconfig -c default_transaction_isolation -v "'repeatable read'";
!\retcode gpstop -u;

----------------------------------------------------------------
-- Coverage for incorrect usage of the gp_hot_standby_snapshot_restore_point_name GUC
----------------------------------------------------------------

-- no RP name is given, should error out.
!\retcode gpconfig -r gp_hot_standby_snapshot_restore_point_name;
!\retcode gpstop -u;
-1S: reset gp_hot_standby_snapshot_restore_point_name;
-1S: select * from hs_rp_basic;

-- standby uses an invalid RP name, should complain
-1S: set gp_hot_standby_snapshot_restore_point_name = 'rp-non-exist';
-1S: select * from hs_rp_basic;

-- primary creates a same-name RP, standby won't create snapshot for it,
-- and still use the old snapshot
1: select sum(1) from gp_create_restore_point('rp1');
-1S: set gp_hot_standby_snapshot_restore_point_name = 'rp1';
-1S: select * from hs_rp_basic;

----------------------------------------------------------------
-- Coverage for gp_max_restore_point_snapshots
----------------------------------------------------------------

1q:
2q:
-1Sq:

-- set the GUC to 1
!\retcode gpconfig -c gp_max_restore_point_snapshots -v 1;
!\retcode gpstop -ar;

-- primary creates more RP, but standby won't create snapshot for it
1: select sum(1) from gp_create_restore_point('rp-good');
1: select sum(1) from gp_create_restore_point('rp-bad');

-1S: set gp_hot_standby_snapshot_restore_point_name = 'rp-good';
-1S: select * from hs_rp_basic;

-1S: set gp_hot_standby_snapshot_restore_point_name = 'rp-bad';
-1S: select * from hs_rp_basic;

1q:
-1Sq:

!\retcode gpconfig -r gp_max_restore_point_snapshots;
!\retcode gpstop -ar;

----------------------------------------------------------------
-- Tests that simulates out-of-sync WAL replays on standby coordinator and segments
----------------------------------------------------------------
-- syncrep needs to be turned off for these tests
1: set synchronous_commit = off;

1: select sum(1) from gp_create_restore_point('rp3');
!\retcode gpconfig -c gp_hot_standby_snapshot_restore_point_name -v rp3;
!\retcode gpstop -u;

--
-- Case 1: standby coordinator WAL replay is delayed
--
-1S: select pg_wal_replay_pause();
1: select sum(1) from gp_create_restore_point('rp-qe-only');

-- This RP is only replayed on (some) QEs so far, so the query will fail.
-1S: set gp_hot_standby_snapshot_restore_point_name = 'rp-qe-only';
-1S: select * from hs_rp_basic;
-- And if we use 'rp3', it should work.
-1S: set gp_hot_standby_snapshot_restore_point_name = 'rp3';
-1S: select * from hs_rp_basic;

-- resume replay on QD
-1S: select pg_wal_replay_resume();

--
-- Case 2: standby segment WAL replay is delayed
--
-1S: select pg_wal_replay_pause() from gp_dist_random('gp_id');
1: select sum(1) from gp_create_restore_point('rp-qd-only');

-- This RP is only replayed on QD so far, so the query will fail.
-1S: set gp_hot_standby_snapshot_restore_point_name = 'rp-qd-only';
-1S: select * from hs_rp_basic;
-- And if we use 'rp3', it should work.
-1S: set gp_hot_standby_snapshot_restore_point_name = 'rp3';
-1S: select * from hs_rp_basic;

--resume
-1S: select pg_wal_replay_resume() from gp_dist_random('gp_id');

-- reset for the rest of tests
1: reset synchronous_commit;

----------------------------------------------------------------
-- Crash recovery test
----------------------------------------------------------------

1: create table hs_rp_crash(a int, b text);

--
-- Case 1: standby coordinator/segment restart, and replay from a checkpoint 
-- that's behind the RP which they are going to use.
--

-- completed tx before RP, will be seen
1: insert into hs_rp_crash select i,'complete_before_rp' from generate_series(1,5)i;
-- in-progress tx on RP, won't be seen
2: begin;
2: insert into hs_rp_crash select i,'in_progress_at_rp' from generate_series(1,5)i;

-- create RP, and set it at system-level (so we won't clean the snapshot up during restart)
1: select sum(1) from gp_create_restore_point('rptest_crash1');
!\retcode gpconfig -c gp_hot_standby_snapshot_restore_point_name -v rptest_crash1;
!\retcode gpstop -u;

-- completed tx after RP, won't be seen
1: insert into hs_rp_crash select i,'complete_after_rp' from generate_series(1,5)i;

2: end;

-- make sure that after restart, standby replays from a point *after* the RP
1: checkpoint;

-- standby coordinator restarts
-1S: select gp_inject_fault('exec_simple_query_start', 'panic', dbid) from gp_segment_configuration where content=-1 and role='m';
-1S: select 1;
-1Sq:

-- sees expected result corresponding to the RP
-1S: select * from hs_rp_crash;

-- standby segment restarts, and still can use previous RP too

-- seg0 restarts
-1S: select gp_inject_fault('exec_mpp_query_start', 'panic', dbid) from gp_segment_configuration where content=0 and role='m';
-1S: select * from hs_rp_crash;

-- sees expected result corresponding to the RP
-1S: select * from hs_rp_crash;

--
-- Case 2: standby coordinator/segment restart, and replay from a checkpoint 
-- that's behind the RP which they are going to use.
-- The effect should be the same as Case 1.
--

1: truncate hs_rp_crash;

-- completed tx before RP, will be seen
1: insert into hs_rp_crash select i,'complete_before_rp' from generate_series(1,5)i;
-- in-progress tx on RP, won't be seen
2: begin;
2: insert into hs_rp_crash select i,'in_progress_at_rp' from generate_series(1,5)i;

-- make sure that after restart, standby replays from a point *before* the RP
1: checkpoint;

-- create RP, and set it at system-level (so we won't clean the snapshot up during restart)
1: select sum(1) from gp_create_restore_point('rptest_crash2');
!\retcode gpconfig -c gp_hot_standby_snapshot_restore_point_name -v rptest_crash2;
!\retcode gpstop -u;

-- completed tx after RP, won't be seen
1: insert into hs_rp_crash select i,'complete_after_rp' from generate_series(1,5)i;

2: end;
2q:

-- standby coordinator restarts
-1S: select gp_inject_fault('exec_simple_query_start', 'panic', dbid) from gp_segment_configuration where content=-1 and role='m';
-1S: select 1;
-1Sq:

-- sees expected result corresponding to the RP
-1S: select * from hs_rp_crash;

-- standby segment restarts, and still can use previous RP too

-- seg0 restarts
-1S: select gp_inject_fault('exec_mpp_query_start', 'panic', dbid) from gp_segment_configuration where content=0 and role='m';
-1S: select * from hs_rp_crash;

-- sees expected result corresponding to the RP
-1S: select * from hs_rp_crash;

----------------------------------------------------------------
-- Snapshot conflict test
----------------------------------------------------------------

1: create table hs_rp_conflict(a int);

-- The primary inserts some rows, creates an RP, then deletes & vacuums all the rows.
-- The standby query, using that RP, will conflict and be cancelled.
1: insert into hs_rp_conflict select * from generate_series(1,10);
1: select sum(1) from gp_create_restore_point('rp_conflict');
-1S: set gp_hot_standby_snapshot_restore_point_name = 'rp_conflict';
1: delete from hs_rp_conflict;
1: vacuum hs_rp_conflict;
1q:

-- The RP is invalidated and the snapshot deleted, the query will fail
-1S: select count(*) from hs_rp_conflict;
-1Sq:

-- Because the VACUUM invalidates the latest RP, it effectively also invalidated all 
-- RPs prior to that. So segments shouldn't have any snapshots left on disk.
-- In order to run the pg_ls_dir, first set the snapshot mode to inconsistent (since all RPs/snapshots are gone).
!\retcode gpconfig -c gp_hot_standby_snapshot_mode -v inconsistent;
!\retcode gpstop -u;
-1S: select gp_segment_id, pg_ls_dir('pg_snapshots') from gp_dist_random('gp_id');
-- coordinator still has it because snapshot conflict only happens on segments.
-1S: select pg_ls_dir('pg_snapshots');
-1Sq:
-- restart it should clean up all snapshots on coordinator.
!\retcode gpstop -ar;
-1S: select pg_ls_dir('pg_snapshots');

!\retcode gpconfig -r default_transaction_isolation;
!\retcode gpconfig -r debug_print_snapshot_dtm --skipvalidation;
!\retcode gpstop -u;
