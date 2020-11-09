import http from "k6/http";
import { check, group, sleep } from "k6";
import { Rate } from "k6/metrics";

// A custom metric to track failure rates
var failureRate = new Rate("check_failure_rate");

export default function() {
    let response = http.post("http://localhost:8080/sweets/v1/cs/configs?dataId=mydata&group=mygroup&content=hello");

    // check() returns false if any of the specified conditions fail
    let checkRes = check(response, {
        "status is 200": (r) => r.status === 200
    });

    // We reverse the check() result since we want to count the failures
    failureRate.add(!checkRes);
};