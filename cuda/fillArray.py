import numpy as np
from timeit import default_timer as timer
from numba import jit

@jit("(float32)", target="gpu")
def fillArrayWithGPU(a):
  for k in range(a.size):
    a[k] += 1
  

def main():
  N = 100000000
  a = np.ones(N, dtype=np.float32)
  start = timer()
  fillArrayWithGPU(a)
  print("time consumption: ", timer() - start)


main()