var assert = require('assert');
var defaults = require('levelup-defaults');
var bytewise = require('bytewise');
var type = require('component-type');
var after = require('after');
var streamToArray = require('stream-to-array');

module.exports = Pathwise;

function Pathwise(db){
  assert(db, 'db required');
  this._db = defaults(db, {
    keyEncoding: bytewise,
    valueEncoding: 'json'
  });
}

Pathwise.prototype.put = function(path, obj, opts, fn){
  if (typeof opts == 'function') {
    fn = opts;
    opts = {};
  }
  var batch = opts.batch || this._db.batch();
  this._write(batch, path, obj, fn);
  if (opts.batch) setImmediate(fn);
  else batch.write(fn);
};

Pathwise.prototype._write = function(batch, key, obj, fn){
  var self = this;
  switch(type(obj)) {
    case 'object':
      var keys = Object.keys(obj);
      var next = after(keys.length, fn);
      keys.forEach(function(k){
        self._write(batch, key.concat(k), obj[k], next);
      });
      break;
    case 'array':
      this._write(batch, key, arrToObj(obj), fn);
      break;
    default:
      batch.put(bytewise.encode(key), JSON.stringify(obj));
      break;
  }
};

Pathwise.prototype.batch = function(ops, fn) {
  var self = this;
  var batch = this._db.batch();
  var next = after(ops.length, function(err){
    if (err) return fn(err);
    batch.write(fn);
  });
  ops.forEach(function(op){
    if (op.type == 'put') self.put(op.path, op.data, { batch: batch }, next);
    else if (op.type == 'del') self.del(op.path, { batch: batch }, next);
  });
};

Pathwise.prototype.get = function(path, fn){
  var ret = {};
  var el = ret;

  streamToArray(this._db.createReadStream({
    start: path,
    end: path.concat(undefined)
  }), function(err, data){
    if (err) return fn(err);

    data.forEach(function(kv){
      var segs = kv.key.slice(path.length);
      if (segs.length) {
        segs.forEach(function(seg, idx){
          if (!el[seg]) {
            if (idx == segs.length - 1) {
              el[seg] = kv.value;
            } else {
              el[seg] = {};
            }
          }
          el = el[seg];
        });
        el = ret;
      } else {
        ret = kv.value;
      }
    });
    fn(null, ret);
  });
};

Pathwise.prototype.del = function(path, opts, fn){
  if (typeof opts == 'function') {
    fn = opts;
    opts = {};
  }
  var batch = opts.batch || this._db.batch();

  streamToArray(this._db.createKeyStream({
    start: path,
    end: path.concat(undefined)
  }), function(err, keys){
    if (err) return fn(err);
    keys.forEach(function(key){ batch.del(bytewise.encode(key)) });
    if (opts.batch) fn();
    else batch.write(fn);
  });
};

Pathwise.prototype.children = function(path, fn) {
  streamToArray(this._db.createReadStream({
    start: path,
    end: path.concat(undefined)
  }), function(err, kv){
    if (err) return fn(err);
    fn(null, kv.map(function(_kv){
      return _kv.key[path.length] || _kv.value;
    }));
  });
}

Pathwise.prototype.someChildren = function(path,opts, fn) {
  var end = path.concat(undefined)
  if(opts.lte)end = path+opts.lte
  streamToArray(this._db.createReadStream({
    start: path+opts.gte,
    end: end
  }), function(err, kv){
    if (err) return fn(err);
    fn(null, kv.map(function(_kv){
      return _kv.key[path.length] || _kv.value;
    }));
  });
}


function arrToObj(arr){
  var obj = {};
  arr.forEach(function(el, idx){
    obj[idx] = el;
  });
  return obj;
}

