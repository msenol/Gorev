#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const https = require('https');
const os = require('os');

/**
 * Postinstall script to download Gorev binaries from GitHub releases
 * This runs automatically after npm install
 */

const GITHUB_REPO = 'msenol/Gorev';
const BINARIES_DIR = path.join(__dirname, 'binaries');

/**
 * Safely unlink a file without throwing errors on Windows
 */
function safeUnlink(filePath) {
    try {
        if (fs.existsSync(filePath)) {
            // On Windows, try to make file writable first
            if (os.platform() === 'win32') {
                try {
                    fs.chmodSync(filePath, 0o666);
                } catch (e) {
                    // Ignore chmod errors
                }
            }
            fs.unlinkSync(filePath);
        }
    } catch (err) {
        // Ignore unlink errors - file might be in use or permission denied
        console.warn(`Warning: Could not remove file ${filePath}: ${err.message}`);
    }
}

function getPlatformInfo() {
    const platform = os.platform();
    const arch = os.arch();

    let goPlatform, goArch, binaryName, downloadName;

    switch (platform) {
        case 'win32':
            goPlatform = 'windows';
            binaryName = 'gorev.exe';
            break;
        case 'darwin':
            goPlatform = 'darwin';
            binaryName = 'gorev';
            break;
        case 'linux':
            goPlatform = 'linux';
            binaryName = 'gorev';
            break;
        default:
            throw new Error(`Unsupported platform: ${platform}`);
    }

    switch (arch) {
        case 'x64':
            goArch = 'amd64';
            break;
        case 'arm64':
            goArch = 'arm64';
            break;
        default:
            throw new Error(`Unsupported architecture: ${arch}`);
    }

    downloadName = `gorev-${goPlatform}-${goArch}${platform === 'win32' ? '.exe' : ''}`;
    const platformDir = `${goPlatform}-${goArch}`;

    return {
        platform: goPlatform,
        arch: goArch,
        binaryName,
        downloadName,
        platformDir,
        fullPlatform: `${goPlatform}-${goArch}`
    };
}

function downloadFile(url, dest) {
    return new Promise((resolve, reject) => {
        console.log(`Downloading ${url}...`);

        const file = fs.createWriteStream(dest);

        const request = https.get(url, (response) => {
            // Handle redirects
            if (response.statusCode === 301 || response.statusCode === 302) {
                const redirectUrl = response.headers.location;
                console.log(`Redirected to: ${redirectUrl}`);
                return downloadFile(redirectUrl, dest).then(resolve).catch(reject);
            }

            if (response.statusCode !== 200) {
                fs.unlinkSync(dest);
                reject(new Error(`Download failed with status code: ${response.statusCode}`));
                return;
            }

            response.pipe(file);

            file.on('finish', () => {
                file.close();
                console.log('Download completed successfully');
                resolve();
            });
        });

        request.on('error', (err) => {
            safeUnlink(dest);
            reject(err);
        });

        file.on('error', (err) => {
            safeUnlink(dest);
            reject(err);
        });
    });
}

async function getLatestReleaseVersion() {
    return new Promise((resolve, reject) => {
        const options = {
            hostname: 'api.github.com',
            path: `/repos/${GITHUB_REPO}/releases/latest`,
            headers: {
                'User-Agent': 'gorev-npm-installer'
            }
        };

        const req = https.request(options, (res) => {
            let data = '';

            res.on('data', (chunk) => {
                data += chunk;
            });

            res.on('end', () => {
                try {
                    const release = JSON.parse(data);
                    if (release.tag_name) {
                        resolve(release.tag_name);
                    } else {
                        reject(new Error('No tag_name found in release data'));
                    }
                } catch (err) {
                    reject(new Error(`Failed to parse release data: ${err.message}`));
                }
            });
        });

        req.on('error', reject);
        req.end();
    });
}

