var express = require('express');
var router = express.Router();

/* GET home page. */
router.get('/', function(req, res, next) {
  res.json({ info: 'Greetings Earthlings' });
  // res.render('index', { title: 'Express' }); if we want to render some html
});

module.exports = router;
