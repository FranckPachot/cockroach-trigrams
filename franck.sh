git clone git@github.com:FranckPachot/cockroach-trigrams.git
rm data/*
PGOPTIONS="-c yb_enable_upsert_mode=on" PATH="$PATH":/cygdrive/c/Program\ Files/Go/bin 

make crdb-trgrm && ./crdb-trgrm load
yb0 -c 'select count(*) from foods'
make crdb-trgrm && ./crdb-trgrm query "melk"

make crdb-trgrm && PGOPTIONS="-c pg_trgm.similarity_threshold=1 -c yb_enable_expression_pushdown=on" ./crdb-trgrm query "EGGS"
