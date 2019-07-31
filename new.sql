--
-- Greenplum Database database dump
--

SET gp_default_storage_options = '';
SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;

--
-- Name: partition_test; Type: SCHEMA; Schema: -; Owner: bdoil
--

CREATE SCHEMA partition_test;


ALTER SCHEMA partition_test OWNER TO bdoil;

--
-- Name: s"'/;,1s; Type: SCHEMA; Schema: -; Owner: bdoil
--

CREATE SCHEMA "s""'/;,1s";


ALTER SCHEMA "s""'/;,1s" OWNER TO bdoil;

--
-- Name: s2; Type: SCHEMA; Schema: -; Owner: bdoil
--

CREATE SCHEMA s2;


ALTER SCHEMA s2 OWNER TO bdoil;

--
-- Name: schema2; Type: SCHEMA; Schema: -; Owner: bdoil
--

CREATE SCHEMA schema2;


ALTER SCHEMA schema2 OWNER TO bdoil;

--
-- Name: sm"1; Type: SCHEMA; Schema: -; Owner: bdoil
--

CREATE SCHEMA "sm""1";


ALTER SCHEMA "sm""1" OWNER TO bdoil;

--
-- Name: sm1; Type: SCHEMA; Schema: -; Owner: bdoil
--

CREATE SCHEMA sm1;


ALTER SCHEMA sm1 OWNER TO bdoil;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: p3_sales; Type: TABLE; Schema: partition_test; Owner: bdoil; Tablespace: 
--

