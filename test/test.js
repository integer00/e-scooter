
import http from "k6/http";
import exec from 'k6/execution';

import { check, group, sleep } from 'k6';

export const options = {
  // iterations: 20,
  thresholds: {
    http_req_failed: [{
      threshold: 'rate<0.0001',
      // abortOnFail: true,
    }
  ], // http errors should be less than 1%
    http_req_duration: ['p(95)<200'], // 95% of requests should be below 200ms
  },
  scenarios: {
    scootersfull: {
      executor: 'constant-vus',
      vus: 1,
      duration: '30s',
    },
  },
  
};

const endpoint = "http://localhost:8080"
const users = new Map([
  [1, "alice"],
  [2, "bob"],
  [3, "me"],
]);

export default function () {
  const cookies = http.cookieJar();
  const username = users.get(exec.vu.idInInstance)


  group('get token', function() {

    console.log(username)

    const req_payload = JSON.stringify({
      userid: username
    });

    const req = http.post(endpoint + "/login",req_payload);

    console.log(req_payload)
    check(req, {
      'respond is 200': (r) => r.status === 200,
    });
  
    const cookiesForURL = cookies.cookiesForURL(req.url);
    check(null, {
      "have token": () => cookiesForURL.token.length > 0,
    });

    const req2 = http.get(endpoint + "/scooters")
    check(req2, {
      'should get scooters': (r) => r.status != 401,
    })
    console.log(req2.body)
  })

  group('scooter-pipeline', function() {
    const req3_payload = JSON.stringify({
      scooterid: "kappa_ride",
      userid: username
    });

    const req3 = http.post(endpoint + "/bookscooter",req3_payload)
    check(req3, {
      'scooter is booked': (r) => r.status === 200,
    })
    console.log(req3.body)

    const req4 = http.post(endpoint + "/start",req3_payload)
    check(req4, {
      'scooter started': (r) => r.status === 200,
    })
    console.log(req4.body)

    //ride for a while
    sleep(4)

    const req5 = http.post(endpoint + "/stop",req3_payload)
    check(req5, {
      'scooter is stopped': (r) => r.status === 200,
    })

    const req6 = http.get(endpoint + "/history")
    check(req6, {
      'history': (r) => r.status === 200,
    })

  })
}

