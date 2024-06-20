# Query Optimization

Greenplum employs two specialized optimizers: Postgres-Based Planner and GPORCA, each tailored to specific types of workloads:

- **Postgres-based Planner** is optimized for transactional workloads. 
- **GPORCA**  is designed for analytical and hybrid transactional-analytical workloads. 

For a given user query, the optimizers explore a vast search space of equivalent execution plans. They utilize table statistics and a cardinality estimation model to predict the number of rows involved in each operation of the query. Based on these predictions, a cost estimation model assigns a cost value to each of the execution plans. The execution plan with the lowest estimated cost is then selected as the optimal plan.

## <a id="optimizer"></a> Optimizer Hints

In database query optimization, optimizer hints are directives provided to the query optimizer to influence the execution plan of SQL queries. These hints can override the optimizer's default behavior, addressing specific issues such as inaccurate row estimates, suboptimal scan operations, inappropriate join types, and inefficient join orders. This topic explores various categories of optimizer hints and discusses common issues they aim to resolve.

>**Note** Greenplum currently does not support motion hints.

### <a id="cardinality"></a> Cardinality Hints

When the optimizer inaccurately estimates the number of rows output from joins, it can lead to suboptimal query plans, such as choosing broadcast over redistribution or preferring MergeJoin over HashJoin.

Cardinality hints adjust the estimated number of rows output from query operations, particularly useful when the optimizer's estimates are inaccurate due to missing statistics or stale data.

**Example:**

```sql
/*+ Rows(t1 t2 t3 #42) */ SELECT * FROM t1, t2, t3; -- sets row estimate to 42
/*+ Rows(t1 t2 t3 +42) */ SELECT * FROM t1, t2, t3; -- adds 42 to original row estimate
/*+ Rows(t1 t2 t3 -42) */ SELECT * FROM t1, t2, t3; -- subtracts 42 from original row estimate
/*+ Rows(t1 t2 t3 *42) */ SELECT * FROM t1, t2, t3; -- multiplies 42 with original row estimate
```

### <a id="table-access"></a> Table Access Hints

Stale statistics or misestimation of cost/cardinality can cause the optimizer to choose inefficient scan operations. Global User Configurations (GUCs) are often too broad a solution for varying query requirements.

Table access hints specify the scan type and index to use for tables in a query, helping to optimize data retrieval strategies.

**Example:**

```sql
/*+ SeqScan(t1) */ SELECT * FROM t1 WHERE t1.a > 42; -- force sequential scan
/*+ IndexScan(t1 my_index) */ SELECT * FROM t1 WHERE t1.a > 42; -- force index scan
/*+ IndexOnlyScan(t1) */ SELECT * FROM t1 WHERE t1.a > 42; -- force index-only scan
/*+ BitmapScan(t1 my_bitmap_index) */ SELECT * FROM t1 WHERE t1.a > 42; -- force bitmap index scan
```

### <a id="join-type"></a> Join Type Hints

Unexpected computational skew may cause hash joins to spill onto disk, affecting query performance. Users may know that index nested loop joins perform better for specific queries.

Join type hints specify the physical join operator to use, overriding the optimizer's choice based on user knowledge or specific query characteristics.

**Example:**

```sql
/*+ HashJoin(t1 t2) */ SELECT * FROM t1, t2; -- force hash join
/*+ NestLoop(t1 t2) */ SELECT * FROM t1, t2; -- force nested loop join
/*+ MergeJoin(t1 t2) */ SELECT * FROM t1 FULL JOIN t2 ON t1.a = t2.a; -- force merge join
```

### <a id="join-order"></a> Join Order Hints

Inefficient join orders selected by the optimizer can stem from insufficient statistics, unknown correlated columns, or bugs in the optimizer. Given the vast number of potential join orders, the optimizer does not evaluate every possibility.

Join order hints specify the preferred join order of tables in a query, guiding the optimizer to avoid inefficient join strategies.

**Example:**

```sql
/*+ Leading(t1 t2 t3) */ SELECT * FROM t1, t2, t3;
/*+ Leading(t1 (t3 t2)) */ SELECT * FROM t1, t2, t3;
```

### <a id="join-order"></a> Plan Hints Logging

To log used and unused hints, set `pg_hint_plan.debug_print` to `on` and `client_min_messages` to `log`.

**Example:**

```sql
set pg_hint_plan.debug_print = on;
set client_min_messages = 'log';
/*+SeqScan(t1)IndexScan(t2)*/
EXPLAIN (COSTS false) SELECT * FROM t1, t2 WHERE t1.id = t2.id;
LOG:  pg_hint_plan:
used hint:
SeqScan(t1)
IndexScan(t2)
not used hint:
duplication hint:
error hint:

                 QUERY PLAN                 
--------------------------------------------
 Hash Join
   Hash Cond: (t1.id = t2.id)
   ->  Seq Scan on t1
   ->  Hash
         ->  Index Scan using t2_pkey on t2
(5 rows) 
```

### <a id="best-practice"></a> Best Practices for Optimizer Hints

Best practices for using optimizer hints include but not limited to:

Use Hints to Mitigate Bad Plan Choices
:   Apply hints to address specific issues such as inaccurate row estimates, suboptimal scan operations, inappropriate join types, or inefficient join orders.

Test Plan Performance
:   Before applying hints in production, thoroughly test query performance to ensure the selected plan improves overall execution time and resource usage.

Use Plan Hints as Temporary Workarounds
:   While hints can resolve immediate performance issues, they should be considered temporary solutions. Regularly revisit and refine hints as underlying data characteristics evolve or data volume change.

Log Issues for Permanent Fixes
:   Where possible, log issues with the optimizer's behavior, so VMware can address the underlying issues with optimizer algorithms or statistics to reduce the reliance on hints for performance tuning.

Manage Conflicts with GUCs
:   Global User Configurations (GUCs) can override optimizer hints, potentially causing conflicts. For instance, a hint specifying an index scan `/*+ IndexScan(mytable myindex) */` will be overwritten by the GUC `SET optimizer_enable_indexscan=off;` Ensure that GUC settings align with hint directives. 

>**Note** `pg_hint_plan` allows you to specify GUCs inside the hint.
> 
> For example:
>
> ```
> /*Set(enable_indexscan off)*/
> EXPLAIN (COSTS false) SELECT * FROM t1, t2 WHERE t1.id = t2.id;
> ```
