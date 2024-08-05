from sqlalchemy import Column, Integer, String, Float
from base import Base
import requests
import os
import csv
import datetime
current_dir = os.path.dirname(os.path.abspath(__file__))


class AgencyBudget(Base):
    __tablename__ = 'agency_budget'
    id = Column(Integer, primary_key=True, index=True, autoincrement=True)
    agency = Column(String(255))
    budget = Column(Float)


def fetch_agency_codes():
    base_url = "https://api.usaspending.gov/api/v2/agency/awards/count/"
    agency_codes = {}
    page = 1
    while True:
        response = requests.get(base_url, params={'limit': 100, 'page': page})
        data = response.json()
        results = data['results'][0] if data['results'] else []

        if not results:
            break

        for r in results:
            agency_codes[r['awarding_toptier_agency_name']] = r['awarding_toptier_agency_code']

        page += 1

    file_name = os.path.join(current_dir, 'data', 'agency_codes.csv')
    with open(file_name, mode='w', newline='') as file:
        writer = csv.writer(file)
        writer.writerow(['Agency Name', 'Agency Code'])
        for agency_name, agency_code in agency_codes.items():
            writer.writerow([agency_name, agency_code])


def fetch_budget_resources():
    agency_codes = {}
    agency_file_name = os.path.join(current_dir, 'data', 'agency_codes.csv')
    with open(agency_file_name, mode='r') as file:
        reader = csv.DictReader(file)
        for row in reader:
            agency_codes[row['Agency Name']] = row['Agency Code']

    agency_budgets = {}
    for name, code in agency_codes.items():
        try:
            print(f'budget resource for {name}')
            response = requests.get(f'https://api.usaspending.gov/api/v2/agency/{code}/budgetary_resources/')
            data = response.json()
            # response layout: https://api.usaspending.gov/api/v2/agency/012/budgetary_resources/
            # the index 0 has the most recent year
            budget = data['agency_data_by_year'][0]['agency_budgetary_resources']
            if budget:
                agency_budgets[name] = budget
        except Exception as e:
            print(f'could not find budget resource for {name}, error: {e}')

    budget_file_name = os.path.join(current_dir, 'data', 'agency_resources.csv')
    with open(budget_file_name, mode='w', newline='') as file:
        writer = csv.writer(file)
        writer.writerow(['Agency', 'Budget'])
        for agency, budget in agency_budgets.items():
            writer.writerow([agency, budget])


def record_updated_at():
    file_name = os.path.join(current_dir, 'data', 'data_updates.txt')
    today = str(datetime.date.today())
    with open(file_name, 'a') as file:
        file.write(f'\nagency budget data updated on: {today}')


def add_agency_budgets(session):
    agency_file_name = os.path.join(current_dir, 'data', 'agency_resources.csv')
    with open(agency_file_name, mode='r') as file:
        reader = csv.DictReader(file)
        for row in reader:
            agency_budget = AgencyBudget(agency=row['Agency'], budget=row['Budget'])
            session.add(agency_budget)
    session.commit()
    print('agency budget table updated')


def get_agency_data(session):
    # populates data/agency_codes.csv and data/agency_resources.csv
    fetch_agency_codes()
    fetch_budget_resources()
    record_updated_at()

    # adds data from the csv files to the database
    add_agency_budgets(session)
