#!/usr/bin/env bash


curl http://127.0.0.1:40403/json_rpc -d '
{ 
    "jsonrpc": "2.0", 
    "id": "1", 
    "method": "transfer", 
    "params": { 
            "transfers": [ 
                { 
                    "amount": 1,
                    "destination": "detoi1qxszv4ell3de4ur8lsrfys9k4dhgzw9kgu3cp25z4yu5382h6wtzc29pvfz92x774klw7y352euqwdxwuw", 
                    "payload_rpc": [ 
                        { 
                            "name":"D", 
                            "datatype":"U", 
                            "value":16045690981402826360
                        }, 
                        { 
                            "name":"NAME", 
                            "datatype":"S", 
                            "value":"dero" 
                        }, 
                        { 
                            "name":"CURR", 
                            "datatype":"S", 
                            "value":"usd" 
                        }
                    ] 
                } 
            ] 
        } 
}' -H 'Content-Type: application/json'
