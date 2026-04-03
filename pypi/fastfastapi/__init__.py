import sys
from importlib.metadata import version, PackageNotFoundError

try:
    __version__ = version("fastfastapi")
except PackageNotFoundError:
    __version__ = "0.0.0"


def main() -> None:
    from ._binary import ensure_binary
    import subprocess

    binary = ensure_binary(__version__, "fastfastapi")
    result = subprocess.run([str(binary)] + sys.argv[1:])
    sys.exit(result.returncode)


def ffa() -> None:
    from ._binary import ensure_binary
    import subprocess

    binary = ensure_binary(__version__, "ffa")
    result = subprocess.run([str(binary)] + sys.argv[1:])
    sys.exit(result.returncode)
