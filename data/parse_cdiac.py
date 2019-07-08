import pandas as pd
import os, json, sys, re
import load_data

filename = 'csv_data/cdiac_fossil_fuel_cement_national.csv'


def parse_recent_data(df):
	return df[df.Year == 2014]


if __name__ == "__main__":
	df = pd.read_csv(filename)
	recent_data = parse_recent_data(df)
	load_data.export_to_json(recent_data, 'countries_emissions_2014')