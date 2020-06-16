module gitlab2.eeecs.qub.ac.uk/40178456/qse/ui

go 1.13

replace (
	gitlab2.eeecs.qub.ac.uk/40178456/qse/ad => ../ad
	gitlab2.eeecs.qub.ac.uk/40178456/qse/search => ../search
)

require (
	github.com/766b/chi-prometheus v0.0.0-20180509160047-46ac2b31aa30
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/go-chi/cors v1.0.0
	github.com/google/go-cmp v0.3.1
	github.com/prometheus/client_golang v1.3.0
	gitlab2.eeecs.qub.ac.uk/40178456/qse/ad v0.0.0
	gitlab2.eeecs.qub.ac.uk/40178456/qse/search v0.0.0
)
