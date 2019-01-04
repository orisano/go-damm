def main():
    txt = open("damm.txt").read().strip()
    smatrice = txt.split("\n\n")
    print("package damm")
    print("")
    print("var matrices map[uint8][][]uint8 = map[uint8][][]uint8{")
    for smat in smatrice:
        mat = smat.split("\n")[1:]
        print("\t%d: [][]uint8{" % len(mat))
        for row in mat:
            print("\t\t{%s}," % ",".join(str(int(x, 10)) for x in row.strip().split(" ")))
        print("\t},")
    print("}")

if __name__ == "__main__":
    main()
