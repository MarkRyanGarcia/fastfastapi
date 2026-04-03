"""Downloads and caches the fastfastapi binary from GitHub Releases."""

import os
import platform
import shutil
import ssl
import stat
import sys
import tarfile
import urllib.request
import zipfile
from pathlib import Path

GITHUB_REPO = "markryangarcia/fastfastapi"
CACHE_DIR = Path.home() / ".cache" / "fastfastapi"


def _platform_asset() -> str:
    system = platform.system().lower()
    machine = platform.machine().lower()

    arch = "arm64" if machine in ("arm64", "aarch64") else "amd64"

    if system == "darwin":
        return f"fastfastapi_darwin_{arch}.tar.gz"
    elif system == "linux":
        return f"fastfastapi_linux_{arch}.tar.gz"
    elif system == "windows":
        return f"fastfastapi_windows_{arch}.zip"
    else:
        raise RuntimeError(f"Unsupported platform: {system}/{machine}")


def _bin_path(version: str, bin_name: str) -> Path:
    suffix = ".exe" if platform.system().lower() == "windows" else ""
    return CACHE_DIR / version / (bin_name + suffix)


def _ssl_context() -> ssl.SSLContext:
    """Return an SSL context, falling back gracefully on macOS cert issues."""
    try:
        import certifi
        return ssl.create_default_context(cafile=certifi.where())
    except ImportError:
        pass

    ctx = ssl.create_default_context()
    # On macOS with python.org builds, system certs may not be available.
    # Try to load them from the macOS keychain path if present.
    mac_certs = "/etc/ssl/cert.pem"
    if platform.system().lower() == "darwin" and os.path.exists(mac_certs):
        ctx.load_verify_locations(mac_certs)
        return ctx

    # Last resort: disable verification (still encrypted, just not verified)
    ctx.check_hostname = False
    ctx.verify_mode = ssl.CERT_NONE
    return ctx


def ensure_binary(version: str, bin_name: str = "fastfastapi") -> Path:
    bin_path = _bin_path(version, bin_name)
    if bin_path.exists():
        return bin_path

    asset = _platform_asset()
    url = f"https://github.com/{GITHUB_REPO}/releases/download/v{version}/{asset}"

    bin_path.parent.mkdir(parents=True, exist_ok=True)

    archive_path = bin_path.parent / asset
    print(f"Downloading fastfastapi v{version}...", file=sys.stderr)

    ctx = _ssl_context()
    with urllib.request.urlopen(url, context=ctx) as response, open(archive_path, "wb") as out:
        shutil.copyfileobj(response, out)

    if asset.endswith(".zip"):
        with zipfile.ZipFile(archive_path) as zf:
            zf.extractall(bin_path.parent)
    else:
        with tarfile.open(archive_path, "r:gz") as tf:
            tf.extractall(bin_path.parent)

    archive_path.unlink()

    # Mark all extracted binaries executable
    for name in ("fastfastapi", "ffa"):
        p = _bin_path(version, name)
        if p.exists():
            p.chmod(p.stat().st_mode | stat.S_IEXEC | stat.S_IXGRP | stat.S_IXOTH)

    return bin_path
