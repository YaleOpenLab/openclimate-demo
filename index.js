const level = require('level'),
    levelgraph = require('levelgraph'),
    jsonld = require('levelgraph-jsonld'),
    db = jsonld(levelgraph(level('./db', { createIfEmpty: true, base: 'https://schema.org' }))),
    fetch = require('node-fetch'),
    fs = require('fs-extra'),
    express = require('express'),
    cors = require('cors'),
    api = express(),
    port = process.env.PORT || 3001,
    http = require('http').Server(api),
    https_redirect = function(req, res, next) {
        if (process.env.NODE_ENV === 'production') {
            if (req.headers['x-forwarded-proto'] != 'https') {
                return res.redirect('https://' + req.headers.host + req.url);
            } else {
                return next();
            }
        } else {
            return next();
        }
    },
    iso3166 = require('./data/iso3166')

// API

api.use(https_redirect);
api.use(cors())
api.get('/dump', (req, res, next) => {
    var dump = {}
    res.setHeader('Content-Type', 'application/json')
    db.get([], function(err, obj) {
        dump = obj,
            res.send(JSON.stringify({
                dump
            }, null, 3))
    });
});
/*
api.get('/countries', (req, res, next) => {
    var countries = {}
    res.setHeader('Content-Type', 'application/json')
    store.get(['ref'], function(err, obj) {
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
        earth.annual_global_emissions = 'Lotsa GT'
        earth.avg_temp_increase = '0.1'
        earth.total_pledges = 4300
        earth.remaining_emissions = {
            remaining_budget: 'fewer GT',
            source: 'unk'
        }
        res.send(JSON.stringify({
            nation
        }, null, 3))
    });
});
api.get('/nations', (req, res, next) => { //list of country codes 
    res.setHeader('Content-Type', 'application/json')
    let list = {}
    let keys = Object.keys(ref.three_code)
    for (i = 0; i < keys.length; i++) {
        list[keys[i]] = ref.countries[ref.three_code[keys[i]]].f
    }
    res.send(JSON.stringify({
        list
    }, null, 3))


});
api.get('/nation/:id', (req, res, next) => { //country info and list subnational actors - only working for USA with 6 states
    let id = req.params.id.toUpperCase()
    res.setHeader('Content-Type', 'application/json')
    let rP = DBget2('info', id)
        //eP = DBget()
    Promise.all([rP])
        .then(function(v) {
            let nat = v[0],
                nation = {},
                coun = ref.countries[ref.three_code[id]]
            nation.id = coun.n
            nation.code = coun.c
            nation.display_name = coun.f
                //nation.polygon = coun.p
                //nation.flag_img = coun.f
            nation.subs = Object.keys(nat.sub)
            res.send(JSON.stringify({
                nation
            }, null, 3))
        });
});
*/
api.get('/nation/:id', (req, res, next) => { //country info and list subnational actors - only working for USA with 6 states
    let id = req.params.id.toUpperCase()
    res.setHeader('Content-Type', 'application/json')
    let num = parseInt(id)
    let query = ''
    if (typeof num == 'number') {
        query = '@ISO3166-1-numeric'
    } else if (id.length == 2) {
        query = '@ISO3166-1-Alpha-2'
    } else if (id.length == 3) {
        query = '@ISO3166-1-Alpha-3'
    }
    if (query) {
        db.jsonld.get(pld[id], { '@context': pld[query] }, function(err, obj) {
            if (err) {
                res.send(JSON.stringify({
                    err
                }, null, 3))
            } else {
                res.send(JSON.stringify({
                    obj
                }, null, 3))
            }
        });
    } else {
        res.send(JSON.stringify({
            Error: 'Not a valid query'
        }, null, 3))
    }
});
/*
api.get('/national-emissions/:id', (req, res, next) => {
    let id = req.params.id.toUpperCase()
    res.setHeader('Content-Type', 'application/json')
    let rP = DBget3('info', id, 'wby')
        //eP = DBget()
    Promise.all([rP])
        .then(function(v) {
            let em = v[0],
                nation = {},
                coun = ref.countries[ref.three_code[id]]
            nation.id = coun.n
            nation.code = coun.c
            nation.display_name = coun.f
                //nation.polygon = coun.p
                //nation.flag_img = coun.f
            nation.emissions = em
                //nation.emissions.source_name = em.sn
                //nation.emissions.source = em.s
                //nation.emissions.total_ghg_emissions = em.t
                //nation.emissions.land_based_sinks = em.s
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
*/
http.listen(port, function() {
    console.log(`DB API listening on port 3001`);
});