CREATE TABLE partition_test.p3_sales (
    id integer,
    year integer,
    month integer,
    day integer,
    region text
)
 DISTRIBUTED BY (id) PARTITION BY RANGE(year)
          SUBPARTITION BY RANGE(month)
                  SUBPARTITION BY LIST(region) 
          (
          START (2011) END (2012) EVERY (1) WITH (tablename='p3_sales_1_prt_2', appendonly='false')
                  (
                  START (1) END (2) EVERY (1) WITH (tablename='p3_sales_1_prt_2_2_prt_2', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_2_2_prt_2_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_2_2_prt_2_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_2_2_prt_2_3_prt_other_regions', appendonly='false')
                          ), 
                  START (2) END (3) EVERY (1) WITH (tablename='p3_sales_1_prt_2_2_prt_3', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_2_2_prt_3_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_2_2_prt_3_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_2_2_prt_3_3_prt_other_regions', appendonly='false')
                          ), 
                  START (3) END (4) EVERY (1) WITH (tablename='p3_sales_1_prt_2_2_prt_4', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_2_2_prt_4_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_2_2_prt_4_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_2_2_prt_4_3_prt_other_regions', appendonly='false')
                          ), 
                  START (4) END (5) EVERY (1) WITH (tablename='p3_sales_1_prt_2_2_prt_5', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_2_2_prt_5_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_2_2_prt_5_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_2_2_prt_5_3_prt_other_regions', appendonly='false')
                          ), 
                  START (5) END (6) EVERY (1) WITH (tablename='p3_sales_1_prt_2_2_prt_6', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_2_2_prt_6_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_2_2_prt_6_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_2_2_prt_6_3_prt_other_regions', appendonly='false')
                          ), 
                  START (6) END (7) EVERY (1) WITH (tablename='p3_sales_1_prt_2_2_prt_7', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_2_2_prt_7_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_2_2_prt_7_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_2_2_prt_7_3_prt_other_regions', appendonly='false')
                          ), 
                  START (7) END (8) EVERY (1) WITH (tablename='p3_sales_1_prt_2_2_prt_8', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_2_2_prt_8_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_2_2_prt_8_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_2_2_prt_8_3_prt_other_regions', appendonly='false')
                          ), 
                  START (8) END (9) EVERY (1) WITH (tablename='p3_sales_1_prt_2_2_prt_9', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_2_2_prt_9_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_2_2_prt_9_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_2_2_prt_9_3_prt_other_regions', appendonly='false')
                          ), 
                  START (9) END (10) EVERY (1) WITH (tablename='p3_sales_1_prt_2_2_prt_10', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_2_2_prt_10_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_2_2_prt_10_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_2_2_prt_10_3_prt_other_regions', appendonly='false')
                          ), 
                  START (10) END (11) EVERY (1) WITH (tablename='p3_sales_1_prt_2_2_prt_11', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_2_2_prt_11_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_2_2_prt_11_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_2_2_prt_11_3_prt_other_regions', appendonly='false')
                          ), 
                  START (11) END (12) EVERY (1) WITH (tablename='p3_sales_1_prt_2_2_prt_12', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_2_2_prt_12_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_2_2_prt_12_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_2_2_prt_12_3_prt_other_regions', appendonly='false')
                          ), 
                  START (12) END (13) EVERY (1) WITH (tablename='p3_sales_1_prt_2_2_prt_13', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_2_2_prt_13_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_2_2_prt_13_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_2_2_prt_13_3_prt_other_regions', appendonly='false')
                          ), 
                  DEFAULT SUBPARTITION other_months  WITH (tablename='p3_sales_1_prt_2_2_prt_other_months', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_2_2_prt_other_months_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_2_2_prt_other_months_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_2_2_prt_other_months_3_prt_other_regions', appendonly='false')
                          )
                  ), 
          START (2012) END (2013) EVERY (1) WITH (tablename='p3_sales_1_prt_3', appendonly='false')
                  (
                  START (1) END (2) EVERY (1) WITH (tablename='p3_sales_1_prt_3_2_prt_2', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_3_2_prt_2_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_3_2_prt_2_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_3_2_prt_2_3_prt_other_regions', appendonly='false')
                          ), 
                  START (2) END (3) EVERY (1) WITH (tablename='p3_sales_1_prt_3_2_prt_3', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_3_2_prt_3_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_3_2_prt_3_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_3_2_prt_3_3_prt_other_regions', appendonly='false')
                          ), 
                  START (3) END (4) EVERY (1) WITH (tablename='p3_sales_1_prt_3_2_prt_4', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_3_2_prt_4_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_3_2_prt_4_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_3_2_prt_4_3_prt_other_regions', appendonly='false')
                          ), 
                  START (4) END (5) EVERY (1) WITH (tablename='p3_sales_1_prt_3_2_prt_5', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_3_2_prt_5_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_3_2_prt_5_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_3_2_prt_5_3_prt_other_regions', appendonly='false')
                          ), 
                  START (5) END (6) EVERY (1) WITH (tablename='p3_sales_1_prt_3_2_prt_6', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_3_2_prt_6_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_3_2_prt_6_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_3_2_prt_6_3_prt_other_regions', appendonly='false')
                          ), 
                  START (6) END (7) EVERY (1) WITH (tablename='p3_sales_1_prt_3_2_prt_7', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_3_2_prt_7_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_3_2_prt_7_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_3_2_prt_7_3_prt_other_regions', appendonly='false')
                          ), 
                  START (7) END (8) EVERY (1) WITH (tablename='p3_sales_1_prt_3_2_prt_8', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_3_2_prt_8_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_3_2_prt_8_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_3_2_prt_8_3_prt_other_regions', appendonly='false')
                          ), 
                  START (8) END (9) EVERY (1) WITH (tablename='p3_sales_1_prt_3_2_prt_9', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_3_2_prt_9_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_3_2_prt_9_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_3_2_prt_9_3_prt_other_regions', appendonly='false')
                          ), 
                  START (9) END (10) EVERY (1) WITH (tablename='p3_sales_1_prt_3_2_prt_10', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_3_2_prt_10_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_3_2_prt_10_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_3_2_prt_10_3_prt_other_regions', appendonly='false')
                          ), 
                  START (10) END (11) EVERY (1) WITH (tablename='p3_sales_1_prt_3_2_prt_11', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_3_2_prt_11_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_3_2_prt_11_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_3_2_prt_11_3_prt_other_regions', appendonly='false')
                          ), 
                  START (11) END (12) EVERY (1) WITH (tablename='p3_sales_1_prt_3_2_prt_12', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_3_2_prt_12_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_3_2_prt_12_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_3_2_prt_12_3_prt_other_regions', appendonly='false')
                          ), 
                  START (12) END (13) EVERY (1) WITH (tablename='p3_sales_1_prt_3_2_prt_13', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_3_2_prt_13_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_3_2_prt_13_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_3_2_prt_13_3_prt_other_regions', appendonly='false')
                          ), 
                  DEFAULT SUBPARTITION other_months  WITH (tablename='p3_sales_1_prt_3_2_prt_other_months', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_3_2_prt_other_months_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_3_2_prt_other_months_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_3_2_prt_other_months_3_prt_other_regions', appendonly='false')
                          )
                  ), 
          DEFAULT PARTITION outlying_years  WITH (tablename='p3_sales_1_prt_outlying_years', appendonly='false')
                  (
                  START (1) END (2) EVERY (1) WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_2', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_2_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_2_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_2_3_prt_other_regions', appendonly='false')
                          ), 
                  START (2) END (3) EVERY (1) WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_3', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_3_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_3_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_3_3_prt_other_regions', appendonly='false')
                          ), 
                  START (3) END (4) EVERY (1) WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_4', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_4_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_4_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_4_3_prt_other_regions', appendonly='false')
                          ), 
                  START (4) END (5) EVERY (1) WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_5', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_5_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_5_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_5_3_prt_other_regions', appendonly='false')
                          ), 
                  START (5) END (6) EVERY (1) WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_6', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_6_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_6_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_6_3_prt_other_regions', appendonly='false')
                          ), 
                  START (6) END (7) EVERY (1) WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_7', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_7_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_7_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_7_3_prt_other_regions', appendonly='false')
                          ), 
                  START (7) END (8) EVERY (1) WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_8', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_8_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_8_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_8_3_prt_other_regions', appendonly='false')
                          ), 
                  START (8) END (9) EVERY (1) WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_9', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_9_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_9_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_9_3_prt_other_regions', appendonly='false')
                          ), 
                  START (9) END (10) EVERY (1) WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_10', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_10_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_10_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_10_3_prt_other_regions', appendonly='false')
                          ), 
                  START (10) END (11) EVERY (1) WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_11', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_11_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_11_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_11_3_prt_other_regions', appendonly='false')
                          ), 
                  START (11) END (12) EVERY (1) WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_12', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_12_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_12_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_12_3_prt_other_regions', appendonly='false')
                          ), 
                  START (12) END (13) EVERY (1) WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_13', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_13_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_13_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_13_3_prt_other_regions', appendonly='false')
                          ), 
                  DEFAULT SUBPARTITION other_months  WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_other_months', appendonly='false')
                          (
                          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_other_months_3_prt_usa', appendonly='false'), 
                          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_other_months_3_prt_asia', appendonly='false'), 
                          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales_1_prt_outlying_years_2_prt_other_m_3_prt_other_regions', appendonly='false')
                          )
                  )
          );
 ALTER TABLE partition_test.p3_sales ALTER PARTITION outlying_years 
