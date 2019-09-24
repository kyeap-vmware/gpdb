/*-------------------------------------------------------------------------
 *
 * distributed_restore_point.c
 *
 * Distributed Restore Point
 *
 * Portions Copyright (c) 2019-Present Pivotal Software, Inc.
 *
 * src/backend/access/transam/distributed_restore_point.c
 *
 *-------------------------------------------------------------------------
 */

#include "postgres.h"

#include "access/distributed_restore_point.h"
#include "access/xlog.h"
#include "cdb/cdbdisp_query.h"
#include "cdb/cdbdispatchresult.h"
#include "cdb/cdbvars.h"
#include "utils/builtins.h"
#include "utils/timestamp.h"

static XLogRecPtr XLogDistributeRestorePoint(TimestampTz drp_dtime);

/*
 * Dispatcher function to create and log a distributed restore point.
 * This will internally dispatch a corresponding function to log a
 * distributed restore point in the executors.
 *
 * Returns the distributed timestamp (dtime) of the distributed restore point
 * XLOG record as type bigint.
 */
Datum
create_distributed_restore_point(PG_FUNCTION_ARGS)
{
	TimestampTz drp_dtime;
	char command[MAXFNAMELEN + 100];
	CdbPgResults cdb_pgresults = {NULL, 0};

	if (!IS_QUERY_DISPATCHER() && Gp_role != GP_ROLE_UTILITY)
		ereport(ERROR,
			(errcode(ERRCODE_GP_COMMAND_ERROR),
				(errmsg("must be invoked by a query dispatcher"))));

	if (!superuser())
		ereport(ERROR,
			(errcode(ERRCODE_INSUFFICIENT_PRIVILEGE),
				(errmsg("must be superuser to create a distributed restore point"))));

	if (RecoveryInProgress())
		ereport(ERROR,
			(errcode(ERRCODE_OBJECT_NOT_IN_PREREQUISITE_STATE),
				(errmsg("recovery is in progress"),
					errhint("WAL control functions cannot be executed during recovery."))));

	if (!XLogIsNeeded())
		ereport(ERROR,
			(errcode(ERRCODE_OBJECT_NOT_IN_PREREQUISITE_STATE),
				errmsg("WAL level not sufficient for creating a distributed restore point"),
				errhint("wal_level must be set to \"archive\" or \"hot_standby\" at server start.")));


	drp_dtime = GetCurrentTimestamp();
	sprintf(command, "SELECT pg_catalog.create_distributed_restore_point_on_segments(%ld)", drp_dtime);
	CdbDispatchCommand(command, DF_NONE, &cdb_pgresults);
	if (cdb_pgresults.numResults == 0)
		ereport(ERROR,
			(errmsg("did not receive results from query executors")));

	XLogDistributeRestorePoint(drp_dtime);

	PG_RETURN_INT64(drp_dtime);
}

/*
 * Executor function to create and log a distributed restore point.
 */
void
create_distributed_restore_point_on_segments(PG_FUNCTION_ARGS)
{
	if (Gp_role != GP_ROLE_EXECUTE)
		ereport(ERROR,
			(errcode(ERRCODE_GP_COMMAND_ERROR),
				(errmsg("must be invoked by a query executor"))));

	XLogDistributeRestorePoint(PG_GETARG_INT64(0));
}

static XLogRecPtr
XLogDistributeRestorePoint(TimestampTz drp_dtime)
{
	XLogRecPtr RecPtr;
	xl_distributed_restore_point xlrec;

	xlrec.drp_ltime = GetCurrentTimestamp();
	xlrec.drp_dtime = drp_dtime;
	xlrec.drp_contentid = GpIdentity.segindex;

	XLogBeginInsert();
	XLogRegisterData((char *) (&xlrec), sizeof(xl_distributed_restore_point));

	RecPtr = XLogInsert(RM_XLOG_ID, XLOG_DISTRIBUTED_RESTORE_POINT);

	ereport(LOG,
		(errmsg("distributed restore point \"%ld\" created at %X/%X",
			drp_dtime, (uint32) (RecPtr >> 32), (uint32) RecPtr)));

	return RecPtr;
}