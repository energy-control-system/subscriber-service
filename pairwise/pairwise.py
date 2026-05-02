# /// script
# requires-python = ">=3.10"
# dependencies = [
#     "allpairspy>=2.5",
# ]
# ///
"""
Генератор тест-кейсов методом попарного тестирования (pairwise / allpairs).

Использование:
    uv run pairwise.py <входной_файл> [выходной_файл]

Формат входного файла (разделитель — табуляция):
    параметр1\tпараметр2\tпараметр3
    значение1a\tзначение2a\tзначение3a
    значение1b\tзначение2b\tзначение3b

Количество значений у параметров может различаться — пустые ячейки игнорируются.
"""

import csv
import sys
from itertools import combinations

from allpairspy import AllPairs


def parse_input(filepath: str) -> tuple[list[str], list[list[str]]]:
    """Парсинг входного файла: возвращает имена параметров и их значения."""
    with open(filepath, encoding="utf-8") as f:
        lines = [line.rstrip("\n\r") for line in f if line.strip()]

    if len(lines) < 2:
        raise ValueError("Входной файл должен содержать заголовок и хотя бы одну строку значений")

    headers = lines[0].split("\t")
    num_params = len(headers)

    param_values: list[list[str]] = [[] for _ in range(num_params)]

    for line in lines[1:]:
        parts = line.split("\t")
        for i in range(min(len(parts), num_params)):
            value = parts[i].strip()
            if value and value not in param_values[i]:
                param_values[i].append(value)

    return headers, param_values


def generate_pairwise(param_values: list[list[str]]) -> list[list[str]]:
    """Генерация тест-кейсов алгоритмом AllPairs (попарное покрытие)."""
    return [list(case) for case in AllPairs(param_values)]


def compute_new_pairs_per_case(
    headers: list[str],
    test_cases: list[list[str]],
) -> list[int]:
    """Для каждого тест-кейса подсчитывает количество НОВЫХ пар, которые он покрывает."""
    seen: set[tuple[int, int, str, str]] = set()
    result: list[int] = []

    for case in test_cases:
        new_count = 0
        for i, j in combinations(range(len(headers)), 2):
            pair = (i, j, case[i], case[j])
            if pair not in seen:
                seen.add(pair)
                new_count += 1
        result.append(new_count)

    return result


def compute_pairing_details(
    headers: list[str],
    test_cases: list[list[str]],
) -> list[dict]:
    """Формирует детали покрытия пар: какие пары покрыты какими тест-кейсами."""
    details: list[dict] = []

    for i, j in combinations(range(len(headers)), 2):
        pair_map: dict[tuple[str, str], list[int]] = {}

        for case_idx, case in enumerate(test_cases, 1):
            key = (case[i], case[j])
            pair_map.setdefault(key, []).append(case_idx)

        for (val1, val2), cases in sorted(pair_map.items()):
            details.append(
                {
                    "var1": headers[i],
                    "var2": headers[j],
                    "value1": val1,
                    "value2": val2,
                    "appearances": len(cases),
                    "cases": ", ".join(map(str, cases)),
                }
            )

    return details


def compute_coverage_stats(
    headers: list[str],
    param_values: list[list[str]],
    test_cases: list[list[str]],
) -> dict:
    """Вычисляет статистику покрытия пар."""
    total_pairs = 0
    covered_pairs = 0

    for i, j in combinations(range(len(headers)), 2):
        expected = set()
        for vi in param_values[i]:
            for vj in param_values[j]:
                expected.add((vi, vj))
        total_pairs += len(expected)

        actual = set()
        for case in test_cases:
            actual.add((case[i], case[j]))
        covered_pairs += len(actual & expected)

    return {
        "total_pairs": total_pairs,
        "covered_pairs": covered_pairs,
        "coverage_pct": covered_pairs / total_pairs * 100 if total_pairs else 0,
    }


def write_output(
    filepath: str,
    headers: list[str],
    test_cases: list[list[str]],
    pairings: list[int],
    details: list[dict],
) -> None:
    """Записывает результат в TSV-файл (совместим с Excel)."""
    with open(filepath, "w", encoding="utf-8", newline="") as f:
        writer = csv.writer(f, delimiter="\t")

        writer.writerow(["TEST CASES"])
        writer.writerow(["case"] + headers + ["pairings"])
        for idx, (case, pairing) in enumerate(zip(test_cases, pairings), 1):
            writer.writerow([idx] + case + [pairing])

        writer.writerow([])
        writer.writerow(["PAIRING DETAILS"])
        writer.writerow(["var1", "var2", "value1", "value2", "appearances", "cases"])
        for d in details:
            writer.writerow([d["var1"], d["var2"], d["value1"], d["value2"], d["appearances"], d["cases"]])


def print_table(headers: list[str], test_cases: list[list[str]], pairings: list[int]) -> None:
    """Красивый вывод таблицы тест-кейсов в консоль."""
    col_widths = [max(4, len("case"))]
    for i, h in enumerate(headers):
        max_val = max((len(case[i]) for case in test_cases), default=0)
        col_widths.append(max(len(h), max_val))
    col_widths.append(len("pairings"))

    header_row = f"  {'case':>{col_widths[0]}}"
    for i, h in enumerate(headers):
        header_row += f"  {h:<{col_widths[i + 1]}}"
    header_row += f"  {'pairings':>{col_widths[-1]}}"

    print(header_row)
    print("  " + "-" * (len(header_row) - 2))

    for idx, (case, pairing) in enumerate(zip(test_cases, pairings), 1):
        row = f"  {idx:>{col_widths[0]}}"
        for i, val in enumerate(case):
            row += f"  {val:<{col_widths[i + 1]}}"
        row += f"  {pairing:>{col_widths[-1]}}"
        print(row)


def main() -> None:
    if len(sys.argv) < 2:
        print(__doc__)
        sys.exit(1)

    input_file = sys.argv[1]
    output_file = sys.argv[2] if len(sys.argv) > 2 else "test_cases.tsv"

    # --- Чтение входных данных ---
    headers, param_values = parse_input(input_file)

    print(f"Входной файл:  {input_file}")
    print(f"Выходной файл: {output_file}")
    print()
    print(f"Параметров: {len(headers)}")
    for h, vals in zip(headers, param_values):
        print(f"  {h}: {len(vals)} значений → {vals}")

    full_count = 1
    for vals in param_values:
        full_count *= len(vals)

    # --- Генерация тест-кейсов ---
    test_cases = generate_pairwise(param_values)

    print()
    print(f"Полный перебор:        {full_count} тест-кейсов")
    print(f"Попарное покрытие:     {len(test_cases)} тест-кейсов")
    print(f"Сокращение:            {100 - len(test_cases) * 100 / full_count:.1f}%")

    # --- Статистика покрытия ---
    stats = compute_coverage_stats(headers, param_values, test_cases)
    print(f"Покрытие пар:          {stats['covered_pairs']}/{stats['total_pairs']} ({stats['coverage_pct']:.1f}%)")

    # --- Вывод таблицы ---
    pairings = compute_new_pairs_per_case(headers, test_cases)
    details = compute_pairing_details(headers, test_cases)

    print()
    print("=" * 70)
    print("TEST CASES")
    print("=" * 70)
    print_table(headers, test_cases, pairings)

    # --- Запись в файл ---
    write_output(output_file, headers, test_cases, pairings, details)
    print()
    print(f"Результат записан в: {output_file}")


if __name__ == "__main__":
    main()
