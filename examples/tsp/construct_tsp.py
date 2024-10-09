import random
import numpy as np
import matplotlib.pyplot as plt

random.seed(1337)

# Generate 20 random points
points = [[random.random(), random.random()] for _ in range(20)]
# Find the distance matrix
dist_mat = []

for i in range(20):
    dists = []

    for j in range(20):
        dists.append(np.linalg.norm(np.array(points[i]) - np.array(points[j])))

    dist_mat.append(dists)

with open("dist_mat", "w") as f:
    for row in dist_mat:
        for dist in row:
            f.write(f"{dist} ")

        f.write("\n")

# Edges output by the algorithm
edges = """
(11, 3)
(3, 10)
(10, 13)
(13, 6)
(6, 15)
(15, 12)
(12, 0)
(0, 18)
(18, 1)
(1, 19)
(19, 8)
(8, 16)
(16, 14)
(14, 9)
(9, 4)
(4, 5)
(5, 17)
(17, 7)
(7, 2)
(2, 11)
""".split("\n")[1:-1]
edges = [(int(edge.split(", ")[0][1:]), int(edge.split(", ")[1][:-1])) for edge in edges]
# Plot the points with the selected path
x = [p[0] for p in points]
y = [p[1] for p in points]

plt.scatter(x, y, color="blue")

for i, j in edges:
    plt.plot([x[i], x[j]], [y[i], y[j]], color="red")

plt.xlabel("X-axis")
plt.ylabel("Y-axis")
plt.title("TSP Solution Found by ACO Algorithm")
plt.show()
