
all: build clean

build:
	cp golang-reston.tex golang-reston-light.tex
	cp golang-reston.tex golang-reston-dark.tex
	sed -i 's/\%\\usecolortheme{owl}/\\usecolortheme{dove}/' golang-reston-light.tex
	sed -i 's/\%\\usecolortheme{owl}/\\usecolortheme{owl}/' golang-reston-dark.tex
	pdflatex golang-reston-light.tex
	pdflatex golang-reston-light.tex
	pdflatex golang-reston-dark.tex
	pdflatex golang-reston-dark.tex	
	rm -f golang-reston-light.tex golang-reston-dark.tex

clean:
	rm -f *.aux *.log *.nav *.out *.snm *.syn *.toc *.vrb
