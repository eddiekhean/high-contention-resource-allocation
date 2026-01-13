import seedrandom from "seedrandom";
import { Cell, MazeResponse } from "./types";
import { findLongestPathBoundaryPoints } from "./bfs";

export class MazeGenerator {
    rows: number;
    cols: number;
    grid: Cell[][];
    rng: () => number;

    constructor(rows: number, cols: number, seed?: number) {
        this.rows = rows;
        this.cols = cols;
        this.rng = seed ? seedrandom(seed.toString()) : Math.random;
        this.grid = [];
        this.initGrid();
    }

    private initGrid() {
        for (let r = 0; r < this.rows; r++) {
            const row: Cell[] = [];
            for (let c = 0; c < this.cols; c++) {
                row.push({
                    x: c,
                    y: r,
                    visited: false,
                    walls: { top: true, right: true, bottom: true, left: true },
                });
            }
            this.grid.push(row);
        }
    }

    private neighbors(r: number, c: number) {
        const res: any[] = [];
        if (r > 0 && !this.grid[r - 1][c].visited)
            res.push(["top", r - 1, c]);
        if (r < this.rows - 1 && !this.grid[r + 1][c].visited)
            res.push(["bottom", r + 1, c]);
        if (c > 0 && !this.grid[r][c - 1].visited)
            res.push(["left", r, c - 1]);
        if (c < this.cols - 1 && !this.grid[r][c + 1].visited)
            res.push(["right", r, c + 1]);
        return res;
    }

    private removeWall(r1: number, c1: number, r2: number, c2: number, dir: string) {
        const opp: Record<string, keyof Cell["walls"]> = { top: "bottom", bottom: "top", left: "right", right: "left" };
        this.grid[r1][c1].walls[dir as keyof Cell["walls"]] = false;
        this.grid[r2][c2].walls[opp[dir]] = false;
    }

    private dfs() {
        const stack: [number, number][] = [];
        const sr = Math.floor(this.rng() * this.rows);
        const sc = Math.floor(this.rng() * this.cols);

        this.grid[sr][sc].visited = true;
        stack.push([sr, sc]);

        while (stack.length) {
            const [r, c] = stack[stack.length - 1];
            const neigh = this.neighbors(r, c);

            if (neigh.length) {
                const [dir, nr, nc] = neigh[Math.floor(this.rng() * neigh.length)];
                this.removeWall(r, c, nr, nc, dir);
                this.grid[nr][nc].visited = true;
                stack.push([nr, nc]);
            } else {
                stack.pop();
            }
        }
    }

    private braid(loopRatio: number) {
        if (loopRatio <= 0) return;

        for (let r = 0; r < this.rows; r++) {
            for (let c = 0; c < this.cols; c++) {
                const cell = this.grid[r][c];
                const wallCount = Object.values(cell.walls).filter(Boolean).length;

                if (wallCount === 3 && this.rng() < loopRatio) {
                    const candidates: any[] = [];
                    if (r > 0 && cell.walls.top) candidates.push(["top", r - 1, c]);
                    if (r < this.rows - 1 && cell.walls.bottom) candidates.push(["bottom", r + 1, c]);
                    if (c > 0 && cell.walls.left) candidates.push(["left", r, c - 1]);
                    if (c < this.cols - 1 && cell.walls.right) candidates.push(["right", r, c + 1]);

                    if (candidates.length) {
                        const [dir, nr, nc] =
                            candidates[Math.floor(this.rng() * candidates.length)];
                        this.removeWall(r, c, nr, nc, dir);
                    }
                }
            }
        }
    }

    generate(loopRatio = 0): MazeResponse {
        this.dfs();
        this.braid(loopRatio);

        const [start, end] = findLongestPathBoundaryPoints(this.grid, this.rows, this.cols);

        return {
            rows: this.rows,
            cols: this.cols,
            start: { x: start[1], y: start[0] },
            end: { x: end[1], y: end[0] },
            cells: this.grid.flat().map(({ visited, ...c }) => c),
        };
    }
}
