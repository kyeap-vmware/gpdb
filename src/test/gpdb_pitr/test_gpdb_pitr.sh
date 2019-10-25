#!/usr/bin/env bash

## ==================================================================
## Required: A fresh gpdemo cluster with mirrors sourced.
##
## This script tests and showcases a very simple Point-In-Time
## Recovery scenario by setting a restore point on all segments and
## invoking them on the standby master and mirror segments.
##
## Warning: This test will expropriate the gpdemo cluster
##          mirrors. Your gpdemo cluster will need to be recreated
##          after this test is run.
## ==================================================================

## Run tests via pg_regress with given test name(s)
run_tests()
{
    ../regress/pg_regress --dbname=gpdb_pitr_database --use-existing --init-file=../regress/init_file $1
    if [ $? != 0 ]; then
        exit 1
    fi
}

# Store standby master and mirror segment data directories.
DATADIR="${MASTER_DATA_DIRECTORY%*/*/*}"
STANDBY=${DATADIR}/standby
MIRROR1=${DATADIR}/dbfast_mirror1/demoDataDir0
MIRROR2=${DATADIR}/dbfast_mirror2/demoDataDir1
MIRROR3=${DATADIR}/dbfast_mirror3/demoDataDir2

# Create test database.
createdb gpdb_pitr_database

# Run setup test.
run_tests gpdb_pitr_setup

# Stop the cluster (send immediate shutdown to prevent CHECKPOINT).
gpstop -ai

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

# Set environment to use new cluster.
export PGPORT=7001
export MASTER_DATA_DIRECTORY=$STANDBY

# Reconfigure the segment configuration so that the replicas are
# recognized as primary segments.
PGOPTIONS="-c gp_session_role=utility" run_tests reconfigure_segment_config

# Restart the cluster to get the MPP parts working.
gpstop -air

# Run validation test.
run_tests gpdb_pitr_validate

## Print unnecessary success output.
echo "SUCCESS! GPDB Point-In-Time Recovery worked."
