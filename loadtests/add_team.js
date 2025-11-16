import http from 'k6/http';
import { check } from 'k6';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

const serviceURL = 'http://localhost:8080';

export let options = {
    vus: 100,
    duration: '30s',
};

export default function () {
    const teamName = 'team' + uuidv4();

    const body = {
        members: [],
        team_name: teamName,
    };
    const membersCount = Math.floor(Math.random() * (20 - 5 + 1)) + 5;
    for (let i = 0; i < membersCount; i++) {
        body.members.push({
            is_active: true,
            user_id: uuidv4(),
            username: "username" + uuidv4(),
        });
    }

    // ---------- POST /team/add ----------
    let res = http.post(
        `${serviceURL}/team/add`,
        JSON.stringify(body),
        { headers: { 'Content-Type': 'application/json' } }
    );

    check(res, {
        'status is 201': (r) => r.status === 201,
    });

    // ---------- GET /team/get ----------
    res = http.get(`${serviceURL}/team/get?team_name=${teamName}`);

    check(res, {
        'status is 200': (r) => r.status === 200,
        'correct team name': (r) => JSON.parse(r.body).team_name === teamName,
    });
}
