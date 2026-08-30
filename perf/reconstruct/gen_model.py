import sys

n = int(sys.argv[1]) if len(sys.argv) > 1 else 1000

print("# Synthetic sysl model for reproducing a frozen (arr-ai/frozen) performance")
print("# regression. Purely generated data -- no real application names/structure.")
print()
for i in range(n):
    nxt = (i + 1) % n
    print(f"App{i} [~service, owner=\"team{i % 20}\"]:")
    print(f"    @description = \"Synthetic service {i}\"")
    print(f"    Get(id <: int) [~rest]:")
    print(f"        App{nxt} <- Get")
    print(f"        return ok <: App{i}.Data")
    print()
    print(f"    !type Data:")
    print(f"        id <: int [~pk]")
    print(f"        name <: string:")
    print(f"            @description = \"Synthetic field\"")
    print(f"        value <: int?")
    print()
