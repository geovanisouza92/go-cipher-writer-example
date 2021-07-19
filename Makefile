clean:
	rm classic-no-buffer.txt classic-buffered.txt new-no-buffer.txt new-bufferred.txt || true

run:
	for i in 1 2 3 4 5 ; do go test -benchmem -bench "^(BenchmarkWriterClassicNoBuffer)$$" >> classic-no-buffer.txt ; done
	sed -i 's/BenchmarkWriterClassicNoBuffer/BenchmarkWriter/g' classic-no-buffer.txt
	for i in 1 2 3 4 5 ; do go test -benchmem -bench "^(BenchmarkWriterClassicBufferred)$$" >> classic-bufferred.txt ; done
	sed -i 's/BenchmarkWriterClassicBufferred/BenchmarkWriter/g' classic-bufferred.txt
	for i in 1 2 3 4 5 ; do go test -benchmem -bench "^(BenchmarkWriterNewNoBuffer)$$" >> new-no-buffer.txt ; done
	sed -i 's/BenchmarkWriterNewNoBuffer/BenchmarkWriter/g' new-no-buffer.txt
	for i in 1 2 3 4 5 ; do go test -benchmem -bench "^(BenchmarkWriterNewBufferred)$$" >> new-bufferred.txt ; done
	sed -i 's/BenchmarkWriterNewBufferred/BenchmarkWriter/g' new-bufferred.txt

compare:
	benchstat classic-no-buffer.txt classic-bufferred.txt
	benchstat classic-no-buffer.txt new-no-buffer.txt
	benchstat classic-no-buffer.txt new-bufferred.txt
	benchstat classic-bufferred.txt new-no-buffer.txt
	benchstat classic-bufferred.txt new-bufferred.txt
	benchstat new-no-buffer.txt new-bufferred.txt

bench: clean run compare
