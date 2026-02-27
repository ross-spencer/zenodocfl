"""Rudimentary validation from:

* https://github.com/crs4/rocrate-validator
"""

import argparse
import json
import os
import sys

from typing import Final

from rocrate_validator import services, models

crate_json: Final[str] = "ro-crate-metadata.json"


def check_path(path: str):
    """Strip JSON from path if exists.

    RO-CRATE validator reads from the path root and not the MD file.
    """
    if not path.endswith(crate_json):
        return path
    return path.replace(crate_json, "", 1)


def check_paths_exist(path: str) -> bool:
    """Check the paths listed in the RO-CRATE exist."""

    rocrate = {}
    with open(os.path.join(path, crate_json), "r") as json_file:
        rocrate = json.loads(json_file.read())
    try:
        files = rocrate["@graph"][1]["hasPart"]
    except (IndexError, KeyError) as err:
        print("files section does not exist in ro-crate:", err)
        return False

    noexist = []
    for idx, item in enumerate(files):
        file = os.path.join(path, item["@id"])
        if os.path.exists(file):
            continue
        noexist.append(file)
    if noexist != []:
        for item in noexist:
            print(f"item does not exist: {item}", file=sys.stderr)
        return False
    print(f"{idx} files correctly added to RO-CRATE", file=sys.stderr)
    return True


def validate(path: str, level: models.Severity, check_paths: bool = False) -> bool:
    """Validate the given path."""

    path = check_path(path)

    # Validate, see:
    # https://github.com/crs4/rocrate-validator?tab=readme-ov-file#programmatic-validation
    settings = services.ValidationSettings(
        rocrate_uri=path,
        profile_identifier="ro-crate-1.1",
        # Alternatives below are OPTIONAL and RECOMMENDED.
        requirement_severity=models.Severity.REQUIRED,
    )

    result = services.validate(settings)

    if not result.has_issues():
        if not check_paths:
            print("RO-Crate metadata is valid!", file=sys.stderr)
            return 0
        if check_paths_exist(path) is True:
            print("RO-Crate is valid!", file=sys.stderr)
            return 0

    print("RO-Crate is invalid!", file=sys.stderr)
    for issue in result.get_issues():
        # Every issue object has a reference to the check that
        # failed, the severity of the issue, and a message
        # describing the issue.
        print(
            f'Detected issue of severity {issue.severity.name} with check "{issue.check.identifier}": {issue.message}',
            file=sys.stderr,
        )
    return 1


def main():
    """Primary entry point for this script."""

    parser = argparse.ArgumentParser(
        prog="validater",
        description="validate RO-CRATE using rocrate-validator",
        epilog="part of zenodocfl... 📦",
    )

    parser.add_argument(
        "--path",
        "-p",
        help="provide a path to a RO-CRATE",
        required=True,
    )

    parser.add_argument(
        "--check-paths",
        "-c",
        help="check file paths exist",
        action="store_true",
    )

    parser.add_argument(
        "--level",
        "-l",
        help="check file paths exist",
        required=False,
        default="REQUIRED",
        type=str,
    )

    args = parser.parse_args()

    level = models.Severity.REQUIRED
    try:
        level = levels = {
            "OPTIONAL": models.Severity.OPTIONAL,
            "RECOMMENDED": models.Severity.OPTIONAL,
            "REQUIRED": models.Severity.REQUIRED,
        }[args.level.upper()]
    except KeyError:
        pass

    sys.exit(validate(args.path, level, args.check_paths))


if __name__ == "__main__":
    main()
