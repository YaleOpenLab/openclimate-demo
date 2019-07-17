import pandas as pd
import os, json, sys, re


def export_to_json(df, filename):
	"""
	Exports the entire panda dataframe into json format and names
	the json file 'filename'.
	"""

	json_dirname = 'json_data/'
	json_data = df.to_json(json_dirname + filename + '.json', orient='index')


def export_all_to_json(df_list):
	for item in df_list:
		_, filename = os.path.split(item[1])
		filename = filename.split('.')[0]
		export_to_json(item[0], filename)


def load_file():
	"""
	Looks for the csv files in the local directory, reads in each 
	of the CSV files, and returns panda dataframes for each one.
	"""
	dirname, _ = os.path.split(os.path.abspath(sys.argv[0]))
	dirname = os.path.join(dirname, 'data')
	csv_files = ['data/' + file for file in os.listdir(dirname)]
	csv_dfs = [(pd.read_csv(csv_file), csv_file) for csv_file in csv_files]
	return csv_dfs


if __name__ == "__main__":
	csv_dfs = load_file()
	export_all_to_json(csv_dfs)