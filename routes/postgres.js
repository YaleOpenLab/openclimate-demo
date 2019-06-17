var express = require('express');
var router = express.Router();
var users = require('./postgres/users')

router.get('/users', users.getUsers)
router.get('/users/:id', users.getUserById)
router.post('/users', users.createUser)
router.put('/users/:id', users.updateUser)
router.delete('/users/:id', users.deleteUser)

module.exports = router;
