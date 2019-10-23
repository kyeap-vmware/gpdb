#!/usr/bin/env bash

## ==================================================================
## Required: sourced gpdemo cluster with default mirror setup
##
## TODO: change to use pg_basebackup instead of the gpdemo mirrors
## ==================================================================

## Run tests via pg_regress with given test name(s)
run_tests()
{
    ../regress/pg_regress --dbname=gpdb_pitr_database --use-existing --init-file=../regress/init_file $1
    if [ $? != 0 ]; then
        exit 1
    fi
}

DATADIR="${MASTER_DATA_DIRECTORY%*/*/*}"
STANDBY=${DATADIR}/standby
MIRROR1=${DATADIR}/dbfast_mirror1/demoDataDir0
MIRROR2=${DATADIR}/dbfast_mirror2/demoDataDir1
MIRROR3=${DATADIR}/dbfast_mirror3/demoDataDir2

# Create test database.
createdb gpdb_pitr_database

# Run setup test.
run_tests gpdb_pitr_setup

# force kill cluster
# gpstop -ai // TODO: why does MinRecoveryPoint mess up on standby master?
pkill -9 post
rm -rf /tmp/.s.PGSQL*
sleep 2

# Update mirror recovery.confs to set restore point and unset primaryconninfo.
for dir in $MIRROR1 $MIRROR2 $MIRROR3 $STANDBY; do
#touch $dir/recovery.conf
echo "standby_mode = 'on'
recovery_target_name = 'foo'
primary_conninfo = ''" > $dir/recovery.conf
done

# Start standby master and mirrors. These replicas will automatically
# promote.
pg_ctl start -D $STANDBY
pg_ctl start -D $MIRROR1
pg_ctl start -D $MIRROR2
pg_ctl start -D $MIRROR3
sleep 5s

# Get rid of all rows in gp_segment_configuration to simulate new
# cluster initialized from archives.
PGOPTIONS="-c gp_session_role=utility" psql -p 7001 postgres -c "
    set allow_system_table_mods=true;
    truncate gp_segment_configuration;
    "

# Add standby and mirror segments back into gp_segment_configuration table.
PGOPTIONS="-c gp_session_role=utility" psql -p 7001 postgres -c "
    select pg_catalog.gp_add_segment(8::smallint, -1::smallint, 'p', 'p', 'n', 'u', 7001, '$(hostname)', '$(hostname)', '$STANDBY');
    select pg_catalog.gp_add_segment(5::smallint, 0::smallint, 'p', 'p', 'n', 'u', 7002, '$(hostname)', '$(hostname)', '$MIRROR1');
    select pg_catalog.gp_add_segment(6::smallint, 1::smallint, 'p', 'p', 'n', 'u', 7003, '$(hostname)', '$(hostname)', '$MIRROR2');
    select pg_catalog.gp_add_segment(7::smallint, 2::smallint, 'p', 'p', 'n', 'u', 7004, '$(hostname)', '$(hostname)', '$MIRROR3');
    "

# Set environment to use new cluster.
export PGPORT=7001
export MASTER_DATA_DIRECTORY=$STANDBY

# Restart the cluster to get the MPP parts working.
gpstop -air

# Run validation test.
run_tests gpdb_pitr_validate

## Print unnecessary success output.
echo "SUCCESS! GPDB Point-In-Time Recovery worked."
