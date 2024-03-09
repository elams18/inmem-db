from typing import List, Optional, Union
import sys


class Person:
    def __init__(self, name: str, gender: str) -> None:
        self.name = name
        self.gender = gender


class RelationshipManager:
    def __init__(self) -> None:
        self.relationships = []
        self.relationship_map = {
            ("father", "female"): "daughter",
            ("mother", "female"): "daughter",
            ("father", "male"): "son",
            ("mother", "male"): "son",
            ("son", "male"): "father",
            ("daughter", "male"): "father",
            ("son", "female"): "mother",
            ("daughter", "female"): "mother",
            ("husband", "male"): "wife",
            ("wife", "male"): "husband",
            ("husband", "female"): "wife",
            ("wife", "female"): "husband",
        }
        self.plural_mapper = {
            "sons": "son",
            "son": "son",
            "daughters": "daughter",
            "daughter": "daughter",
            "husbands": "husband",
            "husband": "husband",
            "wives": "wife",
            "wife": "wife",
            "fathers": "father",
            "father": "father",
            "mothers": "mother",
            "mother": "mother",
        }
        self.allowed_relations = []

    def add_relationship(
        self, person1: Person, relation: str, person2: Person
    ) -> Optional[int]:
        if relation not in self.allowed_relations:
            print(f"{relation} not in added relation")
            return -1

        self.relationships.append([person1, relation, person2])
        # print(person1.name, relation, person2.name)
        inverse_relation = self.get_inverse_relation(person2.gender, relation)
        # print(inverse_relation, person2.gender)
        if inverse_relation:
            self.relationships.append([person2, inverse_relation, person1])

    def get_inverse_relation(self, p2_gender: str, relation: str) -> Optional[str]:
        # print(p2_gender, relation, self.relationship_map[(relation, p2_gender)])
        if (relation, p2_gender) in self.relationship_map.keys():
            return self.relationship_map[(relation, p2_gender)]
        return None

    def add_relations(self, relation: str) -> None:
        self.allowed_relations.append(relation.lower())


class FamilyTree:
    def __init__(self) -> None:
        self.people = {}
        self.relationships = RelationshipManager()

    def add_person(self, name: str, gender: str) -> Optional[int]:
        if name not in self.people:
            self.people[name] = Person(name, gender)
        else:
            print(f"Person with name {name} already exists")
            return -1

    def add_relationship(
        self, name1: str, relationship: str, name2: str
    ) -> Optional[int]:
        relationship = relationship.lower()
        # print(self.people, name2, name1)
        if name1 in self.people and name2 in self.people:
            person1 = self.people[name1]
            person2 = self.people[name2]
            res = self.relationships.add_relationship(person1, relationship, person2)
            if res == -1:
                return res
        else:
            print(f"One or more people are not present")
            return -1

    def count_relationship(self, name: str, relationship_type: str) -> int:
        if name in self.people:
            person = self.people[name]
            count = sum(
                [
                    1
                    for p, r, p2 in self.relationships.relationships
                    if p2 == person and r == relationship_type
                ]
            )
            return count
        else:
            print(f"Person with name {name} is not present")
            return -1

    def get_relationship(self, name: str, relation: str) -> Union[List[str], int]:
        if name in self.people:
            person = self.people[name]
            # assuming that there may be more than one relation
            for p, r, p2 in self.relationships.relationships:
                print(p.name, r, p2.name)

            related_people = [
                p.name
                for p, r, p2 in self.relationships.relationships
                if p2 == person and r == relation
            ]
            return related_people
        else:
            print("Person with name {name} is not present")
            return -1


def test():
    # Example usage:
    family_tree = FamilyTree()

    # Add persons to the family tree
    family_tree.add_person("John", "male")
    family_tree.add_person("Jane", "female")
    family_tree.add_person("Tom", "male")
    family_tree.add_person("Lucy", "female")

    # Establish relationships
    family_tree.add_relationship("John", "father", "Tom")
    family_tree.add_relationship("Jane", "mother", "Tom")
    family_tree.add_relationship("John", "husband", "Jane")
    # family_tree.add_relationship("Jane", "wife", "John")

    # Count relationships
    assert family_tree.count_relationship("John", "father") == 1  # Should print 0
    assert family_tree.count_relationship("John", "husband") == 1  # Should print 1
    assert family_tree.count_relationship("Jane", "husband") == 0  # Should print 1
    assert family_tree.get_relationship("Tom", "father") == ["John"]


def test_case_2():
    family_tree = FamilyTree()
    family_tree.add_person("John", "male")
    family_tree.add_person("Jane", "female")
    family_tree.add_person("Tom", "male")
    family_tree.add_relationship("John", "father", "Tom")
    print(family_tree.count_relationship("John", "father"))
    assert family_tree.count_relationship("John", "father") == 1


if __name__ == "__main__":
    if len(sys.argv) == 2 and sys.argv[1] == "test":
        test_case_2()
        exit()
    print("family_tree CLI Tool to store connections between family members.")
    print("Usage: As mentioned in DOCS")
    family_tree = FamilyTree()
    while True:
        print("> ", end="")
        command_args = input().split()
        n: int = len(command_args)
        if n == 0:
            continue
        if command_args[0] != "family_tree":
            print("Usage: family_tree [options]")
            continue

        if n > 1 and command_args[1] == "add":
            if n != 4 and n != 5:
                print("Usage: family_tree add [person|realtionship] <name> <gender>|M")
                continue
            por = command_args[2]
            if por == "person":
                if n == 4:
                    gender: str = "male"
                elif n == 5:
                    gender: str = command_args[4]
                else:
                    print(
                        "Usage: family_tree add [person|realtionship] <name> <gender>|M"
                    )
                    continue
                name = command_args[3]
                res = family_tree.add_person(name, gender)
                if res == -1:
                    continue
                print(f"Added person with name: {name} in the family_tree")
                continue
            elif por == "relationship":
                name = command_args[3]
                res = family_tree.relationships.add_relations(name)
                if res == -1:
                    continue
                print(f"Relationship: {name} added into the family tree")
            else:
                print("Usage: family_tree add [person|realtionship] <name>")
        elif n > 1 and command_args[1] == "connect":
            if n != 7:
                print("Usage: family_tree connect <name1> as <relationship> of <name2>")
                continue
            name = command_args[2]
            relationship = command_args[4]
            name2 = command_args[6]
            res = family_tree.add_relationship(name, relationship, name2)
            if res == -1:
                continue
            print(f"Added relationship: {name} and {name2} : {relationship}")
        elif n > 1 and command_args[1] == "count":
            if n != 5:
                print("Usage: family_tree count <plural_relationship> of <name>")
                continue
            plural_relationship = command_args[2]
            relationship = family_tree.relationships.plural_mapper.get(
                plural_relationship, None
            )
            name = command_args[4]
            if not relationship:
                print(f"{plural_relationship} not recognised as valid relationship")
                continue
            res = family_tree.count_relationship(name, relationship)
            if res == -1:
                print("Usage: family_tree count <plural_relationship> of <name>")
            print(f"Count: {res}")
        elif n > 1 and command_args[1] in family_tree.relationships.allowed_relations:
            if n != 4:
                print("Usage: family_tree <relationship> of <name>")
                continue
            relationship = command_args[1]
            name = command_args[3]
            res = family_tree.get_relationship(name, relationship)
            if res != -1:
                print(f"{relationship} of {name} is {res}")
        else:
            print("Usage: family_tree options")
