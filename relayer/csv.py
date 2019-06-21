import sys
import pandas as pd
import datetime
import json

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
	'capacity',
	'energy_generated',
	'energy_unit',
	]


def group_row_by_col(data, row, col):
	'''
	Aggregate rows, grouped by the categories of a specified column.
	Columns must be categorical variables. Specify the row and column names
	using strings.
	'''
	grouped = data.loc[:,row].groupby(data[col])
	return grouped.sum().to_json()


def get_totals(data):

	totals = {
		'total_co2e_emissions': str(data.loc[:,'co2e'].sum()),
		'total_energy_generated': str(data.loc[:,'energy_generated'].sum()),
	}

	return totals


def get_summary_data(data):

	emissions_by_region_df = group_row_by_col(data, 'co2e', 'region')
	emissions_by_asset_type_df = group_row_by_col(data, 'co2e', 'asset_type')

	energy_by_region_df = group_row_by_col(data, 'energy_generated', 'region')
	energy_by_asset_type_df = group_row_by_col(data, 'energy_generated', 'asset_type')

	capacity_by_asset_df = data.loc[:,['asset_name','capacity']]
	capacity_by_asset_df.loc[:,'capacity'].astype(str)
	capacity_by_asset_df = capacity_by_asset_df.to_json()

	summary_dict = {

		'emissions_by_region_df': emissions_by_region_df,
		'emissions_by_asset_type_df': emissions_by_asset_type_df,

		'energy_by_region_df': energy_by_region_df,
		'energy_by_asset_type_df': energy_by_asset_type_df,

		'capacity_by_asset_df': capacity_by_asset_df,

		'totals': get_totals(data),

	}

	return json.dumps(summary_dict)


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
	# file = input("Enter csv filename: ")
	file = 'example.csv'
	data = load_file(file)
	data = clean_data(data)

	print(get_summary_data(data))

	# return export_to_json(data)

	# test = group_emissions_by_attr(data, 'region')
	# print(test)

if __name__ == '__main__':
	main()




