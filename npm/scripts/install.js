const os = require('os');
const path = require('path');
const fs = require('fs');
const { execSync } = require('child_process');

const version = require('../package.json').version;
const platform = os.platform();
const arch = os.arch();

// Map Node.js platforms to release archive OS names (matches GoReleaser + install.sh)
const osMap = {
  win32: 'Windows',
  darwin: 'macOS',
  linux: 'Linux'
};

// Map Node.js arch to release archive arch names
const archMap = {
  x64: 'x86_64',
  arm64: 'arm64'
};

const releaseOs = osMap[platform];
const releaseArch = archMap[arch];

if (!releaseOs || !releaseArch) {
  console.error(`Unsupported platform/architecture: ${platform}/${arch}`);
  process.exit(1);
}

const ext = platform === 'win32' ? '.zip' : '.tar.gz';
const filename = `windmist_${version}_${releaseOs}_${releaseArch}${ext}`;
const downloadUrl = `https://github.com/Nithwin/WindMist/releases/download/v${version}/${filename}`;

const distDir = path.join(__dirname, '..', 'dist');
if (!fs.existsSync(distDir)) {
  fs.mkdirSync(distDir, { recursive: true });
}

console.log(`Downloading WindMist v${version} for ${releaseOs}/${releaseArch}...`);
console.log(`URL: ${downloadUrl}`);

try {
  const tmpFile = path.join(os.tmpdir(), filename);
  execSync(`curl -L -o "${tmpFile}" "${downloadUrl}"`, { stdio: 'inherit' });
  
  if (ext === '.zip') {
    execSync(`tar -xf "${tmpFile}" -C "${distDir}"`, { stdio: 'inherit' });
  } else {
    execSync(`tar -xzf "${tmpFile}" -C "${distDir}"`, { stdio: 'inherit' });
  }
  
  console.log('WindMist installed successfully.');
} catch (e) {
  console.warn('Failed to download binary from GitHub Releases.');
  console.warn('This is expected if the release is not published yet.');
  console.warn(e.message);
}
