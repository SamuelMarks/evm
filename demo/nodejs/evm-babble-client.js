http = require('http');

var EVMDAG1Client = function(host, port) {
    this.host = host
    this.port = port
}

request = function(options, callback) {
    return http.request(options, (resp) => {
        log(FgYellow, util.format('%s %s:%s%s', 
        options.method, 
        options.host,
        options.port,
        options.path));

        let data = '';
        
        // A chunk of data has been received.
        resp.on('data', (chunk) => {
            data += chunk;
        });
        
        // The whole response has been received. Process the result.
        resp.on('end', () => {
            callback(data);
        });   
    })
}

// class methods
EVMDAG1Client.prototype.getAccount = function(address) {
    var options = {
        host: this.host,
        port: this.port,
        path: '/account/' + address,
        method: 'GET'
      };
    
    return new Promise((resolve, reject) => {
        req = request(options, resolve)
        req.on('error', (err) => reject(err))
        req.end()
    })
} 

EVMDAG1Client.prototype.getAccounts = function() {
    var options = {
        host: this.host,
        port: this.port,
        path: '/accounts',
        method: 'GET'
      };
    
    return new Promise((resolve, reject) => {
        req = request(options, resolve)
        req.on('error', (err) => reject(err))
        req.end()
    })
}  

EVMDAG1Client.prototype.call = function(tx) {
    var options = {
        host: this.host,
        port: this.port,
        path: '/call',
        method: 'POST'
      };
    
    return new Promise((resolve, reject) => {
        req = request(options, resolve)
        req.write(tx)
        req.on('error', (err) => reject(err))
        req.end()
    })
} 

EVMDAG1Client.prototype.sendTx = function(tx) {
    var options = {
        host: this.host,
        port: this.port,
        path: '/tx',
        method: 'POST'
    };
  
    return new Promise((resolve, reject) => {
        req = request(options, resolve)
        req.write(tx)
        req.on('error', (err) => reject(err))
        req.end()
    })
}

EVMDAG1Client.prototype.sendRawTx = function(tx) {
    var options = {
        host: this.host,
        port: this.port,
        path: '/rawtx',
        method: 'POST'
    };
  
    return new Promise((resolve, reject) => {
        req = request(options, resolve)
        req.write(tx)
        req.on('error', (err) => reject(err))
        req.end()
    })
}

EVMDAG1Client.prototype.getReceipt = function(txHash) {
    var options = {
        host: this.host,
        port: this.port,
        path: '/tx/' + txHash,
        method: 'GET'
      };
    
    return new Promise((resolve, reject) => {
        req = request(options, resolve)
        req.on('error', (err) => reject(err))
        req.end()
    })
} 

module.exports = EVMDAG1Client;
