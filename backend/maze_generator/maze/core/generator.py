import random
from typing import List, Dict, Optional, Tuple
from .models import MazeResponse, Cell, Walls, Point
from .utils import find_longest_path_boundary_points

class MazeGenerator:
    def __init__(self, rows: int, cols: int, seed: Optional[int] = None):
        self.rows = rows
        self.cols = cols
        self.seed = seed
        self.grid: List[List[Dict[str, bool]]] = []
        
        if seed is not None:
            random.seed(seed)

    def _init_grid(self):
        # Initialize grid with all walls closed (True)
        self.grid = []
        for r in range(self.rows):
            row = []
            for c in range(self.cols):
                row.append({
                    'top': True,
                    'right': True,
                    'bottom': True,
                    'left': True,
                    'visited': False
                })
            self.grid.append(row)

    def _get_unvisited_neighbors(self, r: int, c: int) -> List[Tuple[str, int, int]]:
        neighbors = []
        # Top
        if r > 0 and not self.grid[r-1][c]['visited']:
            neighbors.append(('top', r-1, c))
        # Bottom
        if r < self.rows - 1 and not self.grid[r+1][c]['visited']:
            neighbors.append(('bottom', r+1, c))
        # Left
        if c > 0 and not self.grid[r][c-1]['visited']:
            neighbors.append(('left', r, c-1))
        # Right
        if c < self.cols - 1 and not self.grid[r][c+1]['visited']:
            neighbors.append(('right', r, c+1))
        return neighbors

    def _remove_wall(self, r1: int, c1: int, r2: int, c2: int, direction: str):
        # Remove wall from current cell
        self.grid[r1][c1][direction] = False
        
        # Remove opposite wall from neighbor
        opposite = {
            'top': 'bottom',
            'bottom': 'top',
            'left': 'right',
            'right': 'left'
        }
        self.grid[r2][c2][opposite[direction]] = False

    def _recursive_backtracker(self):
        # Start at random cell (or 0,0)
        start_r, start_c = random.randint(0, self.rows - 1), random.randint(0, self.cols - 1)
        self.grid[start_r][start_c]['visited'] = True
        stack = [(start_r, start_c)]
        
        while stack:
            r, c = stack[-1]
            neighbors = self._get_unvisited_neighbors(r, c)
            
            if neighbors:
                direction, nr, nc = random.choice(neighbors)
                self._remove_wall(r, c, nr, nc, direction)
                self.grid[nr][nc]['visited'] = True
                stack.append((nr, nc))
            else:
                stack.pop()

    def _apply_braiding(self, loop_ratio: float):
        if loop_ratio <= 0:
            return

        # Find all dead ends (cells with 3 walls)
        dead_ends = []
        for r in range(self.rows):
            for c in range(self.cols):
                walls_count = sum([
                    self.grid[r][c]['top'],
                    self.grid[r][c]['right'],
                    self.grid[r][c]['bottom'],
                    self.grid[r][c]['left']
                ])
                if walls_count == 3:
                     dead_ends.append((r, c))

        # Shuffle dead ends to ensure randomness
        random.shuffle(dead_ends)

        # Apply braiding based on loop_ratio
        for r, c in dead_ends:
            if random.random() < loop_ratio:
                 # Try to remove a wall to a connected neighbor that is NOT the one we came from
                 # Actually, we just need to connect to ANY valid neighbor that is separated by a wall
                 # Since it's a dead end, 3 walls are closed. 1 is open.
                 # We want to open one of the closed walls that leads to a valid neighbor.
                 
                candidates = []
                # Top
                if r > 0 and self.grid[r][c]['top']:
                    candidates.append(('top', r-1, c))
                # Bottom
                if r < self.rows - 1 and self.grid[r][c]['bottom']:
                    candidates.append(('bottom', r+1, c))
                # Left
                if c > 0 and self.grid[r][c]['left']:
                    candidates.append(('left', r, c-1))
                # Right
                if c < self.cols - 1 and self.grid[r][c]['right']:
                    candidates.append(('right', r, c+1))
                
                if candidates:
                    direction, nr, nc = random.choice(candidates)
                    self._remove_wall(r, c, nr, nc, direction)

    def generate(self, loop_ratio: float = 0.0) -> MazeResponse:
        self._init_grid()
        self._recursive_backtracker()
        self._apply_braiding(loop_ratio)
        
        start_node, end_node = find_longest_path_boundary_points(self.grid, self.rows, self.cols)
        
        cells = []
        for r in range(self.rows):
            for c in range(self.cols):
                # We do NOT include 'visited' in the output
                cell_data = self.grid[r][c]
                cells.append(Cell(
                    x=c, # Note: Output spec says x (col), y (row)
                    y=r,
                    walls=Walls(
                        top=cell_data['top'],
                        right=cell_data['right'],
                        bottom=cell_data['bottom'],
                        left=cell_data['left']
                    )
                ))
        
        return MazeResponse(
            rows=self.rows,
            cols=self.cols,
            start=Point(x=start_node[1], y=start_node[0]), # (r, c) -> (y, x)
            end=Point(x=end_node[1], y=end_node[0]),
            cells=cells
        )
