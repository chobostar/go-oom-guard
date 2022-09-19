#!/usr/bin/env python
import psycopg2
import os

# 210000 * 500 = 105000000 > 100Mb
TUPLE_SIZE=500
TOTAL_ROWS=210000

def str_generator(size=500):
    min_lc = ord(b'a')
    len_lc = 26
    ba = bytearray(os.urandom(size))
    for i, b in enumerate(ba):
        ba[i] = min_lc + b % len_lc
    return ba.decode('utf-8')

def setup_migrations(conn, cursor):
    cursor.execute("""
    CREATE TABLE IF NOT EXISTS test_table(
        id bigserial primary key,
        val text unique
    );
    """)
    conn.commit()

    cursor.execute("TRUNCATE TABLE test_table;")
    conn.commit()

def generate_oom(conn, cursor):
    n = TOTAL_ROWS
    values = [str_generator(size=TUPLE_SIZE) for _ in range(n)]
    # cursor.mogrify() to insert multiple values
    args = ','.join(cursor.mogrify("('{}')".format(i)).decode("utf-8") for i in values)

    # here will be OOM
    cursor.execute("INSERT INTO test_table(val) VALUES " + (args) + " ON CONFLICT DO NOTHING")
    conn.commit()

conn = psycopg2.connect(
    database="postgres",
    user='postgres',
    password='password',
    host='127.0.0.1',
    port='6432'
)

conn.autocommit = True
cursor = conn.cursor()

setup_migrations(conn, cursor)

generate_oom(conn, cursor)

conn.close()