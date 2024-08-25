import subprocess
import os

password = os.getenv('MYSQL_PASSWORD')
dump_file = "/sql_dump/mysql_dump.sql"


def create_sql_dump():
    try:
        dump_cmd = [
            "mysqldump",
            "-u", "root",
            f"-p{password}",
            "-h", "db",
            "portfolio"
        ]
        print(dump_cmd)
        with open(dump_file, 'w') as f:
            result = subprocess.run(dump_cmd, stdout=f, stderr=subprocess.PIPE)

        if result.returncode == 0:
            print(f"Database dump saved to {dump_file}")
        else:
            print(f"Failed to make database dump. Error: {result.stderr.decode('utf-8')}")

    except Exception as e:
        print(f"an error occured makeing dump: {e}")