SET SUBPARTITION TEMPLATE 
          (
          SUBPARTITION usa VALUES('usa') WITH (tablename='p3_sales'), 
          SUBPARTITION asia VALUES('asia') WITH (tablename='p3_sales'), 
          DEFAULT SUBPARTITION other_regions  WITH (tablename='p3_sales')
          );

ALTER TABLE partition_test.p3_sales 
SET SUBPARTITION TEMPLATE 
          (
          START (1) END (2) EVERY (1) WITH (tablename='p3_sales') , 
          START (2) END (3) EVERY (1) WITH (tablename='p3_sales') , 
          START (3) END (4) EVERY (1) WITH (tablename='p3_sales') , 
          START (4) END (5) EVERY (1) WITH (tablename='p3_sales') , 
          START (5) END (6) EVERY (1) WITH (tablename='p3_sales') , 
          START (6) END (7) EVERY (1) WITH (tablename='p3_sales') , 
          START (7) END (8) EVERY (1) WITH (tablename='p3_sales') , 
          START (8) END (9) EVERY (1) WITH (tablename='p3_sales') , 
          START (9) END (10) EVERY (1) WITH (tablename='p3_sales') , 
          START (10) END (11) EVERY (1) WITH (tablename='p3_sales') , 
          START (11) END (12) EVERY (1) WITH (tablename='p3_sales') , 
          START (12) END (13) EVERY (1) WITH (tablename='p3_sales') , 
          DEFAULT SUBPARTITION other_months  WITH (tablename='p3_sales')
          )
