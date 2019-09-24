/*
 * distributed_restore_point.h
 *
 * Distributed Restore Point
 *
 * Portions Copyright (c) 2019-Present Pivotal Software, Inc.
 *
 * src/include/access/distributed_restore_point.h
 */
#ifndef DISTRIBUTED_RESTORE_POINT_H
#define DISTRIBUTED_RESTORE_POINT_H

#include "fmgr.h"
#include "access/xlog_internal.h"

/* logs distributed restore point */
typedef struct xl_distributed_restore_point
{
	TimestampTz	drp_ltime;		/* local timestamp */
	TimestampTz	drp_dtime;		/* distributed timestamp from dispatcher */
	int32		drp_contentid;	/* segment id */
} xl_distributed_restore_point;

/* catalog functions to create distributed restore point */
extern Datum create_distributed_restore_point(PG_FUNCTION_ARGS);
extern void create_distributed_restore_point_on_segments(PG_FUNCTION_ARGS);

#endif /* DISTRIBUTED_RESTORE_POINT_H */
