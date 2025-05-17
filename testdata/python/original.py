#!/usr/bin/env python3
# mypy: ignore-errors
# pylint: disable=unused-import

"""Module docstring that should be stripped"""

import os  # inline
import sys  # another inline

# Empty comment lines
#
#
#

x = """
Multiline string assigned to x
# not a comment
"""

def func(a: int, b: int) -> int:  # trailing comment
    """Function docstring — should be removed"""
    # TODO: something
    return a + b  # end-of-line

# type: list[int]
nums = []  # type: list[int]

def main():
    print(f"Hash in f-string #{len(nums)}")  # comment

# Final line comment 