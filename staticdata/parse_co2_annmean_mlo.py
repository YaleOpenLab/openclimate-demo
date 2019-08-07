import pandas as pd
from load_data import export_to_json


if __name__ == '__main__':
	df = pd.read_csv('co2_annmean_mlo.csv', comment='#')
	dirname = 'earth_data/'
	filename = 'co2_annmean.json'
	df.to_json(dirname + filename, orient='index')