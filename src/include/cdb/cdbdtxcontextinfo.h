/*-------------------------------------------------------------------------
 *
 * cdbdtxcontextinfo.h
 *
 * Portions Copyright (c) 2007-2008, Greenplum inc
 * Portions Copyright (c) 2012-Present VMware, Inc. or its affiliates.
 *
 *
 * IDENTIFICATION
 *	    src/include/cdb/cdbdtxcontextinfo.h
 *
 *-------------------------------------------------------------------------
 */
#ifndef CDBDTXCONTEXTINFO_H
#define CDBDTXCONTEXTINFO_H

#include "access/xlog_internal.h" /* MAXFNAMELEN */
#include "utils/snapshot.h"

#define DtxContextInfo_StaticInit {InvalidDistributedTransactionId,false,GP_SNAPSHOT_MODE_LOCAL,.gpSnapshotInfo.distributedSnapshot=DistributedSnapshot_StaticInit,0,0,0,0}

typedef struct DtxContextInfo
{
	DistributedTransactionId 		distributedXid;
	
	bool							cursorContext;
	
	GpSnapshotMode						gpSnapshotMode;
	union {
		/* Distributed snapshot info, only for GP_SNAPSHOT_MODE_DISTRIBUTED*/
		DistributedSnapshot		 		distributedSnapshot;
		/* Restore point info, only for GP_SNAPSHOT_MODE_RESTOREPOINT */
		char 						rpname[MAXFNAMELEN];
	} gpSnapshotInfo;

	int 							distributedTxnOptions;

	uint32							segmateSync;
	uint32							nestingLevel;

	/* currentCommandId of QD, for debugging only */
	CommandId				 		curcid;	
} DtxContextInfo;

extern DtxContextInfo QEDtxContextInfo;	

extern void DtxContextInfo_Reset(DtxContextInfo *dtxContextInfo);
extern void DtxContextInfo_CreateOnCoordinator(DtxContextInfo *dtxContextInfo, bool inCursor,
										  int txnOptions, Snapshot snapshot);
extern int DtxContextInfo_SerializeSize(DtxContextInfo *dtxContextInfo);

extern void DtxContextInfo_Serialize(char *buffer, DtxContextInfo *dtxContextInfo);
extern void DtxContextInfo_Deserialize(const char *serializedDtxContextInfo,
									   int serializedDtxContextInfolen,
									   DtxContextInfo *dtxContextInfo);

extern void DtxContextInfo_Copy(DtxContextInfo *target, DtxContextInfo *source);
#endif   /* CDBDTXCONTEXTINFO_H */
