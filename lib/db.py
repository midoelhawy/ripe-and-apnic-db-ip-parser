import os
import sqlite3

class SQLiteHandler:
    def __init__(self, db_name):
        self.db_name = db_name

    def create_table(self):
        conn = sqlite3.connect(self.db_name)
        c = conn.cursor()
        c.execute('''CREATE TABLE IF NOT EXISTS ip_data
                     (id INTEGER PRIMARY KEY,
                      first_ip TEXT,
                      last_ip TEXT,
                      first_ip_int INTEGER,
                      last_ip_int INTEGER,
                      ip_count INTEGER,
                      netname TEXT,
                      country TEXT,
                      descr TEXT,
                      mnt_by TEXT)''')
        c.execute('''CREATE INDEX IF NOT EXISTS idx_first_ip_int ON ip_data (first_ip_int)''')
        c.execute('''CREATE INDEX IF NOT EXISTS idx_last_ip_int ON ip_data (last_ip_int)''')
        conn.commit()
        conn.close()

    def insert_data(self, data):
        conn = sqlite3.connect(self.db_name)
        c = conn.cursor()
        for entry in data:
            c.execute('''INSERT INTO ip_data 
                         (first_ip, last_ip, first_ip_int, last_ip_int, ip_count, netname, country, descr, mnt_by)
                         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)''',
                         (entry['first_ip'], entry['last_ip'], entry['first_ip_int'], entry['last_ip_int'],
                          entry['ip_count'], entry.get('netname', None), entry['country'], entry.get('descr', None),
                          entry.get('mnt-by', None)))
        conn.commit()
        conn.close()
