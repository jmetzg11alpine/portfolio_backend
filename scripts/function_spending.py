from sqlalchemy import Column, Integer, String, Float
from base import Base
import requests
from collections import defaultdict
import csv
import os
current_dir = os.path.dirname(os.path.abspath(__file__))


class FunctionSpending(Base):
    __tablename__ = 'function_spending'
    id = Column(Integer, primary_key=True, index=True, autoincrement=True)
    year = Column(Integer)
    name = Column(String(255))
    amount = Column(Float)


def fetch_function_spending(years_to_check):
    functions_to_keep = ['Energy', 'Net Interest', 'Commerce and Housing Credit', 'Transportation', 'Agriculture', 'Health', 'Education, Training, Employment, and Social Services', 'National Defense']

    agency_codes = {}
    agency_file_name = os.path.join(current_dir, 'data', 'agency_codes.csv')
    with open(agency_file_name, mode='r') as file:
        reader = csv.DictReader(file)
        for row in reader:
            agency_codes[row['Agency Name']] = row['Agency Code']

    for year in years_to_check:
        function_spending = defaultdict(int)
        for agency_name, agency_code in agency_codes.items():
            try:
                print(year, agency_name)
                base_url = f"https://api.usaspending.gov/api/v2/agency/{agency_code}/budget_function/"
                params = {'fiscal_year': year}
                response = requests.get(base_url, params=params)
                results = response.json()['results']

                for result in results:
                    function_name = result['name']
                    if function_name in functions_to_keep:
                        amount = result['gross_outlay_amount']
                        function_spending[function_name] += amount
            except Exception as e:
                print(f'could not find for budget functions for {agency_name}, {e}')
        function_file_name = os.path.join(current_dir, 'data', f'functions_{year}.csv')
        with open(function_file_name, mode='w', newline='') as file:
            writer = csv.writer(file)
            writer.writerow(['Function', 'Amount'])
            for function, amount in function_spending.items():
                writer.writerow([function, amount])


def add_function_spending(session):
    for year in range(2017, 2025):
        function_file_name = os.path.join(current_dir, 'data', f'functions_{year}.csv')
        with open(function_file_name, mode='r') as file:
            reader = csv.DictReader(file)
            for row in reader:
                function_entry = FunctionSpending(year=year, name=row['Function'], amount=row['Amount'])
                session.add(function_entry)
    session.commit()
    print('function spending table updated')


def get_function_spending(session):
    # fetch_function_spending(range(2017, 2025))
    add_function_spending(session)
