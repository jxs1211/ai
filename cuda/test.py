import numpy as np
from timeit import default_timer as timer
from numba import vectorize

def multiplyVector(a, b, c):
	for i in range(a.size):
		c[i] = a[i] * b[i]

@vectorize(["float64(float64,float64)"], target = "cuda")

def multiplyVector2(a, b):
	return a * b

def main():
	N = 64000000
	a = np.ones(N, dtype=np.float64)
	b = np.ones(N, dtype=np.float64)
	c = np.ones(N, dtype=np.float64)

	start = timer()
	# multiplyVector(a,b,c)
	multiplyVector2(a, b)
	print(c[: 10])
	print("time consumption: ", timer() - start)


if __name__ == "__main__":
	main()