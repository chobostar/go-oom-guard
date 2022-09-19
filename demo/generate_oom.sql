set work_mem='1GB';
explain analyze select a, max(b), min(c) from generate_series(1,100000000) as a, generate_series(1,100000000) as b, generate_series(1,10000000) as c group by a;