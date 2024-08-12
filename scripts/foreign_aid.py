from sqlalchemy import Column, Integer, String, Float
from base import Base
import os
import csv
import json
from collections import defaultdict
current_dir = os.path.dirname(os.path.abspath(__file__))


class ForeignAid(Base):
    __tablename__ = 'foreign_aid'
    id = Column(Integer, primary_key=True, index=True, autoincrement=True)
    country = Column(String(255))
    amount = Column(Float)
    year = Column(Integer)
    lat = Column(Float)
    lng = Column(Float)


replacements = {
    'China (Hong Kong, S.A.R., P.R.C.)': 'China',
    'China (P.R.C.)': 'China',
    'China (Tibet)': 'China',
    'Burma (Myanmar)': 'Myanmar',
    'Congo (Brazzaville)': 'Republic of the Congo',
    'Congo (Kinshasa)': 'DR Congo',
    'Korea, Democratic Republic of': 'North Korea',
    'Korea, Republic of': 'South Korea',
    'Sao Tome and Principe': 'São Tomé and Príncipe',
    'Slovak Republic': 'Slovakia',
    'Sudan (former)': 'Sudan',
    'West Bank and Gaza': 'Palestine',
    'Micronesia, Federated States of': 'Micronesia',
    'Kiribati, Republic of': 'Kiribati',
    'Curacao': 'Curaçao',
    'Serbia and Montenegro (former)': 'Serbia'
}


exclude = ['World']


def prep_foreign_aid_data():
    lat_lng_mapper_file = os.path.join(current_dir, 'data', 'lat_lng_mapper.json')
    with open(lat_lng_mapper_file, 'r') as json_file:
        lat_lng_mapper = json.load(json_file)
    data = defaultdict(lambda: defaultdict(float))

    foreign_aid_file = os.path.join(current_dir, 'data', 'foreign_aid_all.csv')
    with open(foreign_aid_file, 'r') as csv_file:
        csv_reader = csv.DictReader(csv_file)
        for row in csv_reader:
            country = row['Country Name']
            year = row['Fiscal Year']
            amount = row['Constant Dollar Amount']
            country = replacements.get(country, country)
            try:
                if country not in exclude and 'Region' not in country:
                    data[country][year] += int(amount)
            except Exception as e:
                print(country, e)

    entries = []
    for country, years in data.items():
        for year, amount in years.items():
            entries.append({
                'country': country,
                'year': year,
                'amount': amount,
                'lat': lat_lng_mapper[country]['lat'],
                'lng': lat_lng_mapper[country]['lng']
            })
    return entries


def get_foreign_aid(session):
    entries = prep_foreign_aid_data()

    aids = [
        ForeignAid(
            country=entry['country'],
            year=entry['year'],
            amount=entry['amount'],
            lat=entry['lat'],
            lng=entry['lng'])
        for entry in entries
    ]

    session.bulk_save_objects(aids)
    session.commit()

    print('foreign aid table updated')


# to get the lat and lng of countries.
# works in a jupyter notebook but not in the fastAPI app
# def fetch_country_latlng():
#     url = "https://restcountries.com/v3.1/all"
#     response = requests.get(url)
#     countries = response.json()

#     name_mapper = {
#         'Cape Verde': 'Cabo Verde',
#         'Ivory Coast': "Cote d'Ivoire",
#         'Saint Kitts and Nevis': 'St. Kitts and Nevis',
#         'Saint Lucia': 'St. Lucia',
#         'Saint Vincent and the Grenadines': 'St. Vincent and Grenadines'
#     }

#     lat_lng_mapper = {}
#     for country in countries:
#         name = country['name']['common']
#         if name in name_mapper:
#             name = name_mapper[name]
#         latlng = country['latlng']
#         lat_lng_mapper[country] = {'lat': latlng[0], 'lng': latlng[1]}

#     with open('lat_lng_mapper.json', 'w') as json_file:
#         json.dump(lat_lng_mapper, json_file)
