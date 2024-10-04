const sqlite3 = require('sqlite3').verbose();
const path = require('path');

class ResultWriter {
    constructor(dbFilePath) {
        this.dbFilePath = dbFilePath;
        this.db = null;
    }

    // Opens the database and creates the table if it doesn't exist
    async initialize() {
        return new Promise((resolve, reject) => {
            this.db = new sqlite3.Database(this.dbFilePath, (err) => {
                if (err) {
                    console.error('Error opening database:', err.message);
                    return reject(err);
                } else {
                    this.db.run(`CREATE TABLE IF NOT EXISTS performance_results (
                        perfID INTEGER PRIMARY KEY AUTOINCREMENT,
                        data TEXT NOT NULL,
                        datetime TEXT NOT NULL,
                        ip TEXT NOT NULL
                    )`, (err) => {
                        if (err) {
                            console.error('Error creating table:', err.message);
                            return reject(err);
                        }
                        resolve(); // Resolve when the table is created
                    });
                }
            });
        });
    }

    // Static method to create an instance of ResultWriter
    static async create(dbFilePath) {
        const writer = new ResultWriter(dbFilePath);
        await writer.initialize(); // Initialize DB and create table if needed
        return writer;
    }

    async getResults() {
        return new Promise((resolve, reject) => {
            const query = 'SELECT * FROM performance_results';
            this.db.all(query, [], (err, rows) => {
                if (err) {
                    console.error('Error retrieving results from database:', err.message);
                    return reject(err);
                }
                resolve(rows); // Return the rows retrieved from the database
            });
        });
    }

    // Append performance results and metadata
    async appendResults(metadata, results) {
        const data = { ...metadata, ...results };
        const jsonData = JSON.stringify(data);
        const datetime = new Date().toISOString();
        const ip = metadata.ip;

        return new Promise((resolve, reject) => {
            this.db.run(`INSERT INTO performance_results (data, datetime, ip) VALUES (?, ?, ?)`, 
            [jsonData, datetime, ip], 
            function(err) {
                if (err) {
                    return reject(err);
                }
                resolve(this.lastID); // Returns the ID of the inserted row
            });
        });
    }

    async close() {
        if (this.db) {
            this.db.close((err) => {
                if (err) {
                    console.error('Error closing database:', err.message);
                }
            });
        }
    }
}

module.exports = ResultWriter;
