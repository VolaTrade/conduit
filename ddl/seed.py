import psycopg2 
from time import sleep  

if __name__ == '__main__':
    sleep(5)
    conn = psycopg2.connect(
                            host="conduit_postgres", 
                            dbname="postgres",
                            user="conduiter", 
                            password="password", 
                            port=5432
                        )


    cur = conn.cursor()

    with open("ddl.sql", "r") as fr:
        seed_data: str = fr.read()

        try:
            print("seeding -> ", seed_data)
            cur.execute(seed_data)
            print("Cursor executed")
            conn.commit()
            print("Committed sql transaction")

        finally:
            conn.close()