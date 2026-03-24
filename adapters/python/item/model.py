from dataclasses import dataclass


@dataclass
class Item:
    """Domain entity. Derived from com.example.Item in ports/item/repository.avpr."""

    id: str
    name: str
