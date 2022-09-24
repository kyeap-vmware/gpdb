#include "pg_upgrade_greenplum.h"

#include "postgres_fe.h"

/*
 * Copies of backend xid comparison utilities
 */
#define InvalidTransactionId		((TransactionId) 0)
#define FirstNormalTransactionId	((TransactionId) 3)

#define TransactionIdIsValid(xid)		((xid) != InvalidTransactionId)
#define TransactionIdIsNormal(xid)	((xid) >= FirstNormalTransactionId)

/*
 * TransactionIdPrecedes --- is id1 logically < id2?
 */
static bool
TransactionIdPrecedes(TransactionId id1, TransactionId id2)
{
	/*
	 * If either ID is a permanent XID then we can just do unsigned
	 * comparison.  If both are normal, do a modulo-2^32 comparison.
	 */
	int32		diff;

	if (!TransactionIdIsNormal(id1) || !TransactionIdIsNormal(id2))
		return (id1 < id2);

	diff = (int32) (id1 - id2);
	return (diff < 0);
}

/*
 * GPDB5: Calculate the oldest xid in the old cluster by taking the minimum
 * datfrozenxid across all databases in the old cluster.
 */
void
set_old_cluster_chkpnt_oldstxid()
{
	TransactionId	oldestXid = InvalidTransactionId;

	for (int dbnum = 0; dbnum < old_cluster.dbarr.ndbs; dbnum++)
	{
		DbInfo *active_db = &old_cluster.dbarr.dbs[dbnum];
		if (!TransactionIdIsValid(oldestXid) ||
			TransactionIdPrecedes(active_db->datfrozenxid, oldestXid))
			oldestXid = active_db->datfrozenxid;
	}

	old_cluster.controldata.chkpnt_oldstxid = oldestXid;
}
