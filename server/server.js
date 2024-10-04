const express = require('express');
const crypto = require('crypto');
const fs = require('fs').promises;
const path = require('path');
const basicAuth = require('basic-auth');
const csvWriter = require('csv-writer').createObjectCsvWriter;

const { ResultWriter } = require('./resultswriter');
const { validateSchema } = require('./validate');

class Server {
    constructor(resultWriter, sharedSecret, adminPassword) {
        this.resultWriter = resultWriter;
        this.sharedSecret = sharedSecret;
        this.adminPassword = adminPassword;
        this.app = express();

        this.app.use(express.json());
        this.app.post('/v1/user/performance', this.postResultsHandler.bind(this));
        this.app.get('/v1/user/performance/csv', this.getResultsHandlerCSV.bind(this));
        this.app.get('/example', this.exampleHandler.bind(this));
    }

    async postResultsHandler(req, res) {
        const messageMAC = req.headers['authorization']; // Retrieve HMAC from the Authorization header
        const buf = JSON.stringify(req.body);

        // Validate results
        if (!validateSchema(buf)) {
            return res.status(400).send('Invalid Payload');
        }

        // Create HMAC for verification
        const hmac = crypto.createHmac('sha256', this.sharedSecret);
        hmac.update(buf);
        const expectedMAC = hmac.digest('hex');

        /* // Ensure both buffers have the same length before comparing
        if (Buffer.from(messageMAC).length !== Buffer.from(expectedMAC).length) {
            return res.status(403).send('Forbidden: Invalid HMAC length');
        }

        // Compare the HMACs
        if (!crypto.timingSafeEqual(Buffer.from(expectedMAC), Buffer.from(messageMAC))) {
            return res.status(403).send('Forbidden: Invalid HMAC');
        } */

        const metadata = {
            datetime: new Date().toISOString(),
            ip: req.ip,
        };

        try {
            await this.resultWriter.appendResults(metadata, req.body);
            res.sendStatus(200);
        } catch (error) {
            res.status(500).send('Error writing results');
        }
    }

    async getResultsHandlerCSV(req, res) {
        const credentials = basicAuth(req);
        if (!credentials || credentials.name !== 'admin' || credentials.pass !== this.adminPassword) {
            res.status(401).setHeader('WWW-Authenticate', 'Basic realm="Private"').send('Unauthorized');
            return;
        }

        const fileContents = await fs.readFile(this.resultWriter.getFilename(), 'utf-8');
        const results = fileContents.split('\n').filter(Boolean).map(JSON.parse);

        const headers = results.reduce((acc, result) => {
            result.events.forEach(event => {
                event.attributes.forEach(attr => {
                    acc.add(`${event.eventName}.${attr.name}`);
                });
            });
            return acc;
        }, new Set());

        const csvFilePath = path.join(__dirname, 'results.csv');
        const writer = csvWriter({
            path: csvFilePath,
            header: [{ id: 'sessionId', title: 'Session ID' }, ...Array.from(headers).map(h => ({ id: h, title: h }))],
        });

        await writer.writeRecords(results); // Convert JSON to CSV

        res.download(csvFilePath);
    }

    exampleHandler(req, res) {
        const performanceResults = {
            sessionId: `d1ac887243389d94544e4d9cc5524ab5-${new Date().toISOString()}`,
            userId: 'd1ac887243389d94544e4d9cc5524ab5',
            buildInfo: '1.0.0.7',
            events: [{
                eventName: 'Score',
                attributes: [{ name: 'Score.Num', value: '0' }]
            }]
        };

        const body = JSON.stringify(performanceResults);
        const hmac = crypto.createHmac('sha256', this.sharedSecret);
        hmac.update(body);

        const authHeader = hmac.digest('hex');
        res.send(`curl -v -H 'Content-Type: application/json' -H 'Authorization: ${authHeader}' http://${req.headers.host}/v1/user/performance -d '${body}'`);
    }

    listen(port) {
        this.app.listen(port, () => {
            console.log(`Server listening on port ${port}`);
        });
    }
}

module.exports = Server;
