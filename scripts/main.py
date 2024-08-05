from urllib.parse import quote
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker, declarative_base
import os
import time
from dotenv import load_dotenv
from base import Base
from agency import get_agency_data

load_dotenv()

password = os.getenv('MYSQL_PASSWORD')
password_quoted = quote(password)
DATABASE_URL = os.getenv('MYSQL_URL').replace("<password>", password_quoted)


def connect_to_database(url, max_attemps=10, delay_seconds=1):
    for attempt in range(max_attemps):
        try:
            engine = create_engine(url)
            engine.connect()
            Base.metadata.create_all(bind=engine)
            Session = sessionmaker(bind=engine)
            session = Session()
            print('Database connection successful')
            return session
        except Exception as e:
            print(f'--------------  Attempt {attempt + 1} fail. Retrying in {delay_seconds} seconds  ---------------------')
            print(e)
            time.sleep(delay_seconds)
    raise Exception('-----------------------  Failed to connect to the database after several attemps  ------------------------')
    return None


if __name__ == "__main__":
    session = connect_to_database(DATABASE_URL)
    if session:
        # get_agency_data(session)
        session.close()
