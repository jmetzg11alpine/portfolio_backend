from urllib.parse import quote
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker
import os
import time
from dotenv import load_dotenv
from base import Base
from agency import get_agency_data
from foreign_aid import get_foreign_aid
from function_spending import get_function_spending

load_dotenv()

password = os.getenv('MYSQL_PASSWORD')
password_quoted = quote(password)
DATABASE_URL = os.getenv('MYSQL_URL').replace("<password>", password_quoted)


def connect_to_database(url, max_attemps=10, delay_seconds=5):
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
            print(f'ttempt {attempt + 1} fail')
            print(e)
            time.sleep(delay_seconds)
    raise Exception('Failed to connect to the database after several attemps')
    return None


if __name__ == "__main__":
    session = connect_to_database(DATABASE_URL)
    if session:
        get_agency_data(session)
        get_foreign_aid(session)
        get_function_spending(session)
        session.close()
