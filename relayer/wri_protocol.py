import os, sys
import pandas as pd
import datetime
import json
import re

from units import convert_units


column_names = [
	'asset_name',
	'subnational',
	'country',
	'source_type',
	'activity',
	'scope1_CO2',
	'scope1_CH4',
	'scope1_N2O',
	'scope2_CO2',
	'scope2_CH4',
	'scope2_N2O',
]


def group_row_by_col(data, row, col):
	"""
	Aggregate rows, grouped by the categories of a specified column.
	Columns must be categorical variables. Specify the row and column names
	using strings.
	"""
	grouped = data.loc[:,row].groupby(data[col])
	return grouped.sum().to_json()


def get_summary_data(df):

	scope1_total = df[,'scope1_CO2'].sum() + df[,'scope1_CH4'].sum() + df[,'scope1_N2O'].sum()
	scope2_total = df[,'scope2_CO2'].sum() + df[,'scope2_CH4'].sum() + df[,'scope2_N2O'].sum()


def export_to_json(data):
	"""Exports the entire panda dataframe into json format."""
	return data.to_json('./data.json')


def load_file():
	"""
	Looks for the csv files in the local directory, reads in each 
	of the CSV files, and returns panda dataframes for each one.
	"""
	dirname, _ = os.path.split(os.path.abspath(sys.argv[0]))
	csv_files = [file for file in os.listdir(dirname) if re.match("[0-9]{4}.csv$", file)]
	csv_dfs = [pd.read(csv_file, names=column_names) for csv_file in csv_files]
	return csv_dfs


def main():
	csv_dfs = load_file()


if __name__ == "__main__":
	main()