# FAMILY-TREE

## Overview:
    1. A CLI tool to depict the family structure.
    2. Assumes Unique Names in Family.
    3. Assumes a single family for now -> Immediate enhancement: To add ID mapped into the family

## Design:

This program has a OOP approach where there are three main entities: 
    1. Person
    2. Family Tree
    3. Relationship Manager

*   Person -> Stores the state of the person (name, gender)
*   Relationship Manager -> Strategy to encapsulate different strategies for managing relationships
*   Family Tree -> Facade for creating system of Person objects and their relationships. 

## Future evolution:
    1. Multiple families with IDs mapped when creating family
    2. Issue with same-sex and poly marriages are already addressed, but the configuration can be mapped better with the implementation. For example, now even if we add "father" in relationships, we can query "sons of".
    3. Getting ancestors can be next step. (Eg. great great grandfather) which can be solved by DFS.
