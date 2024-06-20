# pg_hint_plan

The `pg_hint_plan` module allows tweaking PostgreSQL execution plans with hints in SQL comments, such as `/*+ SeqScan(a) */`.

PostgreSQL's cost-based optimizer uses data statistics to estimate the costs of various execution plans and selects the lowest-cost plan. However, it may miss some data properties, such as column correlations, leading to suboptimal plans.

The `pg_hint_plan` module is a Greenplum Database extension.

## <a id="topic_reg"></a>Loading the Module 

To activate the `pg_hint_plan` module, run the following command as a superuser in a session:

```
LOAD 'pg_hint_plan';
```

## <a id="topic_info"></a>Additional Module Documentation 

Refer to the `pg_hint_plan` READMEs in the [GitHub repository](https://github.com/ossc-db/pg_hint_plan) for additional information about this module.
