var fetch = require('node-fetch');
var Pathwise = require('./pathwise');
var level = require('level');
const express = require('express');
const cors = require('cors');
const api = express()
var http = require('http').Server(api);
var cc = 'CC',
    ref = 'REF'
var store = new Pathwise(level('./db', { createIfEmpty: true }));
var https_redirect = function(req, res, next) {
    if (process.env.NODE_ENV === 'production') {
        if (req.headers['x-forwarded-proto'] != 'https') {
            return res.redirect('https://' + req.headers.host + req.url);
        } else {
            return next();
        }
    } else {
        return next();
    }
};

api.use(https_redirect);
api.use(cors())
api.get('/country-references', (req, res, next) => {
    var ref_codes = {}
    res.setHeader('Content-Type', 'application/json')
    store.get([ref], function(err, obj) {
        ref_codes = obj,
            res.send(JSON.stringify({
                ref_codes
            }, null, 3))
    });
});
api.get('/countries', (req, res, next) => {
    var countries = {}
    res.setHeader('Content-Type', 'application/json')
    store.get([cc], function(err, obj) {
        countries = obj,
            res.send(JSON.stringify({
                countries
            }, null, 3))
    });
});
api.get('/co2', (req, res, next) => {
    var co2 = {}
    res.setHeader('Content-Type', 'application/json')
    store.get(['CO2'], function(err, obj) {
        co2 = obj,
            res.send(JSON.stringify({
                co2
            }, null, 3))
    });
});
api.get('/earth', (req, res, next) => {
    res.setHeader('Content-Type', 'application/json')
    store.get(['stats'], function(err, obj) {
        let earth = {}
        earth.atmo_co2_concentration = '412PPM'
        earth.annual_global_emissions = ''
        earth.avg_temp_increase = ''
        earth.total_pledges = ''
        earth.remaining_emissions = {
            remaining_budget: '',
            source: ''
        }
        res.send(JSON.stringify({
            nation
        }, null, 3))
    });
});
api.get('/nation/:id', (req, res, next) => {
    let id = req.params.id
    let nation = {}
    res.setHeader('Content-Type', 'application/json')
    store.get([cc, id], function(err, obj) {
        nation.id = id
        nation.code = obj.c
        nation.display_name = obj.d
        nation.lower_name = obj.n
        nation.polygon = obj.p
        nation.flag_img = obj.f

        res.send(JSON.stringify({
            nation
        }, null, 3))
    });
});
api.get('/national-emissions/:id', (req, res, next) => {
    let id = req.params.id
    res.setHeader('Content-Type', 'application/json')
    let rP = DBget2(cc, id),
        eP = DBget(ref)
    Promise.all([rP, eP])
        .then(function(v) {
            let coun = v[0],
                em = v[1],
                nation = {}
            nation.id = id
            nation.code = coun.c
            nation.display_name = coun.d
            nation.lower_name = coun.n
            nation.polygon = coun.p
            nation.flag_img = coun.f
            nation.emissions = {}
            nation.emissions.source_name = em.sn
            nation.emissions.source = em.s
            nation.emissions.total_ghg_emissions = em.t
            nation.emissions.land_based_sinks = em.s
            nation.emissions.net_ghg = em.n
            res.send(JSON.stringify({
                nation
            }, null, 3))
        });
});
api.get('/national-pledges/:id', (req, res, next) => {
    let id = req.params.id
    res.setHeader('Content-Type', 'application/json')
    let rP = DBget2(cc, id),
        pP = DBget(ref)
    Promise.all([rP, pP])
        .then(function(v) {
            let coun = v[0],
                pl = v[1],
                nation = {}
            nation.id = id
            nation.code = coun.c
            nation.display_name = coun.d
            nation.lower_name = coun.n
            nation.polygon = coun.p
            nation.flag_img = coun.f
            nation.pledges = {}
            nation.pledges.name = 'Nationally Determined Contributions'
            nation.pledges.source = pl.s
            nation.pledges.baseline_year = pl.b
            nation.pledges.baseline_year_emissions = pl.e
            nation.pledges.target_year = pl.y
            nation.pledges.target_emissions = pl.t
            nation.pledges.emission_reduction_conditions = pl.c
            nation.pledges.emissio_reduction_per = pl.p
            res.send(JSON.stringify({
                nation
            }, null, 3))
        });
});
http.listen(3001, function() {
    console.log(`DB API listening on port 3001`);
});
scrapeCSVnBuildRef('https://raw.githubusercontent.com/openclimatedata/global-carbon-budget/master/data/country-definitions.csv', cc, 1, 2, ref)

