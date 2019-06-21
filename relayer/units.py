def mmbtu_to_mwh(qty):
	return float(qty * 0.29307107)


def toe_to_mwh(qty):
	return float(qty * 11.63)


def gigajoule_to_mwh(qty):
	return float(qty / 3.6)


'''Emissions Unit Conversions'''

def lbs_to_mt(qty):
	return float(qty / 2204.62)


def convert_units(qty, unit):

	if unit == 'lbs':
		return lbs_to_mt(qty)

	elif unit == 'mmbtu':
		return mmbtu_to_mwh(qty)
	elif unit == 'toe':
		return toe_to_mwh(qty)

	else:
		raise Exception("Units not valid.")