;


ALTER TABLE partition_test.p3_sales OWNER TO bdoil;

--
-- Name: foo; Type: TABLE; Schema: public; Owner: bdoil; Tablespace: 
--

CREATE TABLE public.foo (
    i integer
)
 DISTRIBUTED BY (i);


ALTER TABLE public.foo OWNER TO bdoil;

--
-- Name: holds; Type: TABLE; Schema: public; Owner: bdoil; Tablespace: 
--

CREATE TABLE public.holds (
    i integer
)
 DISTRIBUTED BY (i);


ALTER TABLE public.holds OWNER TO bdoil;

--
-- Name: sales; Type: TABLE; Schema: public; Owner: bdoil; Tablespace: 
--

CREATE TABLE public.sales (
    id integer,
    date date,
    amt numeric(10,2)
)
 DISTRIBUTED BY (id) PARTITION BY RANGE(date) 
          (
          PARTITION jan17 START ('2017-01-01'::date) END ('2017-02-01'::date) WITH (tablename='sales_1_prt_jan17', appendonly='false' ), 
          PARTITION feb17 START ('2017-02-01'::date) END ('2017-03-01'::date) WITH (tablename='sales_1_prt_feb17', appendonly='false' ), 
          PARTITION mar17 START ('2017-03-01'::date) END ('2017-04-01'::date) WITH (tablename='sales_1_prt_mar17', appendonly='false' ), 
          PARTITION apr17 START ('2017-04-01'::date) END ('2017-05-01'::date) WITH (tablename='sales_1_prt_apr17', appendonly='false' ), 
          PARTITION may17 START ('2017-05-01'::date) END ('2017-06-01'::date) WITH (tablename='sales_1_prt_may17', appendonly='false' ), 
          PARTITION jun17 START ('2017-06-01'::date) END ('2017-07-01'::date) WITH (tablename='sales_1_prt_jun17', appendonly='false' ), 
          PARTITION jul17 START ('2017-07-01'::date) END ('2017-08-01'::date) WITH (tablename='sales_1_prt_jul17', appendonly='false' ), 
          PARTITION aug17 START ('2017-08-01'::date) END ('2017-09-01'::date) WITH (tablename='sales_1_prt_aug17', appendonly='false' ), 
          PARTITION sep17 START ('2017-09-01'::date) END ('2017-10-01'::date) WITH (tablename='sales_1_prt_sep17', appendonly='false' ), 
          PARTITION oct17 START ('2017-10-01'::date) END ('2017-11-01'::date) WITH (tablename='sales_1_prt_oct17', appendonly='false' ), 
          PARTITION nov17 START ('2017-11-01'::date) END ('2017-12-01'::date) WITH (tablename='sales_1_prt_nov17', appendonly='false' ), 
          PARTITION dec17 START ('2017-12-01'::date) END ('2018-01-01'::date) WITH (tablename='sales_1_prt_dec17', appendonly='false')
          );
ALTER TABLE public.sales EXCHANGE PARTITION dec17 WITH TABLE public.sales_1_prt_dec17_external_partition__ WITHOUT VALIDATION; 
DROP TABLE public.sales_1_prt_dec17_external_partition__; 


ALTER TABLE public.sales OWNER TO bdoil;

--
-- Name: t"'/;,1.t; Type: TABLE; Schema: s"'/;,1s; Owner: bdoil; Tablespace: 
--

CREATE TABLE "s""'/;,1s"."t""'/;,1.t" (
    "c""'/;,1.c" integer NOT NULL,
    c2 integer,
    c3 integer NOT NULL
)
 DISTRIBUTED BY ("c""'/;,1.c");


ALTER TABLE "s""'/;,1s"."t""'/;,1.t" OWNER TO bdoil;

--
-- Name: t"'/;,1.t_c"'/;,1.c_seq; Type: SEQUENCE; Schema: s"'/;,1s; Owner: bdoil
--

CREATE SEQUENCE "s""'/;,1s"."t""'/;,1.t_c""'/;,1.c_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE "s""'/;,1s"."t""'/;,1.t_c""'/;,1.c_seq" OWNER TO bdoil;

