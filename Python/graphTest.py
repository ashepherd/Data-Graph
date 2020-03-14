from rdflib import Graph

g = Graph()
g.parse("./test.nt", format="nt")

qres = g.query(
    """SELECT ?s
       WHERE {
          ?s ?p ?o .
       }
       LIMIT 10""")

for row in qres:
    print(row)

print(qres)


