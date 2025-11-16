import http from 'k6/http';
import { check } from 'k6';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

const serviceURL = 'http://localhost:8080';

export let options = {
    vus: 10,
    duration: '30s',
};

export function setup() {
    const authorId = uuidv4();

    const teamBody = {
        members: [
            {
                is_active: true,
                user_id: authorId,
                username: "username-" + uuidv4(),
            },
        ],
        team_name: "team-" + uuidv4(),
    };

    const res = http.post(
        `${serviceURL}/team/add`,
        JSON.stringify(teamBody),
        { headers: { 'Content-Type': 'application/json' } }
    );

    check(res, {
        "team created": (r) => r.status === 201,
    });

    return { authorId };
}

export default function (data) {
    const prId = uuidv4();

    const body = {
        author_id: data.authorId,
        pull_request_id: prId,
        pull_request_name: "PR-" + uuidv4()
    };

    const res = http.post(
        `${serviceURL}/pullRequest/create`,
        JSON.stringify(body),
        { headers: { 'Content-Type': 'application/json' } }
    );

    check(res, {
        'status is 201': (r) => r.status === 201,
    });
}
