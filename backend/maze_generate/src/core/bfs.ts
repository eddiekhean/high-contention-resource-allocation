import { Cell } from "./types";

export function bfs(
    start: [number, number],
    grid: Cell[][],
    rows: number,
    cols: number
) {
    const q: [number, number][] = [start];
    const dist = new Map<string, number>();
    dist.set(start.join(","), 0);

    while (q.length) {
        const [r, c] = q.shift()!;
        const d = dist.get(`${r},${c}`)!;
        const cell = grid[r][c];

        const moves: any[] = [
            [-1, 0, !cell.walls.top],
            [1, 0, !cell.walls.bottom],
            [0, -1, !cell.walls.left],
            [0, 1, !cell.walls.right],
        ];

        for (const [dr, dc, ok] of moves) {
            if (!ok) continue;
            const nr = r + dr, nc = c + dc;
            const key = `${nr},${nc}`;
            if (nr >= 0 && nc >= 0 && nr < rows && nc < cols && !dist.has(key)) {
                dist.set(key, d + 1);
                q.push([nr, nc]);
            }
        }
    }
    return dist;
}

export function findLongestPathBoundaryPoints(
    grid: Cell[][],
    rows: number,
    cols: number
): [[number, number], [number, number]] {
    const boundaries: [number, number][] = [];
    for (let r = 0; r < rows; r++)
        for (let c = 0; c < cols; c++)
            if (r === 0 || c === 0 || r === rows - 1 || c === cols - 1)
                boundaries.push([r, c]);

    const d1 = bfs([0, 0], grid, rows, cols);
    let b1: [number, number] = [0, 0];
    let max = 0;

    for (const p of boundaries) {
        const d = d1.get(p.join(","));
        if (d !== undefined && d > max) {
            max = d;
            b1 = p;
        }
    }

    const d2 = bfs(b1, grid, rows, cols);
    let b2 = b1;
    max = 0;

    for (const p of boundaries) {
        const d = d2.get(p.join(","));
        if (d !== undefined && d > max) {
            max = d;
            b2 = p;
        }
    }

    return [b1, b2];
}
