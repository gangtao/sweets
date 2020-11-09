import http from "k6/http";
import { check, group, sleep } from "k6";
import { Rate } from "k6/metrics";

// A custom metric to track failure rates
var failureRate1 = new Rate("check_failure_rate1");
var failureRate2 = new Rate("check_failure_rate2");

function makeid(length) {
    var result           = '';
    var characters       = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    var charactersLength = characters.length;
    for ( var i = 0; i < length; i++ ) {
       result += characters.charAt(Math.floor(Math.random() * charactersLength));
    }
    return result;
 }

export default function() {

    const dataId = makeid(8)
    const group = makeid(8)
    const content = makeid(10)

    let response1 = http.post(`http://localhost:8080/sweets/v1/cs/configs?dataId=${dataId}&group=${group}&content=${content}`);
    // check() returns false if any of the specified conditions fail
    let checkRes1 = check(response1, {
        "status is 200": (r) => r.status === 200
    });
    // We reverse the check() result since we want to count the failures
    failureRate1.add(!checkRes1);

    let response2 = http.get(`http://localhost:8080/sweets/v1/cs/configs?dataId=${dataId}&group=${group}`);
    // check() returns false if any of the specified conditions fail
    let checkRes2 = check(response2, {
        "status is 200": (r) => r.status === 200,
        "content is present": (r) => r.body.indexOf(content) !== -1,
    });
    // We reverse the check() result since we want to count the failures
    failureRate2.add(!checkRes2);
};