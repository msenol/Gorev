const path = require('path');
const fs = require('fs');

module.exports = {
  run: async function() {
    // Create the mocha test
    const Mocha = require('mocha');

    const mocha = new Mocha({
      ui: 'tdd',
      color: true,
      timeout: 20000
    });

    const testsRoot = path.resolve(__dirname, './unit');

    // Read all test files directly using fs
    const files = fs.readdirSync(testsRoot).filter(f => f.endsWith('.test.js'));

    // Add files to the test suite
    files.forEach(f => mocha.addFile(path.resolve(testsRoot, f)));

    // Run the mocha test
    return new Promise((resolve, reject) => {
      try {
        mocha.run(failures => {
          if (failures > 0) {
            reject(new Error(`${failures} tests failed.`));
          } else {
            resolve();
          }
        });
      } catch (err) {
        console.error(err);
        reject(err);
      }
    });
  }
};
