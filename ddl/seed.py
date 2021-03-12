import psycopg2 
from decouple import config, RepositoryEnv

from time import sleep  

PG_HOST: str = config("PG_HOST")
PG_PORT: int = config("PG_PORT")
PG_DATABASE: int = config("PG_DATABASE")
PG_USER: str = config("PG_USER")
PG_PASSWORD: int = config("PG_PASSWORD")

if __name__ == '__main__':
    sleep(5)
    conn = psycopg2.connect(
                            host=PG_HOST, 
                            dbname=PG_PORT,
                            user=PG_USER, 
                            password=PG_PASSWORD, 
                            port=PG_PORT
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