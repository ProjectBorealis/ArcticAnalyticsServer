const fs = require('fs').promises;
const path = require('path');

class ResultWriter {
    constructor(filename) {
        this.filename = filename;
        this.fileHandle = null; // Keep track of the file handle
    }

    // Opens the file (similar to Go's NewResultWriter)
    static async create(filename) {
        const writer = new ResultWriter(filename);
        writer.fileHandle = await fs.open(filename, 'a'); // Ensure file is created if not exists
        return writer;
    }

    // Append performance results and metadata
    async appendResults(metadata, results) {
        const data = { ...metadata, ...results };
        const jsonLine = JSON.stringify(data) + '\n';
        await fs.appendFile(this.filename, jsonLine); // Append data to file
    }

    getFilename() {
        return this.filename;
    }

    async close() {
        // Close the file handle if it's open
        if (this.fileHandle) {
            await this.fileHandle.close();
            this.fileHandle = null;
        }
    }
}

module.exports = ResultWriter;