// DB
var ref = {}
for (i = 0; i < iso3166.length; i++) {
    var pld = {
        "@context": "https://schema.org",
        "@type": "Country",
        "@ISO3166-1-numeric": iso3166[i]['ISO3166-1-numeric'],
        "@ISO3166-1-Alpha-3": iso3166[i]['ISO3166-1-Alpha-3'],
        "@ISO3166-1-Alpha-2": iso3166[i]['ISO3166-1-Alpha-2'],
        "Country-data": iso3166[i]
    }
    db.jsonld.put(pld, function(err, obj) {
        if (err) {
            console.log(err)
        } else {
            console.log(obj)
        }
    });
}

/*
fs.readFile('csv/country-standards.csv')
    .then(function(file) {
        return file.toString('utf8')
    })
    .then(function(text) {
        let json = {
            num_to_three: {},
            two_to_three: {},
            lowercase_to_three: {},
            alts_to_three: {},
            three_code: {},
            countries: {}
        }
        const rows = text.split('\n')
        for (i = 1; i < rows.length; i++) {
            let col = rows[i].split(/,(?=(?:[^\"]*\"[^\"]*\")*[^\"]*$)/gm)
            if (col[5] != 'N/A') { //simple filter
                let now = parseInt(col[7])
                json.countries[col[0]] = {
                    c: col[5], //code, 3
                    f: col[1], //full name
                    a: [col[2].toLowerCase()], //alternate names
                    n: parseInt(col[0]), //numeric code
                    t: col[4], //two letter code
                    s: parseInt(col[6]), //start year
                    e: now
                }
                json.num_to_three[col[0]] = col[5]
                json.two_to_three[col[4]] = col[5]
                json.three_code[col[5]] = col[0]
                json.lowercase_to_three[col[2].toLowerCase()] = col[5]
                switch (col[0]) {
                    case '842':
                        json.lowercase_to_three['united states of america'] = col[5]
                        json.lowercase_to_three['usa'] = col[5]
                        break;
                    default:
                        break;
                }
            }
        }
        ref = json
        ref.errors = {}
        store.put(['ref'], json, function(err) {
            if (err) {
                console.log(err)
            } else {
                scrapeCSVCO('https://raw.githubusercontent.com/YaleOpenLab/openclimate/master/staticdata/csv_data/cdiac_fossil_fuel_cement_national.csv', 'https://raw.githubusercontent.com/YaleOpenLab/openclimate/master/staticdata/csv_data/consumption_emissions.csv')
            }
        })
        fs.readFile('csv/2017GHbyRegion.csv')
            .then(function(file) {
                return file.toString('utf8')
            })
            .then(function(text) {
                var data = text.split(/(?!\B"[^"]*)\n(?![^"]*"\B)/g),
                    header = data[0],
                    info = header.split(','),
                    schema = []
                ops = []
                for (i = 0; i < info.length; i++) {
                    schema.push(info[i])
                }
                for (i = 1; i < data.length - 1; i++) {
                    row = data[i].split(/,(?=(?:[^\"]*\"[^\"]*\")*[^\"]*$)/gm)
                    if (ref.lowercase_to_three[row[1].toLowerCase()] && row[1] == 'USA') {
                        var ghg = {
                            region: row[0].toLowerCase(),
                            name: row[0],
                            of: ref.lowercase_to_three[row[1].toLowerCase()],
                            year: parseInt(row[7]),
                            pop: row[5],
                            total: row[19]
                        }
                        ops.push({ type: 'put', path: ['info', ref.lowercase_to_three[row[1].toLowerCase()], 'sub', ghg.region], data: ghg })
                        ops.push({
                            type: 'put',
                            path: ['ref', 'countries', ref.three_code[ref.lowercase_to_three[row[1].toLowerCase()]], 'sub'],
                            data: {
                                [ghg.region]: 'regions'
                            }
                        })
                    } else {
                        if (row[1]) { ref.errors[row[1]] = 'regions' }
                    }
                }
                store.batch(ops)
                return
            })
    })
    .catch(function(e) {
        console.log(e)
    })
    */