async function downloadBinary(platformInfo, version) {
    const { downloadName, platformDir, binaryName } = platformInfo;

    // Create binaries directory structure
    const targetDir = path.join(BINARIES_DIR, platformDir);
    fs.mkdirSync(targetDir, { recursive: true });

    const binaryPath = path.join(targetDir, binaryName);

    // Check if binary already exists in the package (bundled binaries)
    if (fs.existsSync(binaryPath)) {
        // Verify the binary is executable and has correct version
        try {
            if (process.platform !== 'win32') {
                fs.chmodSync(binaryPath, 0o755);
            }

            // Check binary version to ensure it matches package version
            const { execSync } = require('child_process');
            try {
                const binaryVersion = execSync(`"${binaryPath}" version`, {
                    encoding: 'utf8',
                    timeout: 5000,
                    stdio: ['ignore', 'pipe', 'ignore']
                }).trim();

                // Extract version from output (format: "Gorev vX.Y.Z" or just "vX.Y.Z")
                const versionMatch = binaryVersion.match(/v?\d+\.\d+\.\d+/);
                const extractedVersion = versionMatch ? versionMatch[0] : null;

                // Remove 'v' prefix from both for comparison
                const normalizedBinaryVersion = extractedVersion ? extractedVersion.replace(/^v/, '') : null;
                const normalizedPackageVersion = version.replace(/^v/, '');

                if (normalizedBinaryVersion === normalizedPackageVersion) {
                    console.log(`✅ Binary already exists with correct version: ${binaryPath} (v${normalizedBinaryVersion})`);
                    console.log('Skipping download (using bundled binary)...');
                    return; // Skip download if binary exists with correct version
                } else {
                    console.log(`⚠️  Existing binary version mismatch: ${normalizedBinaryVersion} != ${normalizedPackageVersion}`);
                    console.log(`Removing outdated binary and downloading new version...`);
                    safeUnlink(binaryPath);
                }
            } catch (versionErr) {
                // Binary exists but version check failed - treat as bundled binary from fresh install
                console.log(`✅ Using bundled binary (version check skipped): ${binaryPath}`);
                return; // Skip download for bundled binaries in fresh installs
            }
        } catch (err) {
            console.log(`Existing binary is invalid, will re-download: ${err.message}`);
            safeUnlink(binaryPath);
        }
    }

    // Construct download URL
    const downloadUrl = `https://github.com/${GITHUB_REPO}/releases/download/${version}/${downloadName}`;

    try {
        await downloadFile(downloadUrl, binaryPath);

        // Make binary executable on Unix systems
        if (process.platform !== 'win32') {
            fs.chmodSync(binaryPath, 0o755);
            console.log(`Made binary executable: ${binaryPath}`);
        }

        console.log(`Successfully installed binary: ${binaryPath}`);
    } catch (err) {
        console.error(`Failed to download binary: ${err.message}`);
        throw err;
    }
}

async function main() {
    console.log('Installing Gorev MCP Server binary...');

    try {
        const platformInfo = getPlatformInfo();
        console.log(`Detected platform: ${platformInfo.fullPlatform}`);

        // Get package version to match with release
        const packageJson = require('./package.json');
        let version = `v${packageJson.version}`;

        console.log(`Looking for release version: ${version}`);

        try {
            await downloadBinary(platformInfo, version);
        } catch (err) {
            console.log(`Failed to download for version ${version}, trying latest release...`);

            // Fallback to latest release if specific version not found
            try {
                const latestVersion = await getLatestReleaseVersion();
                console.log(`Using latest release: ${latestVersion}`);
                await downloadBinary(platformInfo, latestVersion);
            } catch (latestErr) {
                console.error('Failed to download from latest release as well:', latestErr.message);
                throw err; // Throw original error
            }
        }

        console.log('✅ Gorev MCP Server installation completed successfully!');
        console.log('');
        console.log('You can now use it in your MCP configuration:');
        console.log('{');
        console.log('  "mcpServers": {');
        console.log('    "gorev": {');
        console.log('      "command": "npx",');
        console.log('      "args": [');
        console.log('        "@mehmetsenol/gorev-mcp-server",');
        console.log('        "mcp-proxy"');
        console.log('      ]');
        console.log('    }');
        console.log('  }');
        console.log('}');
        console.log('');
        console.log('Note: Daemon auto-starts on first MCP connection. No manual setup required!');

    } catch (err) {
        console.error('❌ Installation failed:', err.message);
        console.error('');
        console.error('You can try:');
        console.error('1. Check your internet connection');
        console.error('2. Manually download from: https://github.com/msenol/Gorev/releases');
        console.error('3. Report the issue at: https://github.com/msenol/Gorev/issues');

        // Don't fail the install completely - let user know they can manually install
        console.error('');
        console.error('⚠️  Installation will continue, but you may need to manually install the binary.');
    }
}

if (require.main === module) {
    main();
}