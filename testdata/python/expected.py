#!/usr/bin/env python3
# mypy: ignore-errors
# pylint: disable=unused-import

"""Module docstring that should be stripped"""

import os
import sys

x = """
Multiline string assigned to x
# not a comment
"""

def func(a: int, b: int) -> int:
    """Function docstring — should be removed"""
    return a + b

# type: list[int]
nums = []  # type: list[int]

def main():
    """
    This is a docstring that should be removed
    """
    print(f"Hash in f-string #{len(nums)}") 