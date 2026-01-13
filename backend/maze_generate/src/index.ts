import express from "express";
import mazeRouter from "./routes/maze.route.js";

const app = express();

app.use(express.json());
app.use("/maze", mazeRouter);

app.listen(3000, () => {
    console.log("Maze service running on :3000");
});
