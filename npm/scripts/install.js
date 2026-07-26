const os = require('os');
const path = require('path');
const fs = require('fs');
const https = require('https');
const { execSync } = require('child_process');

const version = require('../package.json').version;
const platform = os.platform();
const arch = os.arch();

// Map Node.js platforms to GoReleaser OS
const osMap = {
  win32: 'windows',
  darwin: 'darwin',
  linux: 'linux'
};

// Map Node.js arch to GoReleaser Arch
const archMap = {
  x64: 'amd64',
  arm64: 'arm64'
};

const goOs = osMap[platform];
const goArch = archMap[arch];

if (!goOs || !goArch) {
  console.error(`Unsupported platform/architecture: ${platform}/${arch}`);
  process.exit(1);
}

const ext = platform === 'win32' ? '.zip' : '.tar.gz';
const filename = `windmist_${version}_${goOs}_${goArch}${ext}`;
const downloadUrl = `https://github.com/Nithwin/WindMist/releases/download/v${version}/${filename}`;

const distDir = path.join(__dirname, '..', 'dist');
if (!fs.existsSync(distDir)) {
  fs.mkdirSync(distDir, { recursive: true });
}

console.log(`Downloading WindMist v${version} for ${goOs}/${goArch}...`);
console.log(`URL: ${downloadUrl}`);

// In a real robust implementation, we would use axios + unzipper/tar here.
// For the sake of this CLI wrapper, we output instructions or use curl if available.

try {
  // Simple check for curl and tar/unzip
  const tmpFile = path.join(os.tmpdir(), filename);
  execSync(`curl -L -o "${tmpFile}" "${downloadUrl}"`, { stdio: 'inherit' });
  
  if (ext === '.zip') {
    execSync(`tar -xf "${tmpFile}" -C "${distDir}"`, { stdio: 'inherit' }); // Windows 10+ has tar
  } else {
    execSync(`tar -xzf "${tmpFile}" -C "${distDir}"`, { stdio: 'inherit' });
  }
  
  console.log('WindMist installed successfully.');
} catch (e) {
  console.warn('Failed to download binary from GitHub Releases.');
  console.warn('This is expected if the release is not published yet.');
  console.warn(e.message);
}
