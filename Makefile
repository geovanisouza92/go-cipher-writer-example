bench:
	rm old.txt new.txt || true
	for i in 1 2 3 4 5 ; do go test -benchmem -bench "^(BenchmarkWriterClassic)$$" >> old.txt ; done
	sed -i 's/BenchmarkWriterClassic/BenchmarkWriter/g' old.txt
	for i in 1 2 3 4 5 ; do go test -benchmem -bench "^(BenchmarkWriterNew)$$" >> new.txt ; done
	sed -i 's/BenchmarkWriterNew/BenchmarkWriter/g' new.txt
	benchstat old.txt new.txt
