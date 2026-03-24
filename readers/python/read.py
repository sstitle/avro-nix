# /// script
# dependencies = ["fastavro", "typer", "rich"]
# ///

import typer
import fastavro
from rich.console import Console
from rich.table import Table
from rich import box
from pathlib import Path

console = Console()

def main(path: Path = typer.Argument(..., help="Path to the .avro file")):
    if not path.exists():
        console.print(f"[bold red]Error:[/] File not found: {path}")
        raise typer.Exit(1)

    with open(path, "rb") as f:
        avro_reader = fastavro.reader(f)
        records = list(avro_reader)
        schema = avro_reader.writer_schema

    if not records:
        console.print("[yellow]No records found.[/]")
        raise typer.Exit()

    fields = [f["name"] for f in schema["fields"]]

    table = Table(box=box.ROUNDED, header_style="bold cyan", show_lines=True)
    for field in fields:
        table.add_column(field.capitalize())

    for record in records:
        table.add_row(*[str(record[f]) for f in fields])

    console.print()
    console.rule(f"[bold]{schema['name']}[/]")
    console.print(table)
    console.print(f"[dim]{len(records)} record(s)[/]")
    console.print()

if __name__ == "__main__":
    typer.run(main)