--
-- Name: t"'/;,1.t_c"'/;,1.c_seq; Type: SEQUENCE OWNED BY; Schema: s"'/;,1s; Owner: bdoil
--

ALTER SEQUENCE "s""'/;,1s"."t""'/;,1.t_c""'/;,1.c_seq" OWNED BY "s""'/;,1s"."t""'/;,1.t"."c""'/;,1.c";


--
-- Name: t"'/;,1.t_c3_seq; Type: SEQUENCE; Schema: s"'/;,1s; Owner: bdoil
--

CREATE SEQUENCE "s""'/;,1s"."t""'/;,1.t_c3_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE "s""'/;,1s"."t""'/;,1.t_c3_seq" OWNER TO bdoil;

--
-- Name: t"'/;,1.t_c3_seq; Type: SEQUENCE OWNED BY; Schema: s"'/;,1s; Owner: bdoil
--

ALTER SEQUENCE "s""'/;,1s"."t""'/;,1.t_c3_seq" OWNED BY "s""'/;,1s"."t""'/;,1.t".c3;


--
-- Name: t0; Type: TABLE; Schema: s2; Owner: bdoil; Tablespace: 
--

CREATE TABLE s2.t0 (
    a integer,
    b text
)
 DISTRIBUTED BY (a);


ALTER TABLE s2.t0 OWNER TO bdoil;

--
-- Name: foo2; Type: TABLE; Schema: schema2; Owner: bdoil; Tablespace: 
--

CREATE TABLE schema2.foo2 (
    i integer
)
 DISTRIBUTED BY (i);


ALTER TABLE schema2.foo2 OWNER TO bdoil;

--
-- Name: foo3; Type: TABLE; Schema: schema2; Owner: bdoil; Tablespace: 
--

CREATE TABLE schema2.foo3 (
    i integer
)
 DISTRIBUTED BY (i);


ALTER TABLE schema2.foo3 OWNER TO bdoil;

--
-- Name: foo4; Type: TABLE; Schema: schema2; Owner: bdoil; Tablespace: 
--

CREATE TABLE schema2.foo4 (
    c1 text
)
 DISTRIBUTED BY (c1);


ALTER TABLE schema2.foo4 OWNER TO bdoil;

--
-- Name: returns; Type: TABLE; Schema: schema2; Owner: bdoil; Tablespace: 
--

CREATE TABLE schema2.returns (
    id integer,
    date date,
    amt numeric(10,2)
)
 DISTRIBUTED BY (id) PARTITION BY RANGE(date) 
          (
          PARTITION jan17 START ('2017-01-01'::date) END ('2017-02-01'::date) WITH (tablename='returns_1_prt_jan17', appendonly='false' ), 
          PARTITION feb17 START ('2017-02-01'::date) END ('2017-03-01'::date) WITH (tablename='returns_1_prt_feb17', appendonly='false' ), 
          PARTITION mar17 START ('2017-03-01'::date) END ('2017-04-01'::date) WITH (tablename='returns_1_prt_mar17', appendonly='false' ), 
          PARTITION apr17 START ('2017-04-01'::date) END ('2017-05-01'::date) WITH (tablename='returns_1_prt_apr17', appendonly='false' ), 
          PARTITION may17 START ('2017-05-01'::date) END ('2017-06-01'::date) WITH (tablename='returns_1_prt_may17', appendonly='false' ), 
          PARTITION jun17 START ('2017-06-01'::date) END ('2017-07-01'::date) WITH (tablename='returns_1_prt_jun17', appendonly='false' ), 
          PARTITION jul17 START ('2017-07-01'::date) END ('2017-08-01'::date) WITH (tablename='returns_1_prt_jul17', appendonly='false' ), 
          PARTITION aug17 START ('2017-08-01'::date) END ('2017-09-01'::date) WITH (tablename='returns_1_prt_aug17', appendonly='false' ), 
          PARTITION sep17 START ('2017-09-01'::date) END ('2017-10-01'::date) WITH (tablename='returns_1_prt_sep17', appendonly='false' ), 
          PARTITION oct17 START ('2017-10-01'::date) END ('2017-11-01'::date) WITH (tablename='returns_1_prt_oct17', appendonly='false' ), 
          PARTITION nov17 START ('2017-11-01'::date) END ('2017-12-01'::date) WITH (tablename='returns_1_prt_nov17', appendonly='false' ), 
          PARTITION dec17 START ('2017-12-01'::date) END ('2018-01-01'::date) WITH (tablename='returns_1_prt_dec17', appendonly='false' )
          );


