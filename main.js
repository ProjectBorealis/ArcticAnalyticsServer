const fs = require('fs');
const path = require('path');
const dotenv = require('dotenv');
const ResultWriter = require('./server/resultswriter'); // Make sure the path is correct
const Server = require('./server/server');

// Load environment variables
dotenv.config();

const dbPath = path.join(__dirname, 'performance_results.db'); // SQLite database path
const port = 9095;

// Initialize the ResultWriter with the SQLite database path
(async () => {
    let resultWriter; // Declare resultWriter here for scope
    try {
        resultWriter = await ResultWriter.create(dbPath); // Use SQLite database instead of JSON
        const server = new Server(resultWriter, process.env.SHARED_SECRET, process.env.ADMIN_PASSWORD);

        await server.listen(port); // Wait for the server to start
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
