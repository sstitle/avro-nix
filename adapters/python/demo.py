#!/usr/bin/env python3
# /// script
# dependencies = ["rich"]
# ///
"""
Illustrates hexagonal architecture dependency injection.

run() is the application core: it depends only on the Repository port, never on a
concrete adapter. The adapter is selected and injected at the composition root
(here: main). Swapping memory↔RPC requires no change to run().

Usage:
  uv run adapters/python/demo.py [--adapter memory]
  uv run adapters/python/demo.py [--adapter rpc] [--addr host:port]
"""

import argparse
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).parent))

from rich.console import Console
from rich.table import Table

from item.memory import MemoryRepository
from item.model import Item
from item.repository import NotFound, Repository
from item.rpc_client import RpcClient

console = Console()


def run(repo: Repository) -> None:
    """Application core. No adapter types appear here — this is the hexagonal boundary."""

    # Save
    for it in [Item(id="1", name="Widget"), Item(id="2", name="Gadget"), Item(id="3", name="Doohickey")]:
        repo.save(it)

    # List
    console.rule("[bold]All items")
    table = Table("ID", "Name")
    for it in repo.list():
        table.add_row(it.id, it.name)
    console.print(table)

    # Get
    console.rule("[bold]Get item 2")
    found = repo.get("2")
    console.print(f"Found: {found}")

    # Delete, then confirm via typed error
    repo.delete("2")
    console.rule("[bold]After deleting item 2")
    try:
        repo.get("2")
    except NotFound as e:
        console.print(f"[yellow]Not found (expected):[/] {e}")


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--adapter", default="memory", choices=["memory", "rpc"])
    parser.add_argument("--addr", default="localhost:8080", help="RPC server address")
    args = parser.parse_args()

    # Composition root: resolve the port to a concrete adapter.
    repo: Repository
    if args.adapter == "memory":
        repo = MemoryRepository()
    else:
        repo = RpcClient(args.addr)

    console.print(f"[bold]ItemRepository demo[/] — adapter: [cyan]{args.adapter}[/]")
    run(repo)


if __name__ == "__main__":
    main()