function scrapeCSVnBuildRef(raw, location, by, alt, altl) { //url, db namespace, colnum, alt num, altnamespace
    fetch(`${raw}`)
        .then(function(response) {
            return response.text();
        })
        .then(function(text) {
            var header = text.split('\n')[0]
            var data = text.split('\n')
            var info = header.split(',')
            var json = {}
            var ref = {}
            var schema = []
            for (i = 0; i < info.length; i++) {
                schema.push(header[i])
            }
            for (i = 1; i < data.length; i++) {
                row = data[i].split(',')
                if (row[by]) {
                    var named = row[by].toLowerCase()
                    var refn = row[alt].toLowerCase()
                    json[named] = {}
                    ref[refn] = named
                    for (j = 0; j < info.length; j++) {
                        if (typeof row[j] === 'string') {
                            json[named][info[j]] = row[j].toLowerCase() || ''
                        } else {
                            json[named][info[j]] = row[j] || ''
                        }
                    }
                }
            }
            store.put([location], json, function(err) {
                if (err) {
                    console.log(err)
                } else {
                    store.put([altl], ref, function(err) {
                        if (err) {
                            console.log(err)
                        } else {
                            scrapeCSVCO('https://raw.githubusercontent.com/YaleOpenLab/openclimate/master/staticdata/csv_data/cdiac_fossil_fuel_cement_national.csv', 'https://raw.githubusercontent.com/YaleOpenLab/openclimate/master/staticdata/csv_data/consumption_emissions.csv')
                        }
                    })
                }
            })

        })
        .catch(function(e) {
            console.log(e)
        })
}

function scrapeCSVCO(raw, coem) { //built for cdiac_fossil_fuel_cement_national.csv & consumption_emissions.csv
    fetch(`${raw}`)
        .then(function(response) {
            return response.text();
        })
        .then(function(text) {
            var ccP = DBget(cc),
                refP = DBget(ref)
            Promise.all([ccP, refP])
                .then(function(v) {
                    var CC = v[0],
                        REF = v[1],
                        data = text.split('\n'),
                        header = data[0],
                        info = header.split(','),
                        json = {},
                        schema = []
                    for (i = 0; i < info.length; i++) {
                        schema.push(info[i])
                    }
                    for (i = 1; i < data.length; i++) {
                        row = data[i].split(',')
                        if (REF[row[0].toLowerCase()] == null) {
                            REF[row[0].toLowerCase()] = row[0].toLowerCase()
                        }
                        var named = REF[row[0].toLowerCase()]
                        if (json[named] == null) { json[named] = {} }
                        json[named][row[1]] = {}
                        for (j = 2; j < info.length; j++) {
                            json[named][row[1]][schema[j]] = row[j] || 0
                        }
                    }

                    return [CC, REF, json]
                })
                .then(function(x) {
                    var c = x[0],
                        r = x[1],
                        j = x[2],
                        ops = []
                    fetch(`${coem}`, j, r)
                        .then(function(response) {
                            return response.text();
                        })
                        .then(function(t) {
                            var data = t.split('\n'),
                                header = data[0],
                                info = header.split(','),
                                json = j,
                                schema = []
                            for (i = 0; i < info.length; i++) {
                                schema.push(info[i])
                            }
                            for (i = 1; i < data.length; i++) {
                                row = data[i].split(',')
                                if (json[row[0].toLowerCase()] == null) {
                                    json[row[0].toLowerCase()] = {}
                                }
                                if (json[row[0].toLowerCase()][row[1]] == null) {
                                    json[row[0].toLowerCase()][row[1]] = {}
                                }
                                json[row[0].toLowerCase()][row[1]][schema[2]] = parseFloat(row[2])
                            }
                            console.log(json)
                            ops.push({ type: 'put', path: ['REF'], data: r })
                            ops.push({ type: 'put', path: ['CO2'], data: json })
                            store.batch(ops)
                        })
                })
                .catch(function(e) {
                    console.log(e)
                })
        })
}

function DBget(loc) {
    return new Promise(function(resolve, reject) {
        store.get([loc], function(e, a) {
            if (e) {
                reject(e)
            } else if (Object.keys(a).length == 0) {
                resolve(0)
            } else {
                resolve(a)
            }
        });
    })
}

function DBget2(loc, loc2) {
    return new Promise(function(resolve, reject) {
        store.get([loc, loc2], function(e, a) {
            if (e) {
                reject(e)
            } else if (Object.keys(a).length == 0) {
                resolve(0)
            } else {
                resolve(a)
            }
        });
    })
}

function DBget3(loc, loc2, loc3) {
    return new Promise(function(resolve, reject) {
        store.get([loc, loc2, loc3], function(e, a) {
            if (e) {
                reject(e)
            } else if (Object.keys(a).length == 0) {
                resolve(0)
            } else {
                resolve(a)
            }
        });
    })
}

function DBget4(loc, loc2, loc3, loc4) {
    return new Promise(function(resolve, reject) {
        store.get([loc, loc2, loc3, loc4], function(e, a) {
            if (e) {
                reject(e)
            } else if (Object.keys(a).length == 0) {
                resolve(0)
            } else {
                resolve(a)
            }
        });
    })
}