ALTER TABLE schema2.returns OWNER TO bdoil;

--
-- Name: tbl0; Type: TABLE; Schema: schema2; Owner: bdoil; Tablespace: 
--

CREATE TABLE schema2.tbl0 (
    a integer,
    b text
)
 DISTRIBUTED BY (a);


ALTER TABLE schema2.tbl0 OWNER TO bdoil;

--
-- Name: t1; Type: TABLE; Schema: sm"1; Owner: bdoil; Tablespace: 
--

CREATE TABLE "sm""1".t1 (
    a integer
)
 DISTRIBUTED BY (a);


ALTER TABLE "sm""1".t1 OWNER TO bdoil;

--
-- Name: small1; Type: TABLE; Schema: sm1; Owner: bdoil; Tablespace: 
--

CREATE TABLE sm1.small1 (
    a integer
)
 DISTRIBUTED BY (a);


ALTER TABLE sm1.small1 OWNER TO bdoil;

--
-- Name: t"'/;,1.t; Type: TABLE; Schema: sm1; Owner: bdoil; Tablespace: 
--

CREATE TABLE sm1."t""'/;,1.t" (
    a integer
)
 DISTRIBUTED BY (a);


ALTER TABLE sm1."t""'/;,1.t" OWNER TO bdoil;

--
-- Name: t".1; Type: TABLE; Schema: sm1; Owner: bdoil; Tablespace: 
--

CREATE TABLE sm1."t"".1" (
    a integer
)
 DISTRIBUTED BY (a);


ALTER TABLE sm1."t"".1" OWNER TO bdoil;

--
-- Name: t'.1; Type: TABLE; Schema: sm1; Owner: bdoil; Tablespace: 
--

CREATE TABLE sm1."t'.1" (
    a integer
)
 DISTRIBUTED BY (a);


ALTER TABLE sm1."t'.1" OWNER TO bdoil;

--
-- Name: t.1; Type: TABLE; Schema: sm1; Owner: bdoil; Tablespace: 
--

CREATE TABLE sm1."t.1" (
    a integer
)
 DISTRIBUTED BY (a);


ALTER TABLE sm1."t.1" OWNER TO bdoil;

--
-- Name: t.2; Type: TABLE; Schema: sm1; Owner: bdoil; Tablespace: 
--

CREATE TABLE sm1."t.2" (
    a integer
)
 DISTRIBUTED BY (a);


ALTER TABLE sm1."t.2" OWNER TO bdoil;

--
-- Name: tbl0; Type: TABLE; Schema: sm1; Owner: bdoil; Tablespace: 
--

CREATE TABLE sm1.tbl0 (
    a integer,
    b text
)
 DISTRIBUTED BY (a);


ALTER TABLE sm1.tbl0 OWNER TO bdoil;

--
-- Name: tbl1; Type: TABLE; Schema: sm1; Owner: bdoil; Tablespace: 
--

CREATE TABLE sm1.tbl1 (
    a integer,
    b text
)
 DISTRIBUTED BY (a);


ALTER TABLE sm1.tbl1 OWNER TO bdoil;

--
-- Name: c"'/;,1.c; Type: DEFAULT; Schema: s"'/;,1s; Owner: bdoil
--

ALTER TABLE ONLY "s""'/;,1s"."t""'/;,1.t" ALTER COLUMN "c""'/;,1.c" SET DEFAULT nextval('"s""''/;,1s"."t""''/;,1.t_c""''/;,1.c_seq"'::regclass);


--
-- Name: c3; Type: DEFAULT; Schema: s"'/;,1s; Owner: bdoil
--

ALTER TABLE ONLY "s""'/;,1s"."t""'/;,1.t" ALTER COLUMN c3 SET DEFAULT nextval('"s""''/;,1s"."t""''/;,1.t_c3_seq"'::regclass);


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: bdoil
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM bdoil;
GRANT ALL ON SCHEMA public TO bdoil;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- Greenplum Database database dump complete
--

