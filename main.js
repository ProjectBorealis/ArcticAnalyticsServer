const fs = require('fs');
const path = require('path');
const dotenv = require('dotenv');
const ResultWriter = require('./server/resultswriter');
const Server = require('./server/server');

// Load environment variables
dotenv.config();

const dataDir = './results'; // Example data directory
const port = 9095;

// Ensure the results directory exists
if (!fs.existsSync(dataDir)) {
    fs.mkdirSync(dataDir, { recursive: true });
}

(async () => {
    let resultWriter; // Declare resultWriter here for scope
    try {
        resultWriter = await ResultWriter.create(path.join(dataDir, 'results.json'));
        const server = new Server(resultWriter, process.env.SHARED_SECRET, process.env.ADMIN_PASSWORD);

        server.listen(port);
        console.log(`Listening on port ${port}...`);

        // Graceful shutdown handling
        process.on('SIGTERM', async () => {
            console.log('Shutting down...');
            await server.close(); // Assuming server.close is asynchronous
            console.log('Server closed.');

            await resultWriter.close(); // Close the result writer explicitly
            console.log('Result writer closed.');
            process.exit(0);
        });
    } catch (err) {
        console.error('Error starting server:', err);
    }
})();
