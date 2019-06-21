import sys
import pandas as pd
import datetime
from units import convert_units


column_names = [
	'asset_name', 
	'asset_type',
	'city',
	'region',
	'start_date', 
	'end_date',
	'co2e', 
	'co2e_unit',
	'certified', 
	'carbon_intensity', 
	'intensity_unit',
	'energy_generated',
	'energy_unit',
	]


def group_emissions_by_attr(data, attr):
	'''
	Aggregate CO2e emissions, grouped by the categories of a specified column.
	Specify the relevant column using a string.
	'''
	grouped = data.loc[:,'co2e'].groupby(data[attr])
	return grouped.sum()


def get_total_emissions(data):
	return data.loc[:,'co2e'].sum()	


def standardize_co2e_units(data):

	for i, row in data.iterrows():
		if row['co2e_unit'] != 'mt':
			co2e_val = convert_units(row['co2e'], row['co2e_unit'])
			data.at[i, 'co2e'] = co2e_val

	data.loc[:, 'co2e_unit'] = 'mt'
	return data


def standardize_energy_units(data):

	for i, row in data.iterrows():
		if row['energy_unit'] != 'MWh':
			energy_val = convert_units(row['energy_generated'], row['energy_unit'])
			data.at[i, 'energy_generated'] = energy_val

	data.loc[:, 'energy_unit'] = 'MWh'
	return data


def clean_data(data):
	'''
	Standardizes all units in the data to metric tons for CO2e emissions
	and MWh for energy generation.
	'''
	data = standardize_co2e_units(data)
	data = standardize_energy_units(data)
	return(data)


def load_file(file):
	'''Reads in a CSV file and returns a panda dataframe.'''
	return pd.read_csv(file, names=column_names)


def export_to_json(data):
	'''Exports the entire panda dataframe into json format.'''
	return data.to_json('./data.json')


def main():
	file = input("Enter csv filename: ")
	data = load_file(file)
	data = clean_data(data)

	return export_to_json(data)

	# test = group_emissions_by_attr(data, 'region')
	# print(test)

if __name__ == '__main__':
	main()




