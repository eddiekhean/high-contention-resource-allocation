export type Walls = {
    top: boolean;
    right: boolean;
    bottom: boolean;
    left: boolean;
};

export type Cell = {
    x: number;
    y: number;
    walls: Walls;
    visited?: boolean;
};

export type Point = {
    x: number;
    y: number;
};

export type MazeResponse = {
    rows: number;
    cols: number;
    start: Point;
    end: Point;
    cells: Cell[];
};
